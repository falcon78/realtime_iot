const MAX_DATA_NUMS = 100

let bar_chart = null

const currentScheme = String(window.location).split(":")[0]

let websocketHost = (currentScheme === "https" ? "wss://" : "ws://") + window.location.hostname
if (window.location.port.length > 0) {
    websocketHost += ":" + window.location.port
}

const channelId = new URL(window.location).searchParams.get("channelId")

async function getGraphData(id) {
    try {
        const res = await fetch(
            `/api/records/${id}?limit=${MAX_DATA_NUMS}`,
            {headers: {'Content-Type': 'application/json'}}
        )
        return await res.json()
    } catch (e) {
        alert(e)
    }
}

export default {
    name: 'BarComponent',

    setup() {
        return {};
    },

    async mounted() {
        const accessKey = new URL(window.location).searchParams.get("accessKey")

        const socket = new WebSocket(`${websocketHost}/ws/${accessKey}`)
        const response = await getGraphData(channelId)

        let channelData = {
            0: [], // labels
            1: [],
            2: [],
            3: [],
            4: [],
        }
        const colors = ['rgba(226,135,67, 0.9)', 'rgba(42,157,143,0.9)', 'rgba(247,37,133,0.9)', 'rgba(74,78,105,0.9)']
        for (const res of response.records) {
            channelData[0].push(new Date(res.timestamp))
            channelData[1].push(res.channelOne)
            channelData[2].push(res.channelTwo)
            channelData[3].push(res.channelThree)
            channelData[4].push(res.channelFour)
        }

        const datasets = []
        for (let i = 1; i < 5; i++) {
            datasets.push({
                label: "Channel" + i,
                data: channelData[i],
                borderColor: colors[i],
                fill: false,
                borderWidth: 3,
                radius: 1
            })
        }

        let data = {
            type: "line",
            data: {
                labels: channelData[0],
                datasets: datasets
            },
            options: {
                responsive: true,
                normalized: false,
                lineTension: 1,
                scales: {
                    yAxes: [{
                        display: true,
                        ticks: {
                            beginAtZero: true,
                            steps: 10,
                            stepValue: 5,
                            max: Math.ceil(response.max),
                            min: Math.ceil(response.min)
                        }
                    }],
                    xAxes: [{
                        type: 'time',
                        time: {
                            tooltipFormat: 'YYYY-MM-DD HH:mm',
                            displayFormats: {
                                millisecond: 'HH:mm:ss',
                                second: 'HH:mm:ss',
                                minute: 'HH:mm',
                                hour: 'HH'
                            }
                        },
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: 'Time'
                        }
                    }],
                }
            }
        }

        const ctx = document.getElementById('chart');
        bar_chart = new Chart(ctx, data);
        bar_chart.options.tooltips.enabled = false;
        bar_chart.options.animation = false;

        socket.addEventListener('open', function (event) {
            console.log("socket open")
        })
        socket.addEventListener('message', async (event) => {
            const data = JSON.parse(await event.data.text())

            bar_chart.data.datasets[0].data.push(data.channelOne)
            bar_chart.data.datasets[1].data.push(data.channelTwo)
            bar_chart.data.datasets[2].data.push(data.channelThree)
            bar_chart.data.datasets[3].data.push(data.channelFour)
            bar_chart.data.labels.push(new Date(data.timestamp))
            for (let i = 0; i < 4; i++) {
                if (bar_chart.data.datasets[i].data.length > MAX_DATA_NUMS) {
                    bar_chart.data.datasets[i].data.shift()
                }
            }
            if (bar_chart.data.labels.length > MAX_DATA_NUMS) {
                bar_chart.data.labels.shift()
            }

            bar_chart.update()
        })
    },

    template: `
      <div>
          <canvas id="chart"></canvas>
      </div>
    `,
};
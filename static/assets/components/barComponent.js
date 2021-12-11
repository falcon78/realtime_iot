const MAX_DATA_NUMS = 100

let bar_chart = null

let arrayData = []
for (let i = 0; i < MAX_DATA_NUMS; i++) {
    arrayData.push(Math.random() * 200)
}

let labelData = []
for (let i = 0; i < MAX_DATA_NUMS; i++) {
    labelData.push(("data" + Math.random().toString()).substr(0, 10))
}

const channelId = new URL(window.location).searchParams.get("channelId")

async function getGraphData(id) {
    try {
        const res = await fetch(
            `/api/records/${id}`,
            {headers: {'Content-Type': 'application/json'}}
        )
        return await res.json()
    } catch (e) {
        alert(e)
    }
}

export default {
    name: 'One',

    setup() {
        return {};
    },

    async mounted() {
        const response = await getGraphData(channelId)

        let channelData = {
            1: [],
            2: [],
            3: [],
            4: []
        }
        const colors = ['rgba(226,135,67, 0.9)', 'rgba(42,157,143,0.9)', 'rgba(247,37,133,0.9)', 'rgba(74,78,105,0.9)']
        for (const res of response) {
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
                labels: labelData,
                datasets: datasets
            },
            options: {
                responsive: true,
                normalized: true,
                lineTension: 1,
                scales: {
                    yAxes: [
                        {
                            ticks: {
                                beginAtZero: true,
                                padding: 0
                            }
                        }
                    ],
                    x: {
                        type: 'linear'
                    }
                }
            }
        }

        const ctx = document.getElementById('planet-chart');
        bar_chart = new Chart(ctx, data);
        // bar_chart.options.tooltips.enabled = false;
        bar_chart.options.animation = false;

        // setInterval(() => {
        //     bar_chart.data.datasets[0].data.push(Math.random() * 200)
        //     // bar_chart.data.datasets[1].data.push(Math.random() * 200)
        //     bar_chart.data.labels.push(("data" + Math.random().toString()).substr(0, 10))
        //
        //     if (bar_chart.data.datasets[0].data.length > MAX_DATA_NUMS) {
        //         bar_chart.data.datasets[0].data.shift()
        //         // bar_chart.data.datasets[1].data.shift()
        //         bar_chart.data.labels.shift()
        //     }
        //
        //     bar_chart.update()
        // }, 500)
    },

    template: `
      <div>
      <canvas id="planet-chart"></canvas>
      </div>
    `,
};
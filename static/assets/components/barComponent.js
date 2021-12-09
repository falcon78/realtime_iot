import store from '../store.js'

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

let data = {
    type: "line",
    data: {
        labels: labelData,
        datasets: [
            {
                label: "Number of Moons",
                data: arrayData,
                borderColor: "rgba(0, 181, 204, 1)",
                fill: false,
                borderWidth: 3,
                radius: 1
            },
        ]
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

export default {
    name: 'One',

    setup() {
        return { store };
    },

    mounted() {
        const ctx = document.getElementById('planet-chart');
        bar_chart = new Chart(ctx, data);
        // bar_chart.options.tooltips.enabled = false;
        bar_chart.options.animation = false;

        setInterval(() => {
            bar_chart.data.datasets[0].data.push(Math.random() * 200)
            // bar_chart.data.datasets[1].data.push(Math.random() * 200)
            bar_chart.data.labels.push(("data" + Math.random().toString()).substr(0, 10))

            if (bar_chart.data.datasets[0].data.length > MAX_DATA_NUMS) {
                bar_chart.data.datasets[0].data.shift()
                // bar_chart.data.datasets[1].data.shift()
                bar_chart.data.labels.shift()
            }

            bar_chart.update()
        }, 500)
    },

    template: `
    <div>
      <canvas id="planet-chart"></canvas>
    </div> 
    `,
};
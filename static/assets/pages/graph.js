import barComponent from "../components/barComponent.js";

export default {
    name: 'Realtime Graph',
    components: {barComponent},

    setup() {
        const title = 'リアルタイムグラフ'
        return {title}
    },

    template: `
      <div>
      {{ title }}
      <barComponent></barComponent>
      </div>
    `,
};
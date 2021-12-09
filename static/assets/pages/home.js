import barComponent from '../components/barComponent.js'

export default {
    name: 'Home',
    components: {barComponent},

    setup() {
        const title = 'Home page'
        return {title}
    },

    template: `
        <div>
            {{ title }}
            <barComponent></barComponent>
        </div>
    `,
};
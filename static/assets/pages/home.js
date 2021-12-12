export default {
    name: 'Home',

    setup() {
        const title = 'Home page'
        return {title}
    },

    data() {
        return {
            channels: [],
            newChannelName: ""
        }
    },

    mounted() {
        this.fetchChannels()
    },

    methods: {
        fetchChannels() {
            fetch(
                "/api/channels",
                {headers: {'Content-Type': 'application/json'}}
            )
                .then(async (res) => {
                    this.channels = await res.json()
                })
                .catch(err => {
                    alert(err)
                })
        },

        submit() {
            if (this.newChannelName.trim().length === 0) {
                return
            }

            const meta = {
                headers: {'Content-Type': 'application/json'},
                method: "POST"
            }

            fetch(`/api/channel/create/${this.newChannelName}`, meta)
                .then(async (res) => {
                    this.newChannelName = ""
                    this.fetchChannels()
                })
                .catch(err => {
                    alert(err)
                })
        },

        delete(id) {
            const meta = {
                headers: {'Content-Type': 'application/json'},
                method: "DELETE"
            }

            fetch(`/api/channel/delete/${id}`, meta)
                .then(async (res) => {
                    this.fetchChannels()
                })
                .catch(err => {
                    alert(err)
                })
        }
    },

    template: `
      <div>
      <div style="display: flex; justify-content: center">
        <input v-model="newChannelName" placeholder="チャンネル名">
        <button v-on:click="submit">チャンネル追加</button>
      </div>

      <table>
        <thead>
        <tr>
          <th>名前</th>
          <th>アクセスキー</th>
          <th>削除</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="channel in channels" v-bind:key="channel.id">
          <td>
            <a v-bind:href="'/graph' + '?channelId=' + channel.id + '&accessKey=' + channel.accessKey">
              {{channel.name}}
            </a>
          </td>
          <td>{{channel.accessKey}}</td>
          <td><a v-on:click="() => this.delete(channel.id)">delete</a></td>
        </tr>
        </tbody>
      </table>


      </div>
    `,
};
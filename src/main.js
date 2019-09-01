import Vue from 'vue'
import Vuex from 'vuex'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import axios from 'axios'

Vue.use(Vuex)
Vue.config.productionTip = false
Vue.prototype.$http = axios
Vue.prototype.$hostname = (Vue.config.productionTip) ? 'https://www.your-api.com' : 'http://localhost:8090'
axios.defaults.headers.common['Content-Type'] = 'application/json'

const store = new Vuex.Store({
  state: {
    startTime: '12:00',
    endTime: '',
    showAddEventForm: false,
    focus: ''
  },
  mutations: {
    setStartTime: (state, v) => state.startTime = v,
    setEndTime: (state, v) => state.endTime = v,
    setShowAddEventForm: (state, v) => state.showAddEventForm = v,
    setFocus: (state, v) => state.focus = v,
  }
})

new Vue({
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')

import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import axios from 'axios'
import { store } from './_store'

Vue.prototype.$http = axios
axios.defaults.headers.common['Content-Type'] = 'application/json'

new Vue({
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')

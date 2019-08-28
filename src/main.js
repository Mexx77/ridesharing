import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios';

Vue.config.productionTip = false;

Vue.prototype.$http = axios;
Vue.prototype.$hostname = (Vue.config.productionTip) ? 'https://www.your-api.com' : 'http://localhost:8090'

axios.defaults.headers.common['Content-Type'] = 'application/json';

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')

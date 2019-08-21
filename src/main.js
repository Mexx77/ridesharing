import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios';

Vue.config.productionTip = false;

Vue.prototype.$http = axios;
// read only api key
axios.defaults.headers.common['x-apikey'] = '5d5d7d77a592085130522576';
axios.defaults.headers.common['Content-Type'] = 'application/json';

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')

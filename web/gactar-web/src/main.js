import Vue from 'vue'
import App from './App.vue'

import axios from 'axios'
import Buefy from 'buefy'
import VueAxios from 'vue-axios'

require('./app.scss')

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(Buefy)

new Vue({
  render: (h) => h(App),
}).$mount('#app')

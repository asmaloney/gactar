import Vue from 'vue'
import App from './App.vue'

import axios from 'axios'
import { Button, Field, Input, Select, Upload } from 'buefy'
import VueAxios from 'vue-axios'

require('./app.scss')

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(Button)
Vue.use(Field)
Vue.use(Input)
Vue.use(Select)
Vue.use(Upload)

new Vue({
  render: (h) => h(App),
}).$mount('#app')

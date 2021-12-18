import Vue from 'vue'
import App from './App.vue'

// axios
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.use(VueAxios, axios)

// Buefy
import {
  ConfigProgrammatic,
  Button,
  Dropdown,
  Field,
  Input,
  Select,
  Tabs,
  Upload,
} from 'buefy'

Vue.use(Button)
Vue.use(Dropdown)
Vue.use(Field)
Vue.use(Input)
Vue.use(Select)
Vue.use(Tabs)
Vue.use(Upload)
ConfigProgrammatic.setOptions({
  defaultIconPack: 'far',
})

// Fontawesome icons
// Icons are found here: https://fontawesome.com/v5.15/icons?d=gallery&p=2
import { dom, library } from '@fortawesome/fontawesome-svg-core'
import { faCaretSquareDown } from '@fortawesome/free-regular-svg-icons'
import {
  faFileDownload,
  faFileUpload,
  faRunning,
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faCaretSquareDown, faFileDownload, faFileUpload, faRunning)
dom.watch()

Vue.component('FontAwesomeIcon', FontAwesomeIcon)

Vue.config.productionTip = false

import './app.scss'

new Vue({
  render: (h) => h(App),
}).$mount('#app')

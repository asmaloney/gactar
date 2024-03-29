import Vue from 'vue'
import App from './App.vue'

// axios
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.use(VueAxios, axios)

// Buefy
// Currently we cannot import individual components using TypeScript.
// So we have to import all of buefy.
import Buefy from 'buefy'
Vue.use(Buefy, {
  defaultIconPack: 'far',
})

// import {
//   ConfigProgrammatic,
//   Dropdown,
//   Field,
//   Input,
//   Select,
//   Tabs,
//   Upload,
// } from 'buefy'

// Vue.use(Button)
// Vue.use(Dropdown)
// Vue.use(Field)
// Vue.use(Input)
// Vue.use(Select)
// Vue.use(Tabs)
// Vue.use(Upload)
// ConfigProgrammatic.setOptions({
//   defaultIconPack: 'far',
// })

// Fontawesome icons
// Icons are found here: https://fontawesome.com/v5.15/icons?d=gallery&p=2
import { dom, library } from '@fortawesome/fontawesome-svg-core'
import { faCaretSquareDown } from '@fortawesome/free-regular-svg-icons'
import {
  faCaretDown,
  faCaretUp,
  faFileDownload,
  faFileUpload,
  faRunning,
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(
  faCaretDown,
  faCaretUp,
  faCaretSquareDown,
  faFileDownload,
  faFileUpload,
  faRunning
)
dom.watch()

Vue.component('FontAwesomeIcon', FontAwesomeIcon)

Vue.config.productionTip = false

// SCSS
import './app.scss'

// Our internal API
import api from './api'

api.init(parseInt(window.location.port))

new Vue({
  render: (h) => h(App),
}).$mount('#app')

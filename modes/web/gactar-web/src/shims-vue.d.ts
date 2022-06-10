// This lets us import from .vue files without producing errors.
// Without it, vetur will complain about "import App from './App.vue'"
// See: // https://vuejs.github.io/vetur/guide/setup.html#typescript

declare module '*.vue' {
  import Vue from 'vue'
  export default Vue
}

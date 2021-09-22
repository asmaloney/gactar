<template>
  <b-tab-item label="amod">
    <div class="columns buttons">
      <div class="column">
        <b-field class="is-pulled-right" grouped>
          <b-dropdown aria-role="list">
            <template #trigger="{}">
              <b-button
                label="Load Example"
                type="is-info is-light"
                :icon-right="'caret-square-down'"
              />
            </template>

            <b-dropdown-item
              v-for="option in exampleFiles"
              :key="option"
              :value="option"
              aria-role="listitem"
              :focusable="false"
              @click="getExample(option)"
            >
              {{ option }}
            </b-dropdown-item>
          </b-dropdown>

          <save-button
            :code="amodCode"
            default-name="model"
            file-extension="amod"
          />

          <b-field class="file">
            <b-button type="upload is-info is-light is-outlined ml-2">
              <span class="fa fa-file-upload icon-space" />
              <b-upload v-model="fileToLoad" accept=".amod,text/plain">
                Load
              </b-upload>
            </b-button>
          </b-field>
        </b-field>
      </div>
    </div>

    <code-mirror
      :key="count"
      ref="code-editor"
      :amod-code.sync="amodCode"
      mode="amod"
      framework="amod"
    />
  </b-tab-item>
</template>

<script>
import CodeMirror from './CodeMirror'
import SaveButton from './SaveButton'

const localStorageName = 'gactar-code-editor'

export default {
  components: { CodeMirror, SaveButton },

  data() {
    return {
      amodCode: '',
      exampleFiles: [],
      fileToLoad: null,
      loadedFromLocal: false,

      // This is used to prevent caching of the code-mirror data.
      // See https://stackoverflow.com/questions/48400302/vue-js-not-updating-props-in-child-when-parent-component-is-changing-the-propert
      count: 0,
    }
  },

  watch: {
    amodCode() {
      this.$emit('codeChange', this.amodCode)
    },

    // watch for a change in fileToLoad, then load it
    fileToLoad(val) {
      if (val == null) {
        return
      }

      var reader = new FileReader()
      reader.onload = (e) => {
        this.$refs['code-editor'].setCode(e.target.result)
      }
      reader.readAsText(this.fileToLoad)
    },
  },

  created() {
    window.addEventListener('beforeunload', this.beforeWindowUnload)
    window.addEventListener('load', this.onWindowLoad)
  },

  beforeDestroy() {
    window.removeEventListener('load', this.onWindowLoad)
    window.removeEventListener('beforeunload', this.beforeWindowUnload)
  },

  async mounted() {
    await this.getExamples()
    if (!this.loadedFromLocal) {
      this.getExample(this.exampleFiles[0])
    }
  },

  methods: {
    beforeWindowUnload() {
      localStorage.setItem(localStorageName, this.amodCode)
    },

    async getExample(example) {
      try {
        const { data } = await this.$http.get('/examples/' + example)
        this.count += 1
        this.amodCode = data
      } catch (err) {
        this.$emit('showError', err)
      }
    },

    async getExamples() {
      try {
        const { data } = await this.$http.get('/examples/list')
        this.exampleFiles = data.example_list
      } catch (err) {
        this.$emit('showError', err)
      }
    },

    onWindowLoad() {
      // check for a local save and use it instead of loading an example
      var code = localStorage.getItem(localStorageName)
      if (code !== null) {
        this.loadedFromLocal = true
        this.$refs['code-editor'].setCode(code)
      }
    },
  },
}
</script>

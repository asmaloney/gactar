<template>
  <div>
    <h1>
      Andy's Fancy ACT-R Thingamabob (a.k.a.
      <a href="https://github.com/asmaloney/gactar" target="_">gactar</a>)
    </h1>
    <section class="section p-0 pt-4">
      <div class="columns">
        <div class="column is-three-fifths">
          <div class="columns">
            <div class="column">
              <b-dropdown aria-role="list">
                <template #trigger="{}">
                  <b-button
                    label="Load Example"
                    type="is-info"
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
            </div>

            <div class="column">
              <b-field class="is-pulled-right" grouped>
                <b-button type="is-info ml-2" @click="saveCode">
                  <span class="fa fa-file-download icon-space" /> Save To File
                </b-button>

                <b-field class="file is-info ml-2">
                  <b-upload
                    v-model="fileToLoad"
                    class="file-label"
                    accept=".amod,text/plain"
                  >
                    <span class="file-cta">
                      <span class="fa fa-file-upload icon-space" />
                      <span class="file-label">Load From File</span>
                    </span>
                  </b-upload>
                </b-field>
              </b-field>
            </div>
          </div>
        </div>

        <div class="column">
          <b-field label="Goal" label-position="on-border">
            <b-input
              v-model="goal"
              placeholder="(initial goal here)"
              expanded
            />
            <p class="control">
              <b-button class="button is-info" :loading="running" @click="run">
                <span class="fa fa-running icon-space" />Run
              </b-button>
            </p>
          </b-field>
        </div>
      </div>

      <div class="columns">
        <div class="column is-three-fifths">
          <code-mirror
            :key="count"
            ref="code-editor"
            :amod-code.sync="amodCode"
          />
        </div>
        <div class="column">
          <textarea id="results" v-model="results"></textarea>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import CodeMirror from './components/CodeMirror'

const localStorageName = 'gactar-code-editor'

export default {
  components: { CodeMirror },

  data() {
    return {
      amodCode: '',
      exampleFiles: [],
      fileToLoad: null,
      goal: '',
      loadedFromLocal: false,
      running: false,
      results: '',

      // This is used to prevent caching of the code-mirror data.
      // See https://stackoverflow.com/questions/48400302/vue-js-not-updating-props-in-child-when-parent-component-is-changing-the-propert
      count: 0,
    }
  },

  watch: {
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
        this.showError(err)
      }
    },

    async getExamples() {
      try {
        const { data } = await this.$http.get('/examples/list')
        this.exampleFiles = data.example_list
      } catch (err) {
        this.showError(err)
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

    async run() {
      this.running = true
      try {
        const { data } = await this.$http.post('/run', {
          amod: this.amodCode,
          run: this.goal,
        })

        if (data.error) {
          this.showError(data.error)
          return
        }

        this.setResults(data.results)
      } catch (err) {
        this.showError(err)
      }
    },

    saveCode() {
      // Adapted from: https://stackoverflow.com/a/51315312
      var codeAsBlob = new Blob([this.amodCode], {
        type: 'text/plain;charset=utf-8',
      })

      var downloadLink = document.createElement('a')
      downloadLink.download = 'model.amod'
      downloadLink.innerHTML = 'Save File'

      if (window.webkitURL != null) {
        // Chrome allows the link to be clicked without actually adding it to the DOM.
        downloadLink.href = window.webkitURL.createObjectURL(codeAsBlob)
      } else {
        // Firefox requires the link to be added to the DOM before it can be clicked.
        downloadLink.href = window.URL.createObjectURL(codeAsBlob)
        downloadLink.style.display = 'none'
        downloadLink.onclick = (e) => {
          document.body.removeChild(e.target)
        }
        document.body.appendChild(downloadLink)
      }

      downloadLink.click()
    },

    setResults(results) {
      let text = ''
      for (const [key, value] of Object.entries(results)) {
        text += key + '\n' + '---\n'
        text += value
        text += '\n\n'
      }
      this.results = text
      this.running = false
    },

    showError(err) {
      this.results = err
      this.running = false
    },
  },
}
</script>

<style scoped>
.icon-space {
  margin-right: 0.5em;
}
</style>

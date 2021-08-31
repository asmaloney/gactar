<template>
  <div>
    <h1>
      Andy's Fancy ACT-R Thingamabob (a.k.a.
      <a href="https://github.com/asmaloney/gactar" target="_">gactar</a>)
    </h1>
    <section class="section p-0 pt-4">
      <div class="columns">
        <div class="column is-three-fifths">
          <b-field label="Load Example" label-position="on-border">
            <b-select
              v-model="selectedExample"
              placeholder="Select Example"
              :loading="!exampleFiles.length"
              @input="handleSelectExample"
            >
              <option
                v-for="option in exampleFiles"
                :key="option"
                :value="option"
              >
                {{ option }}
              </option>
            </b-select>
          </b-field>
        </div>
        <div class="column">
          <b-field label="Goal" label-position="on-border">
            <b-input
              v-model="goal"
              placeholder="(initial goal here)"
              expanded
            />
            <p class="control">
              <b-button
                class="button is-primary"
                :loading="running"
                @click="run"
              >
                Run
              </b-button>
            </p>
          </b-field>
        </div>
      </div>

      <div class="columns">
        <div class="column is-three-fifths">
          <code-mirror
            :key="count"
            :amod-code.sync="amodCode"
            ref="code-editor"
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
      goal: '',
      loadedFromLocal: false,
      running: false,
      results: '',
      selectedExample: null,

      // This is used to prevent caching of the code-mirror data.
      // See https://stackoverflow.com/questions/48400302/vue-js-not-updating-props-in-child-when-parent-component-is-changing-the-propert
      count: 0,
    }
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
      this.handleSelectExample(this.selectedExample)
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
        if (!this.loadedFromLocal) {
          this.selectedExample = this.exampleFiles[0]
        }
      } catch (err) {
        this.showError(err)
      }
    },

    async handleSelectExample(example) {
      await this.getExample(example)
      this.selectedExample = null
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

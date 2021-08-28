<template>
  <div>
    <h1>
      Andy's Fancy ACT-R Thingamabob (a.k.a.
      <a href="https://github.com/asmaloney/gactar" target="_">gactar</a>)
    </h1>
    <section class="section p-0 pt-4">
      <div class="columns">
        <div class="column is-three-fifths">
          <b-field label="Example" label-position="on-border">
            <b-select
              v-model="selectedExample"
              placeholder="Select an example"
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
          <code-mirror :key="count" :amod-code.sync="amodCode"></code-mirror>
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

export default {
  components: { CodeMirror },

  data() {
    return {
      amodCode: '',
      exampleFiles: [],
      goal: '',
      running: false,
      results: '',
      selectedExample: '',

      // This is used to prevent caching of the code-mirror data.
      // See https://stackoverflow.com/questions/48400302/vue-js-not-updating-props-in-child-when-parent-component-is-changing-the-propert
      count: 0,
    }
  },

  async mounted() {
    await this.getExamples()
    this.handleSelectExample(this.selectedExample)
  },

  methods: {
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
        this.selectedExample = this.exampleFiles[0]
      } catch (err) {
        this.showError(err)
      }
    },

    async handleSelectExample(example) {
      await this.getExample(example)
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
      this.running = false
    },

    setResults(results) {
      let text = ''
      for (const [key, value] of Object.entries(results)) {
        text += key + '\n' + '---\n'
        text += value
        text += '\n\n'
      }
      this.results = text
    },

    showError(err) {
      this.results = err
    },
  },
}
</script>

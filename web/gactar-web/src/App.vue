<template>
  <span>
    <h1>
      <img src="/images/gactar-logo.svg" />
      gactar-web
      <span v-if="version" class="version-number">
        &nbsp;(<a href="https://github.com/asmaloney/gactar" target="_">{{
          version
        }}</a>
        )
      </span>
    </h1>
    <div class="columns">
      <div class="column is-three-fifths code-column">
        <b-tabs
          v-model="activeTab"
          class="custom"
          :animated="false"
          expanded
          :type="'is-boxed'"
        >
          <amod-code-tab @codeChange="codeChange" @showError="showError" />

          <template v-for="tab in tabs">
            <code-tab
              v-if="tab.displayed"
              :key="tab.id"
              :value="tab.id"
              :framework="tab.id"
              :mode="tab.mode"
              :file-extension="tab.fileExtension"
              :model-name="tab.modelName"
              :code="code[tab.id]"
            >
            </code-tab>
          </template>
        </b-tabs>
      </div>

      <div class="column">
        <div class="columns buttons">
          <div class="column">
            <b-field label="Goal" label-position="on-border">
              <b-input
                v-model="goal"
                placeholder="(initial goal here)"
                expanded
              />
              <p class="control">
                <b-button
                  class="button is-info"
                  :loading="running"
                  @click="run"
                >
                  <span class="fa fa-running icon-space" />Run
                </b-button>
              </p>
            </b-field>
          </div>
        </div>

        <div class="columns result">
          <div class="column">
            <textarea id="results" v-model="results"></textarea>
          </div>
        </div>
      </div>
    </div>
  </span>
</template>

<script>
import AmodCodeTab from './components/AmodCodeTab'
import CodeTab from './components/CodeTab'

export default {
  components: { AmodCodeTab, CodeTab },

  data() {
    return {
      activeTab: 0,
      baseTabs: [
        {
          id: 'ccm',
          mode: 'python',
          fileExtension: 'py',
          modelName: '',
          displayed: false,
        },
        {
          id: 'pyactr',
          mode: 'python',
          fileExtension: 'py',
          modelName: '',
          displayed: false,
        },
        {
          id: 'vanilla',
          mode: 'commonlisp',
          fileExtension: 'lisp',
          modelName: '',
          displayed: false,
        },
      ],

      code: {},
      goal: '',
      running: false,
      results: '',
      version: null,
    }
  },

  computed: {
    tabs() {
      return this.baseTabs
    },
  },

  async mounted() {
    this.version = this.getVersion()
  },

  methods: {
    codeChange(newCode) {
      this.code['amod'] = newCode
    },

    async run() {
      this.running = true
      try {
        const { data } = await this.$http.post('/run', {
          amod: this.code['amod'],
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

    async getVersion() {
      try {
        const { data } = await this.$http.get('/version')

        this.version = data.version
      } catch (err) {
        this.showError(err)
      }
    },

    setResults(results) {
      let text = ''
      for (const [key, value] of Object.entries(results)) {
        text += key + '\n' + '---\n'
        text += value.output
        text += '\n\n'

        this.code[key] = value.code

        const index = this.tabs.findIndex((obj) => obj.id == key)
        if (index != -1) {
          this.tabs[index].modelName = value.modelName

          // show our tabs the first time we have code
          if (value.code.length != 0) {
            this.tabs[index].displayed = true
          }
        }
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

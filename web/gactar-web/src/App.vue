<template>
  <span>
    <h1>
      <img src="/images/gactar-logo.svg" />
      gactar-web
      <span v-if="version" class="version-number">
        &nbsp;(<a href="https://github.com/asmaloney/gactar" target="_">
          {{ version }}
        </a>
        )
      </span>
    </h1>
    <div class="tile is-ancestor">
      <div class="tile is-vertical is-7 code-tile">
        <b-tabs
          v-model="activeTab"
          class="code-tabs"
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

      <div class="tile is-vertical is-parent">
        <div v-if="tabs.length > 0" class="tile is-child is-12">
          <b-field label="Select Frameworks" custom-class="is-small">
            <b-checkbox-button
              v-for="tab in tabs"
              v-model="selectedFrameworks"
              type="is-info"
              size="is-small"
              :native-value="tab.id"
              :key="tab.id"
              expanded
              class="ml-1 mr-1"
              @input="frameworkChanged"
            >
              <span>{{ tab.id }}</span>
            </b-checkbox-button>
          </b-field>
        </div>

        <div class="tile is-child is-12">
          <b-field label="Goal" label-position="on-border">
            <b-input
              v-model="goal"
              placeholder="(initial goal here)"
              expanded
            />
            <p class="control">
              <b-button type="is-info" :loading="running" @click="run">
                <span class="fa fa-running icon-space" />Run
              </b-button>
            </p>
          </b-field>
        </div>

        <div class="tile is-child">
          <textarea id="results" v-model="results" expanded></textarea>
        </div>
      </div>
    </div>
  </span>
</template>

<script lang="ts">
import Vue from 'vue'

import api, {
  FrameworkInfo,
  FrameworkInfoList,
  ResultMap,
  RunResult,
  Version,
} from './api'

import AmodCodeTab from './components/AmodCodeTab.vue'
import CodeTab from './components/CodeTab.vue'

interface Tab {
  id: string
  mode: string
  fileExtension: string
  modelName: string
  displayed: boolean
}

type CodeMap = { [key: string]: string }

interface Data {
  activeTab: number
  baseTabs: Tab[]
  code: CodeMap

  goal: string
  running: boolean
  results: string

  availableFrameworks: string[]
  selectedFrameworks: string[]
  version: Version
}

const selectedFrameworksStorageName = 'gactar.selected-frameworks'

export default Vue.extend({
  components: { AmodCodeTab, CodeTab },

  data(): Data {
    return {
      activeTab: 0,
      baseTabs: [],

      code: {},
      goal: '',
      running: false,
      results: '',
      availableFrameworks: [],
      selectedFrameworks: [],
      version: '',
    }
  },

  computed: {
    tabs(): Tab[] {
      return this.baseTabs
    },
  },

  created() {
    window.addEventListener('load', () => {
      this.onWindowLoad()
    })
  },

  mounted() {
    this.loadFrameworks()
    this.loadVersion()
  },

  methods: {
    frameworkChanged() {
      // Save our selected frameworks
      localStorage.setItem(
        selectedFrameworksStorageName,
        JSON.stringify(this.selectedFrameworks)
      )
    },

    codeChange(newCode: string) {
      this.code['amod'] = newCode
    },

    hideTabsNotInUse() {
      this.baseTabs.forEach((tab: Tab) => {
        if (!this.selectedFrameworks.includes(tab.id)) {
          tab.displayed = false
        }
      })
    },

    loadFrameworks() {
      api
        .getFrameworks()
        .then((list: FrameworkInfoList) => {
          list.forEach((info: FrameworkInfo) => {
            // create tab info for each language present on the server
            const tab: Tab = {
              id: info.name,
              mode: info.language,
              fileExtension: info.fileExtension,
              modelName: '',
              displayed: false,
            }

            this.availableFrameworks.push(info.name)

            this.baseTabs.push(tab)
          })
        })
        .catch((err: Error) => {
          this.showError(err.message)
        })
    },

    loadVersion() {
      api
        .getVersion()
        .then((version: Version) => {
          this.version = version
        })
        .catch((err: Error) => {
          this.showError(err.message)
        })
    },

    onWindowLoad() {
      // Load our selected frameworks from local storage (if any)
      var frameworks = localStorage.getItem(selectedFrameworksStorageName)
      if (frameworks === null) {
        this.selectedFrameworks = this.availableFrameworks
      } else {
        this.selectedFrameworks = JSON.parse(frameworks) as string[]

        // Filter the saved list by the available frameworks.
        const availableFrameworks = this.availableFrameworks // need this const because we can't use "this" inside filter
        this.selectedFrameworks = this.selectedFrameworks.filter(function (
          name: string
        ) {
          return availableFrameworks.includes(name)
        })
      }

      window.removeEventListener('load', () => {
        this.onWindowLoad()
      })
    },

    run() {
      this.running = true

      this.hideTabsNotInUse()

      api
        .run(this.code['amod'], this.goal, this.selectedFrameworks)
        .then((results: RunResult) => {
          if ('results' in results) {
            this.setResults(results.results)
          } else {
            this.showError(results.error)
          }
        })
        .catch((err: Error) => {
          this.showError(err.message)
        })
    },

    setResults(results: ResultMap) {
      let text = ''
      for (const [key, value] of Object.entries(results)) {
        text += key + '\n' + '---\n'
        text += value.output
        text += '\n\n'

        this.code[key] = value.code

        const index = this.tabs.findIndex((obj: Tab) => obj.id == key)
        if (index != -1) {
          this.tabs[index].modelName = value.modelName

          // show our tabs the first time we have code
          if (value.code && value.code.length != 0) {
            this.tabs[index].displayed = true
          }
        }
      }

      this.results = text
      this.running = false
    },

    showError(err: string) {
      this.results = err
      this.running = false
    },
  },
})
</script>

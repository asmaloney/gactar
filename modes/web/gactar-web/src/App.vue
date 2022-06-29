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
          <amod-code-tab
            :amod-issues="amodIssues"
            @codeChange="amodCodeChange"
            @showError="showError"
          />

          <template v-for="fw in frameworks">
            <framework-code-tab
              v-if="fw.showTabs"
              :model-name="currentModelName"
              :key="fw.info.name"
              :value="fw.info.name"
              :framework="fw.info"
              :code="fw.code"
            >
            </framework-code-tab>
          </template>
        </b-tabs>
      </div>

      <div class="tile is-vertical is-parent">
        <div v-if="availableFrameworks.length > 0" class="tile is-child is-12">
          <b-field label="Select Frameworks" custom-class="is-small">
            <b-checkbox-button
              v-for="name in availableFrameworks"
              v-model="selectedFrameworks"
              type="is-info"
              size="is-small"
              :native-value="name"
              :key="name"
              expanded
              class="ml-1 mr-1"
              @input="frameworkChanged"
            >
              <span>{{ name }}</span>
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

        <div class="tile is-child results">
          <b-tabs :animated="false" :type="'is-boxed'" size="is-small" expanded>
            <b-tab-item label="All Results">
              <textarea
                v-model="allResults"
                class="textarea is-info has-fixed-size"
                expanded
                readonly
              />
            </b-tab-item>
            <template v-for="fw in frameworks">
              <b-tab-item
                v-if="fw.showTabs"
                :label="fw.info.name"
                :key="fw.info.name"
              >
                <textarea
                  :value="fw.output"
                  class="textarea is-info has-fixed-size"
                  expanded
                  readonly
                />
              </b-tab-item>
            </template>
          </b-tabs>
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
  FrameworkResultMap,
  IssueList,
  RunParams,
  RunResult,
  Version,
} from './api'

import { commentString, issuesToArray } from './utils'

import AmodCodeTab from './components/AmodCodeTab.vue'
import FrameworkCodeTab from './components/FrameworkCodeTab.vue'

interface FrameworkData {
  info: FrameworkInfo

  code: string // generated code
  output: string // output from run

  showTabs: boolean
}

interface Data {
  amodCode: string
  amodIssues: IssueList

  activeTab: number

  currentModelName: string
  goal: string
  running: boolean
  allResults: string

  frameworks: FrameworkData[]
  availableFrameworks: string[]
  selectedFrameworks: string[]
  version: Version
}

const selectedFrameworksStorageName = 'gactar.selected-frameworks'

export default Vue.extend({
  components: { AmodCodeTab, FrameworkCodeTab },

  data(): Data {
    return {
      amodCode: '',
      amodIssues: [],

      activeTab: 0,

      currentModelName: '',
      goal: '',
      running: false,
      allResults: '',

      frameworks: [],
      availableFrameworks: [],
      selectedFrameworks: [],
      version: '',
    }
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
    amodCodeChange(newCode: string) {
      this.amodCode = newCode
    },

    clearResults() {
      this.allResults = ''
    },

    frameworkChanged() {
      // Save our selected frameworks
      localStorage.setItem(
        selectedFrameworksStorageName,
        JSON.stringify(this.selectedFrameworks)
      )
    },

    hideTabsNotInUse() {
      this.frameworks.forEach((data: FrameworkData) => {
        if (!this.selectedFrameworks.includes(data.info.name)) {
          data.showTabs = false
        }
      })
    },

    loadFrameworks() {
      api
        .getFrameworks()
        .then((list: FrameworkInfoList) => {
          list.forEach((info: FrameworkInfo) => {
            const frameworkData: FrameworkData = {
              info: info,
              code: '',
              output: '',
              showTabs: false,
            }
            this.frameworks.push(frameworkData)
            this.availableFrameworks.push(info.name)
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

      this.clearResults()
      this.hideTabsNotInUse()

      const params: RunParams = {
        amod: this.amodCode,
        goal: this.goal,
        frameworks: this.selectedFrameworks,
      }

      api
        .run(params)
        .then((result: RunResult) => {
          if (result.issues) {
            this.showIssues(result.issues)
          }
          if (result.results) {
            this.setResults(result.results)
          }
          this.running = false
        })
        .catch((err: Error) => {
          this.showError(err.message)
          this.running = false
        })
    },

    setResults(results: FrameworkResultMap) {
      let text = '' // text we add to "all results"

      for (const [frameworkName, result] of Object.entries(results)) {
        this.currentModelName = result.modelName

        let frameworkData = this.frameworks.find(
          (item) => item.info.name == frameworkName
        )

        if (frameworkData == null) {
          return
        }

        text += frameworkName + '\n' + '---\n'

        let frameworkOutput = '' // text for the framework tab's results

        if (result.issues) {
          const issueTexts = issuesToArray(result.issues)
          frameworkOutput = issueTexts.join('\n') + '\n\n'
          text += frameworkOutput
        }

        frameworkData.output = frameworkOutput + (result.output || '')

        if (result.output) {
          text += result.output
          text += '\n\n'
        }

        if (result.code) {
          frameworkData.code = result.code
        } else {
          frameworkData.code = commentString(
            frameworkData.info.language,
            '(No code returned from server)'
          )
        }

        // show our tabs the first time we have code
        if (frameworkData.code.length != 0) {
          frameworkData.showTabs = true
        }
      }

      this.allResults += text
    },

    showError(err: string) {
      this.allResults = err
    },

    showIssues(list: IssueList) {
      Vue.set(this, 'amodIssues', list)

      const issueTexts = issuesToArray(list)

      this.allResults += issueTexts.join('\n') + '\n\n'
    },
  },
})
</script>

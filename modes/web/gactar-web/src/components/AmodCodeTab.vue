<template>
  <b-tab-item label="amod">
    <div class="columns buttons">
      <div class="column">
        <b-field class="is-pulled-right" grouped>
          <b-dropdown aria-role="list">
            <template #trigger="{}">
              <b-button type="is-info is-light">
                Load Example
                <span class="far fa-caret-square-down icon-space-left" />
              </b-button>
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
      :amod-code="amodCode"
      :amod-issues="amodIssues"
      mode="amod"
      editorID="amod"
      @editorCodeChange="editorCodeChange"
    />
  </b-tab-item>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

import api, { ExampleList, IssueList } from '../api'

import CodeMirror from './CodeMirror.vue'
import SaveButton from './SaveButton.vue'

interface Data {
  amodCode: string
  exampleFiles: ExampleList
  fileToLoad: File | null
  loadedFromLocal: boolean
  count: number
}

const codeEditorStorageName = 'gactar.code-editor'

export default defineComponent({
  components: { CodeMirror, SaveButton },

  props: {
    amodIssues: {
      type: Array as PropType<IssueList>,
      required: false,
    },
  },

  data(): Data {
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
    fileToLoad(file: File) {
      var reader = new FileReader()
      reader.onload = (ev: ProgressEvent<FileReader>) => {
        if (ev.target != null && typeof ev.target.result === 'string') {
          this.amodCode = ev.target.result
        }
      }
      reader.readAsText(file)
    },
  },

  created() {
    window.addEventListener('beforeunload', () => {
      this.beforeWindowUnload()
    })
    window.addEventListener('load', () => {
      this.onWindowLoad()
    })
  },

  beforeDestroy() {
    window.removeEventListener('load', () => {
      this.onWindowLoad()
    })
    window.removeEventListener('beforeunload', () => {
      this.beforeWindowUnload()
    })
  },

  // Disable lint while waiting for this fix:
  //  https://github.com/vuejs/core/pull/5914
  // eslint-disable-next-line @typescript-eslint/no-misused-promises
  async mounted() {
    await this.getExamples()
    if (!this.loadedFromLocal) {
      await this.getExample(this.exampleFiles[0])
    }
  },

  methods: {
    beforeWindowUnload() {
      localStorage.setItem(codeEditorStorageName, this.amodCode)
    },

    editorCodeChange(code: string) {
      this.$emit('codeChange', code)
    },

    async getExample(example: string) {
      await api
        .getExample(example)
        .then((code: string) => {
          this.count += 1
          this.amodCode = code
        })
        .catch((err: Error) => {
          this.$emit('showError', err)
        })
    },

    async getExamples() {
      await api
        .getExampleList()
        .then((list: ExampleList) => {
          this.exampleFiles = list
        })
        .catch((err: Error) => {
          this.$emit('showError', err)
        })
    },

    onWindowLoad() {
      // check for a local save and use it instead of loading an example
      var code = localStorage.getItem(codeEditorStorageName)
      if (code !== null) {
        this.loadedFromLocal = true
        this.amodCode = code
      }
    },
  },
})
</script>

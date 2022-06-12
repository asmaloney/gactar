<template>
  <b-tab-item class="code-tab" :label="framework.name">
    <div class="columns buttons">
      <div class="column">
        <strong>{{ defaultFileName }}.{{ framework.fileExtension }}</strong>
        (generated for {{ framework.executableName }})
        <b-field class="is-pulled-right">
          <save-button
            :code="code"
            :default-name="defaultFileName"
            :file-extension="framework.fileExtension"
          />
        </b-field>
      </div>
    </div>

    <code-mirror
      :key="count"
      :ref="refName"
      :amod-code="code"
      :mode="framework.language"
      :editorID="framework.name"
      :read-only="true"
    />
  </b-tab-item>
</template>

<script lang="ts">
import Vue, { PropType } from 'vue'

import { FrameworkInfo } from '@/api'

import CodeMirror from './CodeMirror.vue'
import SaveButton from './SaveButton.vue'

interface Data {
  fileToLoad: string | null
  accept: string
  refName: string
  count: number
}

export default Vue.extend({
  components: { CodeMirror, SaveButton },

  props: {
    code: {
      type: String,
      required: true,
    },
    framework: {
      type: Object as PropType<FrameworkInfo>,
      required: true,
    },
    modelName: {
      type: String,
      required: true,
    },
  },

  data(): Data {
    return {
      fileToLoad: null,

      accept: '.' + this.framework.language + ',text/plain',
      refName: 'code-editor-' + this.framework.language,

      // This is used to prevent caching of the code-mirror data.
      // See https://stackoverflow.com/questions/48400302/vue-js-not-updating-props-in-child-when-parent-component-is-changing-the-propert
      count: 0,
    }
  },

  computed: {
    defaultFileName(): string {
      return this.framework.name + '_' + this.modelName
    },
  },

  watch: {
    code() {
      this.count += 1
    },
  },
})
</script>

<template>
  <div>
    <textarea :id="id" v-model="code"></textarea>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import CodeMirror from 'codemirror'

// Add-ons
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/selection/active-line'
import 'codemirror/mode/commonlisp/commonlisp'
import 'codemirror/mode/python/python'

import '../codemirror/amod'

export default Vue.extend({
  props: {
    amodCode: {
      type: String,
      default() {
        return '(code here)'
      },
    },
    mode: {
      type: String,
      required: true,
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
    framework: {
      type: String,
      required: true,
    },
  },

  data() {
    return {
      editor: null,
      id: 'id-' + this.framework,
      code: this.amodCode,
    }
  },

  mounted() {
    this.editor = CodeMirror.fromTextArea(document.getElementById(this.id), {
      lineNumbers: true,
      mode: this.mode,
      theme: 'amod',

      autoCloseBrackets: true,
      autoRefresh: true,
      matchBrackets: true,
      readOnly: this.readOnly,
      styleActiveLine: true,
    })
    this.editor.on('change', this.onCmCodeChange)
  },

  methods: {
    onCmCodeChange() {
      if (this.editor.getValue().length != 0) {
        this.$emit('update:amodCode', this.editor.getValue())
      }
    },

    // Called by the parent to set the code directly
    setCode(code: string) {
      this.editor.setValue(code)
    },
  },
})
</script>

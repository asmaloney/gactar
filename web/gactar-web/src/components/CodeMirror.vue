<template>
  <div>
    <textarea :id="id" v-model="code"></textarea>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import CodeMirror, { Editor } from 'codemirror'

// Add-ons
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/selection/active-line'
import 'codemirror/mode/commonlisp/commonlisp'
import 'codemirror/mode/python/python'

import '../codemirror/amod'

interface Data {
  editor: Editor | null
  id: string
  code: string
}

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

  data(): Data {
    return {
      editor: null,
      id: 'id-' + this.framework,
      code: this.amodCode,
    }
  },

  mounted() {
    const element = document.getElementById(this.id) as HTMLTextAreaElement
    const editor = CodeMirror.fromTextArea(element, {
      lineNumbers: true,
      mode: this.mode,
      theme: 'amod',

      autoCloseBrackets: true,
      autoRefresh: true,
      matchBrackets: true,
      readOnly: this.readOnly,
      styleActiveLine: true,
    })

    editor.on('change', () => {
      this.onCodeChange(editor)
    })

    this.editor = editor
  },

  watch: {
    amodCode(code: string) {
      if (this.editor) {
        this.editor.setValue(code)
      }
    },
  },

  methods: {
    onCodeChange(editor: Editor) {
      if (editor && editor.getValue().length != 0) {
        this.$emit('editorCodeChange', editor.getValue())
      }
    },
  },
})
</script>

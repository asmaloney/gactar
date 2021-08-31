<template>
  <div>
    <textarea id="codemirror" v-model="code"></textarea>
  </div>
</template>

<script>
import CodeMirror from 'codemirror'

// Add-ons
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/edit/matchbrackets'
import 'codemirror/addon/selection/active-line'

require('../codemirror/amod')

export default {
  props: {
    amodCode: {
      type: String,
      default() {
        return '(amod file here)'
      },
    },
  },

  data() {
    return {
      editor: null,
      code: this.amodCode,
    }
  },

  mounted() {
    this.editor = CodeMirror.fromTextArea(
      document.getElementById('codemirror'),
      {
        lineNumbers: true,
        mode: 'amod',
        theme: 'amod',

        autoCloseBrackets: true,
        matchBrackets: true,
        styleActiveLine: true,
      }
    )
    this.editor.on('change', this.onCmCodeChange)
  },

  methods: {
    onCmCodeChange() {
      if (this.editor.getValue().length != 0) {
        this.$emit('update:amodCode', this.editor.getValue())
      }
    },

    // Called by the parent to set the code directly
    setCode(code) {
      this.editor.setValue(code)
    },
  },
}
</script>

<style scoped lang="scss"></style>

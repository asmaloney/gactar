<template>
  <div>
    <textarea id="codemirror" v-model="amodCode"></textarea>
  </div>
</template>

<script>
import CodeMirror from 'codemirror'
import 'codemirror/mode/htmlmixed/htmlmixed.js'
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
    }
  },

  mounted() {
    this.editor = CodeMirror.fromTextArea(
      document.getElementById('codemirror'),
      {
        lineNumbers: true,
        mode: 'amod',
        theme: 'amod',
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
  },
}
</script>

<style scoped lang="scss"></style>

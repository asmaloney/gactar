<template>
  <div>
    <textarea :id="id" v-model="code"></textarea>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import CodeMirror, { Editor } from 'codemirror'

// Add-ons
import 'codemirror/addon/selection/active-line'
import 'codemirror/addon/display/autorefresh'
import 'codemirror/addon/edit/closebrackets'
import 'codemirror/addon/lint/lint'
import { Annotation } from 'codemirror/addon/lint/lint'
import 'codemirror/addon/edit/matchbrackets'

// Modes
import 'codemirror/mode/commonlisp/commonlisp'
import 'codemirror/mode/python/python'

import '../codemirror/amod'

import { Issue, IssueList } from '../api'

interface Data {
  editor: Editor | null
  id: string
  code: string
}

export default defineComponent({
  props: {
    amodCode: {
      type: String,
      default: '(code here)',
    },
    amodIssues: {
      type: Array as PropType<IssueList>,
      required: false,
    },
    mode: {
      type: String,
      required: true,
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
    editorID: {
      type: String,
      required: true,
    },
  },

  data(): Data {
    return {
      editor: null,
      id: `id-${this.editorID}`,
      code: this.amodCode,
    }
  },

  mounted() {
    const element = document.getElementById(this.id) as HTMLTextAreaElement
    const editor = CodeMirror.fromTextArea(element, {
      lineNumbers: true,
      mode: this.mode,
      theme: 'amod',
      lint: {
        lintOnChange: false,
        tooltips: true,
        getAnnotations: this.lint.bind(this),

        options: {},

        // @ts-ignore highlightLines is missing in the typescript interface
        highlightLines: true,
      },
      gutters: ['CodeMirror-lint-markers'],

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

    amodIssues() {
      if (this.editor != null) {
        this.editor.performLint()
      }
    },
  },

  methods: {
    lint(): Annotation[] {
      if (this.amodIssues == null) {
        return []
      }

      var found: Annotation[] = []
      this.amodIssues.forEach((issue: Issue) => {
        if (issue.level == 'info') {
          return
        }

        let from = CodeMirror.Pos(0, 0)
        let to = CodeMirror.Pos(0, 1)

        if (issue.location != null) {
          from.line = issue.location.line - 1
          from.ch = issue.location.columnStart

          to.line = issue.location.line - 1
          to.ch = issue.location.columnEnd
        }

        found.push({
          from: from,
          to: to,
          message: issue.text,
        })
      })

      return found
    },

    onCodeChange(editor: Editor) {
      if (editor && editor.getValue().length != 0) {
        this.$emit('editorCodeChange', editor.getValue())
      }
    },
  },
})
</script>

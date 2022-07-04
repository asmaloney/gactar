<template>
  <b-button type="is-info is-light ml-2" @click="saveCode">
    <span class="fa fa-file-download icon-space" />Save
  </b-button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    code: {
      type: String,
      required: true,
    },
    defaultName: {
      type: String,
      required: true,
    },
    fileExtension: {
      type: String,
      required: true,
    },
  },

  methods: {
    saveCode() {
      // Adapted from: https://stackoverflow.com/a/51315312
      var codeAsBlob = new Blob([this.code], {
        type: 'text/plain;charset=utf-8',
      })

      var downloadLink = document.createElement('a')
      downloadLink.download = this.defaultName + '.' + this.fileExtension
      downloadLink.innerHTML = 'Save File'

      if (window.webkitURL != null) {
        // Chrome allows the link to be clicked without actually adding it to the DOM.
        downloadLink.href = window.webkitURL.createObjectURL(codeAsBlob)
      } else {
        // Firefox requires the link to be added to the DOM before it can be clicked.
        downloadLink.href = window.URL.createObjectURL(codeAsBlob)
        downloadLink.style.display = 'none'
        downloadLink.onclick = (e) => {
          document.body.removeChild(e.target as Node)
        }
        document.body.appendChild(downloadLink)
      }

      downloadLink.click()
    },
  },
})
</script>

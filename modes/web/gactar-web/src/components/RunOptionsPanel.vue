<template>
  <b-collapse
    :open.sync="open"
    class="card"
    animation="slide"
    aria-id="runOptionPanel"
  >
    <template #trigger="props">
      <div
        class="card-header"
        role="button"
        aria-controls="runOptionPanel"
        :aria-expanded="props.open"
      >
        <p class="card-header-title is-size-7">Run Options</p>
        <a class="card-header-icon">
          <b-icon
            pack="fas"
            size="is-tiny"
            type="is-info"
            v-show="!props.open"
            icon="caret-down"
          />
          <b-icon
            pack="fas"
            size="is-tiny"
            type="is-info"
            v-show="props.open"
            icon="caret-up"
          />
        </a>
      </div>
    </template>

    <div class="card-content p-0">
      <div class="tile is-small-gap is-child">
        <div class="tile is-vertical is-parent">
          <div class="tile is-child">
            <div v-if="availableFrameworks.length > 0">
              <b-field label="Select Frameworks" custom-class="is-small">
                <b-checkbox-button
                  v-for="name in availableFrameworks"
                  v-model="runOptions.frameworks"
                  type="is-info"
                  size="is-small"
                  :native-value="name"
                  :key="name"
                  expanded
                  class="ml-1 mr-1"
                  @input="runOptionsChanged"
                >
                  <span>{{ name }}</span>
                </b-checkbox-button>
              </b-field>
            </div>
          </div>

          <div class="tile is-child">
            <b-field label="Logging Level" custom-class="is-small">
              <b-radio-button
                v-for="name in logLevels"
                v-model="runOptions.logLevel"
                type="is-info"
                size="is-small"
                :native-value="name"
                :key="name"
                expanded
                class="ml-1 mr-1"
                @input="runOptionsChanged"
              >
                <span>{{ name }}</span>
              </b-radio-button>

              <b-checkbox
                v-model="runOptions.traceActivations"
                type="is-info"
                expanded
                class="is-small pl-3"
                @input="runOptionsChanged"
              >
                Trace activations
              </b-checkbox>
            </b-field>
          </div>
        </div>
      </div>
    </div>
  </b-collapse>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

import { logLevels, RunOptions } from '../api'

interface Data {
  open: boolean
  runOptions: RunOptions
}

export default defineComponent({
  props: {
    availableFrameworks: {
      type: Array as PropType<string[]>,
      required: true,
    },
    options: {
      type: Object as PropType<RunOptions>,
      required: true,
    },
  },

  data(): Data {
    return { open: true, runOptions: this.options }
  },

  computed: {
    isOpen(): boolean {
      return this.open
    },

    logLevels(): string[] {
      return logLevels
    },
  },

  watch: {
    options: function (newOptions, _) {
      this.runOptions = newOptions
    },
  },

  methods: {
    runOptionsChanged() {
      this.$emit('options-change', this.runOptions)
    },
  },
})
</script>

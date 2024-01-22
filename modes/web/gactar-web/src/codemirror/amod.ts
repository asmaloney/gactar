import CodeMirror, { StringStream } from 'codemirror'

// Implement very basic lexing/parsing of amod files for syntax highlighting

interface State {
  pending: string
  startPattern: boolean // used to get chunk name
  inPattern: boolean // used to check for variables and wildcards

  currentKeywords: object // current list of keywords based on section
}

CodeMirror.defineMode('amod', function () {
  const section_regex = /^~{2}\s*(model|config|init|productions)\s*~{2}/
  const variable_regex = /[?][a-zA-Z0-9_]*/

  const keywords = {
    model: {
      authors: true,
      description: true,
      examples: true,
      name: true,
      nil: true,
    },

    config: {
      chunks: true,
      gactar: true,
      modules: true,
      nil: true,
      true: true,
      false: true, // ;-)
    },

    init: {
      nil: true,
      similar: true,
    },

    productions: {
      and: true,
      any: true,
      buffer_state: true,
      clear: true,
      description: true,
      do: true,
      match: true,
      module_state: true,
      nil: true,
      '!nil': true,
      print: true,
      recall: true,
      set: true,
      stop: true,
      to: true,
      when: true,
      with: true,
    },
  }

  const builtInGlobals = {
    extra_buffers: true,
    goal: true,
    imaginal: true,
    memory: true,
    procedural: true,
    retrieval: true,
  }

  function tokenString(stream: StringStream, state: State): string {
    let current = stream.next()
    while (!stream.eol() && current != state.pending) {
      current = stream.next()
    }

    return 'string'
  }

  function tokenize(stream: StringStream, state: State): string | null {
    const ch = stream.next()

    if (ch == '/') {
      if (stream.eat('/')) {
        stream.skipToEnd()
        return 'comment'
      }
    }

    if (ch === '{' || ch === '}') {
      return 'bracket'
    }

    if (ch == "'" || ch == '"') {
      state.pending = ch
      return tokenString(stream, state)
    }

    if (ch === '[') {
      state.startPattern = true // next id is the chunk name
      return 'bracket'
    }

    if (ch === '*') {
      if (state.inPattern) {
        return 'wildcard'
      }
    }

    if (ch === '?') {
      stream.backUp(1)

      if (stream.match(variable_regex)) {
        return 'variable'
      }
    }

    if (ch === ']') {
      state.inPattern = false
      return 'bracket'
    }

    if (ch === '~') {
      stream.backUp(1)

      const matches = stream.match(section_regex)
      if (matches) {
        const dynamicKey = matches[1] as keyof object
        state.currentKeywords = keywords[dynamicKey]
        return 'section'
      }

      stream.next()
    }

    stream.eatWhile(/[\w-]/)

    const cur = stream.current()

    if (cur in state.currentKeywords) {
      return 'keyword'
    } else if (cur in builtInGlobals) {
      return 'global'
    } else if (state.startPattern) {
      state.startPattern = false
      state.inPattern = true
      return 'chunk-name'
    }

    return null
  }

  return {
    startState: function (): State {
      return {
        pending: '',
        startPattern: false,
        inPattern: false,
        currentKeywords: {},
      }
    },

    token: function (stream: StringStream, state: State): string | null {
      if (stream.eatSpace()) {
        return null
      }
      return tokenize(stream, state)
    },
  }
})

CodeMirror.defineMIME('text/amod', 'gactar')

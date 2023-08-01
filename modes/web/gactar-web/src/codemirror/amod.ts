import CodeMirror, { StringStream } from 'codemirror'

// Implement very basic lexing/parsing of amod files for syntax highlighting

interface State {
  pending: string
  startPattern: boolean // used to get chunk name
  inPattern: boolean // used to check for variables and wildcards
}

CodeMirror.defineMode('amod', function () {
  const section_regex = /^~{2}\s*(model|config|init|productions)\s*~{2}/
  const variable_regex = /[?][a-zA-Z0-9_]*/

  const keywords = {
    and: true,
    authors: true,
    chunks: true,
    clear: true,
    description: true,
    do: true,
    examples: true,
    gactar: true,
    match: true,
    module: true,
    modules: true,
    name: true,
    print: true,
    recall: true,
    set: true,
    similar: true,
    state: true,
    stop: true,
    to: true,
    when: true,
    true: true,
    false: true, // ;-)
    nil: true,
    '!nil': true,
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

      if (stream.match(section_regex)) {
        return 'section'
      }

      stream.next()
    }

    stream.eatWhile(/[\w-]/)

    const cur = stream.current()
    if (cur in keywords) {
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

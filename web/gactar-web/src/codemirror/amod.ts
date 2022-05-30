import CodeMirror, { StringStream } from 'codemirror'

// Implement very basic lexing/parsing of amod files for syntax highlighting

interface State {
  pending: string
  startPattern: boolean
}

CodeMirror.defineMode('amod', function () {
  const section_regex = /^={2}(model|config|init|productions)={2}/
  const variable_regex = /[?][a-zA-Z0-9_]*/

  const keywords = {
    authors: true,
    chunks: true,
    clear: true,
    description: true,
    do: true,
    examples: true,
    gactar: true,
    match: true,
    modules: true,
    name: true,
    print: true,
    recall: true,
    set: true,
    to: true,
    write: true,
  }

  const builtIns = {
    goal: true,
    imaginal: true,
    memory: true,
    procedural: true,
    nil: true,
    '!nil': true,
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

    if (ch === '[') {
      state.startPattern = true // next id is the chunk name
      return 'bracket'
    }

    if (ch === ']') {
      return 'bracket'
    }

    if (ch === '{' || ch === '}') {
      return 'bracket'
    }

    if (ch == "'" || ch == '"') {
      state.pending = ch
      return tokenString(stream, state)
    }

    if (ch === '?') {
      stream.backUp(1)

      if (stream.match(variable_regex)) {
        return 'variable'
      }
    }

    if (ch === '=') {
      stream.backUp(1)

      if (stream.match(section_regex)) {
        return 'header'
      }

      stream.next()
    }

    stream.eatWhile(/[\w-]/)

    const cur = stream.current()
    if (cur in keywords) {
      return 'keyword'
    } else if (cur in builtIns) {
      return 'built-in'
    } else if (state.startPattern) {
      state.startPattern = false
      return 'chunk-name'
    }

    return null
  }

  return {
    startState: function (): State {
      return { pending: '', startPattern: false }
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

import CodeMirror from 'codemirror'

// Implement very basic lexing/parsing of amod files for syntax highlighting

CodeMirror.defineMode('amod', function () {
  const section_regex = /^={2}(model|config|init|productions)={2}/
  const variable_regex = /[?][a-zA-Z0-9_]*/

  const keywords = {
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
    nil: true,
    '!nil': true,
    retrieval: true,
  }

  function tokenString(stream, state) {
    var current = stream.next()
    while (!stream.eol() && current != state.pending) {
      current = stream.next()
    }

    return 'string'
  }

  function tokenize(stream, state) {
    var ch = stream.next()

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

    var cur = stream.current()
    if (cur in keywords) {
      return 'keyword'
    } else if (cur in builtIns) {
      return 'built-in'
    } else if (state.startPattern) {
      state.startPattern = false
      return 'chunk-name'
    }
  }

  return {
    startState: function () {
      var state = {}

      state.pending = false
      state.startPattern = false

      return state
    },
    token: function (stream, state) {
      if (stream.eatSpace()) return null
      return tokenize(stream, state)
    },
  }
})

CodeMirror.defineMIME('text/amod', 'gactar')

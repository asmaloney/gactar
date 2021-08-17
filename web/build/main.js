function runModel() {
  const url = '/run'
  const data = {
    amod: document.getElementById('amod').value,
    run: document.getElementById('run').value,
  }
  const params = {
    headers: { 'content-type': 'application/json; charset=UTF-8' },
    method: 'POST',
    mode: 'cors',
    body: JSON.stringify(data),
  }

  fetch(url, params)
    .then(function (response) {
      if (response.ok) {
        return response.json()
      }

      throw new Error('Unreachable: ' + response.statusText)
    })
    .then((data) => {
      if (data.error) {
        setError(data.error)
        return
      }

      setResults(data.results)
    })
    .catch((error) => {
      clearResults()
      setError(error)
    })
}

function clearResults() {
  document.getElementById('results').textContent = ''
}

function setResults(results) {
  let text = ''
  for (const [key, value] of Object.entries(results)) {
    text += key + '\n' + '---\n'
    text += value
    text += '\n\n'
  }

  document.getElementById('results').textContent = text
}

function setError(error) {
  document.getElementById('results').textContent += error
}

function setAMOD(text) {
  document.getElementById('amod').textContent = text
}

function setRun(text) {
  document.getElementById('run').textContent = text
}

function loadExampleAMOD(example) {
  const url = '/examples/' + example
  const params = {
    headers: { 'content-type': 'text/plain; charset=UTF-8' },
    method: 'GET',
    mode: 'cors',
  }

  fetch(url, params)
    .then((response) => {
      if (response.ok) {
        return response.text()
      }

      throw new Error('Unreachable: ' + response.statusText)
    })
    .then((text) => {
      setAMOD(text)
    })
    .catch((error) => {
      setResults(error)
    })
}

function addExamples(example_list) {
  var select = document.getElementById('examples')

  for (const name of example_list) {
    var option = document.createElement('option')
    option.text = name
    select.add(option)
  }
}

function handleExampleChange() {
  var selectBox = document.getElementById('examples')

  loadExampleAMOD(selectBox.value)
}

function loadExampleList() {
  const url = '/examples/list'
  const params = {
    headers: { 'content-type': 'text/plain; charset=UTF-8' },
    method: 'GET',
    mode: 'cors',
  }

  fetch(url, params)
    .then((res) => res.json())
    .then((res) => {
      addExamples(res.example_list)
      handleExampleChange()
    })
    .catch((error) => {
      setResults(error)
    })
}

window.addEventListener('load', function () {
  loadExampleList()
})

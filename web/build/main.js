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
    .then(function (data) {
      if (data.error) {
        setError(data.error)
        return
      }

      setResults(data.results)
    })
    .catch(function (error) {
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

function loadExampleAMOD() {
  const url = '/examples/count.amod'
  const params = {
    headers: { 'content-type': 'text/plain; charset=UTF-8' },
    method: 'GET',
    mode: 'cors',
  }

  fetch(url, params)
    .then(function (response) {
      if (response.ok) {
        return response.text()
      }

      throw new Error('Unreachable: ' + response.statusText)
    })
    .then(function (text) {
      setAMOD(text)
      setRun('countFrom 2 5 starting')
    })
    .catch(function (error) {
      setResults(error)
    })
}

window.addEventListener('load', function () {
  loadExampleAMOD()
})

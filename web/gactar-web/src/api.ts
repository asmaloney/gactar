import axios from 'axios'

const http = axios.create({
  headers: { 'Content-Type': 'application/json' },
})

// version
export interface Version {
  version: string
}

async function getVersion(): Promise<Version> {
  const response = await http.get<Version>('/api/version')
  return response.data
}

// frameworks
export interface FrameworkInfo {
  name: string
  language: string

  fileExtension: string

  executableName: string

  pythonRequiredPackages?: string[]
}

export type FrameworkInfoList = FrameworkInfo[]

export interface FrameworkInfoResponse {
  frameworks: FrameworkInfoList
}

async function getFrameworks(): Promise<FrameworkInfoList> {
  const response = await http.get<FrameworkInfoResponse>('/api/frameworks')
  return response.data.frameworks
}

// run
export interface Result {
  modelName: string
  filePath: string
  code: string
  output: string
}

export type ResultMap = { [key: string]: Result }

export interface Results {
  results: ResultMap
}

export interface RunError {
  error: string
}

export type RunResult = Results | RunError

async function run(
  amod: string,
  goal: string,
  frameworks: string[]
): Promise<RunResult> {
  const response = await http.post<RunResult>('/api/run', {
    amod: amod,
    goal: goal,
    frameworks: frameworks,
  })
  return response.data
}

// examples
export interface ExampleList {
  example_list: string[]
}

async function getExampleList(): Promise<ExampleList> {
  const response = await http.get<ExampleList>('/api/examples/list')
  return response.data
}

async function getExample(name: string): Promise<string> {
  const response = await http.get<string>('/api/examples/' + name)
  return response.data
}

export default {
  getExample,
  getExampleList,
  getFrameworks,
  getVersion,
  run,
}

import axios from 'axios'

const http = axios.create({
  headers: { 'Content-Type': 'application/json' },
})

// version
export type Version = string

export interface VersionResponse {
  // The current version tag when gactar was built.
  version: string
}

async function getVersion(): Promise<Version> {
  const response = await http.get<VersionResponse>('/api/version')
  return response.data.version
}

// frameworks
export interface FrameworkInfo {
  // Name (id) of the framework.
  name: string

  // Language the framework uses.
  language: string

  // File extension of the intermediate file.
  fileExtension: string

  // Name of the executable that was run.
  executableName: string

  // (Python only) List of packages this framework requires.
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
export interface RunParams {
  // The text of the amod to run.
  amod: string

  // The starting goal.
  goal: string

  // An optional list of frameworks ("all" if not set).
  frameworks?: string[]
}

export interface Result {
  // Name of the model (from the amod text).
  modelName: string

  // Intermediate code file (full path).
  filePath: string

  // Code which was run.
  code?: string

  // Output of run (stdout + stderr).
  output: string
}

export type ResultMap = { [key: string]: Result }

export interface Results {
  results: ResultMap
}

export interface Issue {
  level: string
  text: string
  line: number
  columnStart: number
  columnEnd: number
}

export type IssueList = Issue[]

export interface RunIssues {
  issues: IssueList
}

export type RunResult = Results | RunIssues

async function run(params: RunParams): Promise<RunResult> {
  const response = await http.post<RunResult>('/api/run', params)
  return response.data
}

// examples
// List of example names which are built into the webserver.
export type ExampleList = string[]

export interface ExampleListResponse {
  exampleList: ExampleList
}

async function getExampleList(): Promise<ExampleList> {
  const response = await http.get<ExampleListResponse>('/api/examples/list')
  return response.data.exampleList
}

async function getExample(name: string): Promise<string> {
  const response = await http.get<string>('/api/examples/' + name)
  return response.data
}

// sessions
export interface Session {
  sessionID: number
}

export interface SessionRunParams {
  // The id of the session.
  sessionID: number

  // The ID of the model to run.
  modelID: number

  // The initial contents of the buffers.
  buffers: string

  // An optional list of frameworks ("all" if not set).
  frameworks?: string[]

  // Whether to include the generated code as part of the response.
  includeCode: boolean
}

export interface SessionRunResult extends Result {
  // The id of the session.
  sessionID: number

  // The ID of the model which was run.
  modelID: number
}

export type SessionResultMap = { [key: string]: SessionRunResult }

export interface SessionRunResults {
  results: SessionResultMap
}

async function sessionBegin(): Promise<Session> {
  const response = await http.get<Session>('/api/session/begin')
  return response.data
}

async function sessionEnd(session: Session): Promise<void> {
  await http.put<Session>('/api/session/end', session)
  return
}

async function sessionRun(
  params: SessionRunParams
): Promise<SessionRunResults> {
  const response = await http.post<SessionRunResults>(
    '/api/session/runModel',
    params
  )
  return response.data
}

// models
export interface ModelParams {
  // The amod code to load.
  amod: string

  // The id of the session to load this model in.
  sessionID: number
}

export interface ModelLoadResult {
  // The ID of this model to use in other calls to the API.
  modelID: number

  // The name of the model (comes from the amod code).
  modelName: string

  // The id of the session.
  sessionID: number
}

async function modelLoad(params: ModelParams): Promise<ModelLoadResult> {
  const response = await http.put<ModelLoadResult>('/api/session/begin', params)
  return response.data
}

export default {
  getExample,
  getExampleList,
  getFrameworks,
  getVersion,
  modelLoad,
  run,
  sessionBegin,
  sessionEnd,
  sessionRun,
}

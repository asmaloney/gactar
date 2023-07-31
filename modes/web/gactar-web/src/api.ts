import axios, { AxiosInstance } from 'axios'

let gactarHTTP: AxiosInstance

function init(port: number) {
  gactarHTTP = axios.create({
    headers: { 'Content-Type': 'application/json' },
    baseURL: `http://localhost:${port}`,
  })
}

// version
export type Version = string

export interface VersionResponse {
  // The current version tag when gactar was built.
  version: Version
}

async function getVersion(): Promise<Version> {
  const response = await gactarHTTP.get<VersionResponse>('/api/version')
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
  const response = await gactarHTTP.get<FrameworkInfoResponse>(
    '/api/frameworks'
  )
  return response.data.frameworks
}

// run
export interface RunOptions {
  logLevel: string
  traceActivations: boolean
  randomSeed?: number
}

export interface RunParams {
  // The text of the amod to run.
  amod: string

  // The starting goal.
  goal: string

  // An optional list of frameworks ("all" if not set).
  frameworks?: string[]

  // optional options!
  options?: RunOptions
}

// Location of an issue in the source code.
export interface Location {
  line: number
  columnStart: number
  columnEnd: number
}

export interface Issue {
  // Severity of the issue.
  level: string

  // Text of the issue.
  text: string

  // Location in the code (optional)
  location?: Location
}

export type IssueList = Issue[]

export interface FrameworkResult {
  // Name of the model (from the amod text).
  modelName: string

  // Any issues specific to a framework.
  issues?: IssueList

  // Intermediate code file (full path).
  filePath?: string

  // Code which was run.
  code?: string

  // Output of run (stdout + stderr).
  output?: string
}

export type FrameworkResultMap = { [key: string]: FrameworkResult }

export interface RunResult {
  issues?: IssueList
  results?: FrameworkResultMap
}

async function run(params: RunParams): Promise<RunResult> {
  const response = await gactarHTTP.post<RunResult>('/api/run', params)
  return response.data
}

// examples
// List of example names which are built into the webserver.
export type ExampleList = string[]

export interface ExampleListResponse {
  exampleList: ExampleList
}

async function getExampleList(): Promise<ExampleList> {
  const response = await gactarHTTP.get<ExampleListResponse>(
    '/api/examples/list'
  )
  return response.data.exampleList
}

async function getExample(name: string): Promise<string> {
  const response = await gactarHTTP.get<string>('/api/examples/' + name)
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

export interface SessionRunResult extends FrameworkResult {
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
  const response = await gactarHTTP.get<Session>('/api/session/begin')
  return response.data
}

async function sessionEnd(session: Session): Promise<void> {
  await gactarHTTP.put<Session>('/api/session/end', session)
  return
}

async function sessionRun(
  params: SessionRunParams
): Promise<SessionRunResults> {
  const response = await gactarHTTP.post<SessionRunResults>(
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
  const response = await gactarHTTP.put<ModelLoadResult>(
    '/api/session/begin',
    params
  )
  return response.data
}

export default {
  getExample,
  getExampleList,
  getFrameworks,
  getVersion,
  init,
  modelLoad,
  run,
  sessionBegin,
  sessionEnd,
  sessionRun,
}

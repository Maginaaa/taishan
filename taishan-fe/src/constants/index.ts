/**
 * 请求成功状态码
 */
export const SUCCESS_CODE = 0

/**
 * 请求contentType
 */
export const CONTENT_TYPE: AxiosContentType = 'application/json'

/**
 * 请求超时时间
 */
export const REQUEST_TIMEOUT = 60000

/**
 * 不重定向白名单
 */
export const NO_REDIRECT_WHITE_LIST = ['/login']

/**
 * 不重置路由白名单
 */
export const NO_RESET_WHITE_LIST = ['Redirect', 'Login', 'NoFind', 'Root']

/**
 * 表格默认过滤列设置字段
 */
export const DEFAULT_FILTER_COLUMN = ['expand', 'selection']

/**
 * 是否根据headers->content-type自动转换数据格式
 */
export const TRANSFORM_REQUEST_DATA = true

export const CaseType = {
  HttpCase: 1,
  LogicControl: 11,
  SqlCase: 21
}

export const OperationType = {
  CREATE: 0,
  EDIT: 1,
  COPY: 2,
  DELETE: 3,
  DISABLED: 11,
  ENABLED: 12,
  DEBUG: 21,
  RUN: 22
}

export const PressType = {
  ConcurrentModel: 0,
  StepModel: 1,
  RPSModel: 2
}

export const DataSnapshotType = {
  SlsLog: 1,
  Trace: 2,
  LogSize: 3
}

export const OperationTypeMap = {
  [OperationType.CREATE]: '创建',
  [OperationType.EDIT]: '编辑',
  [OperationType.COPY]: '复制',
  [OperationType.DELETE]: '删除',
  [OperationType.DISABLED]: '禁用',
  [OperationType.ENABLED]: '启用',
  [OperationType.DEBUG]: '调试',
  [OperationType.RUN]: '执行'
}

export const SourceName = {
  Plan: 'Plan',
  SceneCase: 'SceneCase',
  Scene: 'Scene',
  File: 'File',
  NormalizedSwitch: 'NormalizedSwitch'
}

export const SourceNameMap = {
  [SourceName.Plan]: '计划',
  [SourceName.SceneCase]: '用例',
  [SourceName.Scene]: '场景',
  [SourceName.File]: '文件',
  [SourceName.NormalizedSwitch]: '开关'
}

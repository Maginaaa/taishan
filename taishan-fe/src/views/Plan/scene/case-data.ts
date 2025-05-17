export interface SceneInfo {
  plan_id?: number
  scene_id?: number
  scene_name: string
  scene_type: number
  export_data_info?: object | null
  disabled: boolean
  rate?: number
  concurrency?: number
  case_tree: CaseData[]
  variable_list: Variable[]
}

export interface CaseData {
  case_id?: number
  title: string
  parent_id: number
  type: number
  scene_id: number
  sort: number
  disabled: boolean
  extend: HttpCaseExtend
}

export interface HttpCaseExtend {
  assert_form: AssertForm[]
  body: {
    body_type: string
    body_value: string
  }
  case_name: string
  headers_form: HeadersForm[]
  method_type: string
  params_form: HeadersForm[]
  rally_point: RallyPoint
  url: string
  use_scene_variable: boolean
  variable_form: VariableForm[]
  waiting_config: WaitingConfig
  overtime_config: OvertimeConfig
}

export interface AssertForm {
  checking_rule: number
  desc: string
  enable: boolean
  expect_value: string
  extract_express: string
  extract_type: number
}

export interface HeadersForm {
  desc?: string
  enable: boolean
  key: string
  value: string
}

export interface VariableForm {
  enable: boolean
  extract_type: number
  key: string
  value: string
}

export interface RallyPoint {
  concurrency: number
  enable: boolean
  timeout_period: number
}

export interface WaitingConfig {
  post_waiting_switch: boolean
  post_waiting_time: number
  pre_waiting_switch: boolean
  pre_waiting_time: number
}

export interface OvertimeConfig {
  enable: boolean
  timeout_period: number
}

export interface Variable {
  variable_id: number
  variable_name: string
  variable_val: string
  remark: string
}

export interface PlanDebugResult {
  scene_id: number
  scene_name: string
  passed: boolean
  total_requests: number
  success_count: number
  start_time: string
  end_time: string
  total_time: number
  case_result_list: HttpCaseResponse[]
}

export interface HttpCaseResponse {
  case_id: number
  case_name: string
  url: string
  method_type: string
  request_header: HeadersForm[]
  request_content: string
  status_code: number
  response_content: string
  response_header: HeadersForm[]
  assert_res: []
  variable_res: []
  assert_success: boolean
  extract_all_success: boolean
  response_size: {
    body_size: number
    header_size: number
    total_size: number
  }
  response_time: number
  send_bytes: number
  receiver_bytes: number
  request_success: boolean
  err: string
  start_time: number
  end_time: number
}

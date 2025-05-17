export interface ReportSearchReq {
  plan_id: number | null
  create_user_id: number | null
  page: number
  page_size: number
  start_time: string | null
  end_time: string | null
}

export const baseReportSearch: ReportSearchReq = {
  plan_id: null,
  create_user_id: null,
  page: 1,
  page_size: 10,
  start_time: null,
  end_time: null
}

export interface ReportData {
  end: boolean
  stage_total_request_num: number
  stage_total_response_time: number
  stage_avg_response_time: number
  actual_concurrency: number
  stage_send_bytes: number
  stage_received_bytes: number
  total_response_time: number
  total_request_num: number
  total_success_num: number
  api_success_rate: number
  plan_tps: number
}

export interface ReportDetail {
  report_id: number
  report_name: string
  status: boolean
  plan_id: number
  plan_name: string
  press_type: number
  duration: number
  actual_duration: number
  concurrency: number
  machine_num: number
  remark: string
  total_count: number
  success_rate: number
  start_time: string
  end_time: string
  total_send_bytes: number
  total_receive_bytes: number
  scene_press_result_list: SceneData[]
  create_user_id: number
  create_user_name: string
  update_user_name: string
  update_time: string
}

export interface SceneData {
  scene_id: number
  scene_name: string
  total_count: number
  success_count: number
  total_response_time: number
  avg_response_time: number
  min_response_time: number
  max_response_time: number
  total_send_bytes: number
  total_receive_bytes: number
  tps: number
  vum: number
  data: CaseData[]
}

export interface CaseData {
  case_id: number
  case_name: string
  rally_point: number
  total_request_num: number
  total_request_time: number
  success_num: number
  error_num: number
  avg_request_time: number
  max_request_time: number
  min_request_time: number
  custom_request_time_line: number
  fifty_request_time_line: number
  ninety_request_time_line: number
  ninety_five_request_time_line: number
  ninety_nine_request_time_line: number
  custom_request_time_line_value: number
  fifty_request_time_line_value: number
  ninety_request_time_line_value: number
  ninety_five_request_time_line_value: number
  ninety_nine_request_time_line_value: number
  total_send_bytes: number
  total_received_bytes: number
  start_time: number
  end_time: number
  status_code_counter: number[]
  rps: number
  s_rps: number
  e_rps: number
  actual_concurrency: number
  stage_start_time: number
  stage_end_time: number
  stage_total_request_num: number
  stage_total_request_time: number
  stage_avg_response_time: number
  stage_success_num: number
  stage_error_num: number
  stage_send_bytes: number
  stage_received_bytes: number
  stage_rps: number
  stage_s_rps: number
  stage_e_rps: number
}

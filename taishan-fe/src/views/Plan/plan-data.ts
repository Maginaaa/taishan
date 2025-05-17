import { EChartsOption } from 'echarts'

export interface TreeNode {
  label: string
  value: number
  children?: any[]
  leaf?: Boolean
}

export interface SimplePlan {
  plan_name: string
  plan_id: number
}

export interface PlanInfo {
  plan_name: string
  create_user_name: string
  update_user_name: string
  last_press_time: string
  create_time: string
  update_time: string
  press_count: number
  remark: string
  engine_count: number
  press_info: PressInfo
  sampling_info: SamplingInfo
  global_variable: ParamsForm[]
  default_header: ParamsForm[]
  task_info: TaskInfo
  capacity_switch: boolean
  server_info: string[][]
}

export interface PressInfo {
  press_type: number
  start_concurrency: number
  step_size: number
  step_duration: number
  concurrency: number
  duration: number
  scene_list: PressSceneInfo[]
}

export interface PressSceneInfo {
  scene_type: number
  scene_id: number
  scene_name: string
  rate: number
  concurrency: number
}

export interface SamplingInfo {
  sampling_type: number
  sampling_rate: number
}

export interface ParamsForm {
  enable: boolean
  key: string
  value: string
  desc: string
}

export interface TaskInfo {
  task_id: number
  enable: boolean
  cron: string
}

export const pressGraphOptions: EChartsOption = {
  title: {
    text: '压力预估图',
    left: 'center'
  },
  grid: {
    left: 50,
    right: 60
  },
  tooltip: {
    confine: true,
    trigger: 'axis',
    formatter(params) {
      // 使用 map 函数构建每个场景的数据字符串
      return params
        .map((param) => {
          return (
            param.marker + `${param.seriesName} 并发数：${param.data[1]} 时间：${param.data[0]}秒`
          )
        })
        .join('<br/>')
    }
  },
  xAxis: {
    type: 'value',
    name: '时间(秒)',
    minInterval: 1
  },
  yAxis: {
    type: 'value',
    name: '并发数'
  },
  series: []
}

export const emptyPlanInfo: PlanInfo = {
  plan_name: '拼命加载中…',
  create_user_name: '',
  update_user_name: '',
  last_press_time: '',
  create_time: '',
  update_time: '',
  press_count: 0,
  remark: '',
  engine_count: 0,
  press_info: {
    press_type: 0,
    start_concurrency: 0,
    step_size: 0,
    step_duration: 0,
    concurrency: 0,
    duration: 0,
    scene_list: []
  },
  sampling_info: {
    sampling_type: 0,
    sampling_rate: 0
  },
  global_variable: [],
  default_header: [],
  task_info: {
    task_id: 0,
    enable: false,
    cron: '0 0 1 * * ?'
  },
  capacity_switch: false,
  server_info: []
}

<script setup lang="ts">
import { onActivated, onDeactivated, reactive, ref, unref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useIcon } from '@/hooks/web/useIcon'
import { isNumber } from '@/utils/is'
import { useAuthStore } from '@/store/modules/auth'

import { ContentWrap } from '@/components/ContentWrap'
import { ElMessage } from 'element-plus'

import type { ReportDetail } from '@/views/Report/report-data'

import ReportService from '@/api/scene/report'
import SceneService from '@/api/scene/scene'
import SceneReportService from '@/api/report'
import K8sService from '@/api/scene/k8s'

import DownSceneDetail from '@/views/Report/detail/DownSceneDetail.vue'
import SamplingData from '@/views/Report/sampling/SamplingData.vue'
import PlanChart from '@/views/Report/detail/PlanChart.vue'
import MemChart from '@/views/Machine/chart/MemChart.vue'
import CpuChart from '@/views/Machine/chart/CpuChart.vue'
import MachineService from '@/api/scene/machine'
import ApiDistribution from '@/views/Report/detail/ApiDistribution.vue'
import TracesList from '@/views/Report/detail/TracesList.vue'
import UrlDistribution from '@/views/Report/detail/UrlDistribution.vue'

const router = useRouter()
const authStore = useAuthStore()

defineOptions({
  name: 'ReportDetail'
})

const infoFillIcon = useIcon({ icon: 'ep:info-filled', color: 'red' })

const reportService = new ReportService()
const sceneService = new SceneService()
const sceneReportService = new SceneReportService()
const machineService = new MachineService()
const k8sService = new K8sService()

// 通用变量
const samplingDataRef = ref()
let press_status = ref(0)

let report_id = 0
let loading = ref<boolean>(false)
let report_detail = reactive<ReportDetail | any>({})
let report_data = reactive<any>({})

const getReportId = () => {
  report_id = Number(router.currentRoute.value.params.id)
}

const getReportStatus = async () => {
  loading.value = true
  await getReportDetail()
  await getPlanScenes(report_detail.plan_id)
  if (report_detail.status) {
    await getReportData()
    if (Object.keys(report_data).length > 0) {
      press_status.value = report_data.end ? 3 : 2
    } else {
      press_status.value = 1
    }
  } else {
    press_status.value = 4
  }
  loading.value = false
}

// prePressing数据
let progress = ref(0)
const getProgress = () => {
  const interval = setInterval(() => {
    const increase = Math.floor(Math.random() * 10) + 15
    const _val = progress.value + increase
    if (_val > 99) {
      clearInterval(interval)
      progress.value = 99
      return
    }
    progress.value = _val
  }, 1000)
}

// pressing数据
const editIcon = useIcon({ icon: 'eva:edit-2-outline' })

let data_interval: any = null
let detail_interval: any = null
let engine_interval: any = null
let current_duration = ref('00:00:00')
let time_percent = ref<number>(0.0)
let pre_scene_end = ref<boolean>(false)
let total_duration = ref('')
let stop_press_btn_loading = ref(false)
let change_param = reactive({
  scene_id: 0,
  concurrency: 0,
  report_id: 0
})
let chang_con_dialog_visible = ref(false)
let active_tab = ref('case')
let target_rps = ref(0)
let pre_target_rps = ref(0)
let prefix_scene_data_current_count = ref(0)
let prefix_scene_data_count = ref(0)

const calcSchedule = () => {
  if (report_detail.press_type === 4) {
    return
  }
  const schedule_interval = setInterval(() => {
    if (report_data.prefix_scene_end_time !== '' && !pre_scene_end.value) {
      pre_scene_end.value = true
    }
    if (pre_scene_end.value) {
      const startTime = new Date(report_data.prefix_scene_end_time)
      const currentTime = new Date()
      //@ts-ignore
      const secondDiff = Math.floor(Math.abs(currentTime - startTime) / 1000)
      const seconds = String(Math.floor(secondDiff % 60)).padStart(2, '0')
      const minutes = String(Math.floor((secondDiff / 60) % 60)).padStart(2, '0')
      const hours = String(Math.floor(secondDiff / 3600)).padStart(2, '0')
      current_duration.value = `${hours}:${minutes}:${seconds}`
      const percent = parseFloat(((secondDiff * 100) / (report_detail.duration * 60)).toFixed(1))
      if (percent >= 100) {
        if (press_status.value === 2) {
          press_status.value = 3
        }
        clearInterval(schedule_interval)
      }
      time_percent.value = percent > 100 ? 100 : percent
    } else {
      prefix_scene_data_current_count.value = report_data.scenes[0].cases[0].total_request_num
    }
  }, 1000)
}

const stopTest = () => {
  stop_press_btn_loading.value = true
  reportService.stopPress(report_id)
  if (press_status.value === 2) {
    press_status.value = 3
  }
}

function editPlanRps() {
  target_rps.value = pre_target_rps.value
  sceneReportService
    .reportRpsModify({
      report_id: report_id,
      new_rps: target_rps.value
    })
    .then(() => {
      getReportTargetRps()
      ElMessage.success('修改计划RPS成功')
    })
    .finally(() => {})
}

const getSamplingData = (scene) => {
  const param = {
    report_id: report_id,
    scene_id: scene.scene_id
  }
  let _case_list: any = []
  scene.cases.forEach((c) => {
    _case_list.push({ id: c.case_id, value: caseNameMap.value[c.case_id] })
  })
  samplingDataRef.value.drawerOpen(param, _case_list)
}

const preChangeConcurrency = (item) => {
  chang_con_dialog_visible.value = true
  change_param.scene_id = item.scene_id
  change_param.concurrency = item.concurrency
}

const releaseScene = (scene_id) => {
  reportService.releasePreScene({
    report_id: report_id,
    scene_id: scene_id
  })
}

const changeConcurrency = () => {
  chang_con_dialog_visible.value = false
  change_param.report_id = report_id
  if (!isNumber(change_param.concurrency)) {
    ElMessage.warning('请输入正整数')
  }
  change_param.concurrency = Number(change_param.concurrency)
  reportService.updatePress(change_param)
}

let case_graph_dialog_visible = ref(false)
let case_graph_title = ref('')
let case_graph_id = ref<number>(0)
let case_data = reactive<any>({})
let engine_data = ref({
  cpu: [],
  mem: [],
  network: []
})

const getCaseGraph = async (row) => {
  case_graph_dialog_visible.value = true
  case_graph_title.value = row.case_name
  case_graph_id.value = row.case_id
  let res = await reportService.getReportCaseData({
    report_id: report_id,
    case_id: row.case_id
  })
  case_data = Object.assign(case_data, res)
}

const changeChartType = (tabName: string) => {
  if (tabName === 'engine') {
    startGetEngineChartData()
  } else {
    stopGetEngineChartData()
  }
  if (tabName === 'server_info') {
    getPlanServerInfo()
  }
}
const startGetEngineChartData = () => {
  if (!engine_interval) {
    engine_interval = setInterval(getEngineChartData, 3000)
  }
}

const stopGetEngineChartData = () => {
  if (engine_interval) {
    clearInterval(engine_interval)
    engine_interval = null
  }
}

const serverInfo = ref<any>([])
const getPlanServerInfo = async () => {
  const res = await k8sService.getPlanServerInfo(report_detail.plan_id)
  serverInfo.value = [...res]
}
const replicasChange = (val) => {
  k8sService.updateServerReplicas(val)
}

const getEngineChartData = async () => {
  const param = {
    ip_list: unref(report_detail.engine_list)
  }
  engine_data.value = await machineService.getMachineUseInfo(param)
}

// 通用常量&工具函数
const press_type_mp = {
  0: '并发模式',
  1: '阶梯式加压',
  2: 'RPS模式',
  3: 'RPS比例模式',
  4: '固定次数'
}
const sceneNameMap = ref<any>({})
const caseNameMap = ref<any>({})
const getPlanScenes = async (plan_id) => {
  sceneNameMap.value = {}
  caseNameMap.value = {}
  const scene_list = (await sceneService.getSceneList(plan_id)) || []
  for (let scene of scene_list) {
    if (scene.scene_type === 1) {
      prefix_scene_data_count.value = scene.export_data_info.export_times
    }
    sceneNameMap.value[scene.scene_id] = scene.scene_name
    if (scene.case_tree) {
      buildCaseNameMap(scene.case_tree)
    }
  }
}

const getReportDetail = async () => {
  const detail = await reportService.getReportDetail(report_id)
  report_detail = Object.assign(report_detail, detail)
  if (press_status.value === 3 && !detail.status) {
    press_status.value = 4
  }
  if (!detail.status) {
    stopGetReportDetail()
  }
}

const getReportTargetRps = async () => {
  const rps = await reportService.getReportTargetRps(report_id)
  target_rps.value = rps
}

const getReportData = async () => {
  let res = await reportService.getReportData(report_id)
  report_data = Object.assign(report_data, res)
  if (press_status.value === 1 && report_data) {
    press_status.value = 2
  }
  if (press_status.value === 2 && report_data.end) {
    press_status.value = 3
  }
  if (report_data.end) {
    stopGetReportData()
  }
}

function buildCaseNameMap(case_tree) {
  if (!case_tree) return
  for (let cs of case_tree) {
    caseNameMap.value[cs.case_id] = cs.title
    if (cs.children) {
      buildCaseNameMap(cs.children)
    }
  }
}
const convertByte = (val) => {
  if (!val) {
    return 0
  }
  if (val >= 1024 && val < 1024 ** 2) {
    return (val / 1024).toFixed(2) + 'MB'
  } else if (val >= 1024 ** 2 && val < 1024 ** 3) {
    return (val / 1024 ** 2).toFixed(2) + 'GB'
  } else if (val >= 1024 ** 3) {
    return (val / 1024 ** 3).toFixed(2) + 'TB'
  } else {
    return val.toFixed(2) + 'KB'
  }
}
const getTotalDuration = () => {
  const duration = report_detail.duration
  const hours = String(Math.floor(duration / 60)).padStart(2, '0')
  const minutes = String(duration % 60).padStart(2, '0')
  return `${hours}:${minutes}:00`
}

const convertReqNum = (val) => {
  if (!val) {
    return 0
  }
  if (val < 10000) {
    return val
  } else {
    return (val / 10000).toFixed(1) + 'w'
  }
}

const jumpPlanDetail = (planId) => {
  router.push({
    name: 'PlanDetail',
    params: {
      id: planId
    }
  })
}

const startGetReportData = () => {
  if (!data_interval) {
    data_interval = setInterval(getReportData, 5000)
  }
}

const stopGetReportData = () => {
  if (data_interval) {
    clearInterval(data_interval)
    data_interval = null
  }
  stopGetEngineChartData()
}

const startGetReportDetail = () => {
  if (!detail_interval) {
    detail_interval = setInterval(getReportDetail, 3000)
  }
}

const stopGetReportDetail = () => {
  if (detail_interval) {
    clearInterval(detail_interval)
    detail_interval = null
  }
}

const tracesListRef = ref()
const getTraces = async (row) => {
  tracesListRef.value.drawerOpen(row.case_id, report_id)
}

watch(press_status, (n, o) => {
  if (o >= n) {
    ElMessage.error('状态变更错误')
    return
  }
  switch (n) {
    case 1:
      getProgress()
      startGetReportData()
      return
    case 2:
      startGetReportData()
      calcSchedule()
      if (report_detail.press_type === 2 || report_detail.press_type === 3) {
        getReportTargetRps()
      }
      total_duration.value = getTotalDuration()
      return
    case 3:
      startGetReportDetail()
      stopGetReportData()
      return
    case 4:
      stopGetReportDetail()
      loading.value = true
      getReportData()
      loading.value = false
      return
    default:
      return null
  }
})

getReportId()
getReportStatus()

onActivated(() => {
  if (press_status.value === 2) {
    active_tab.value = 'case'
    startGetReportDetail()
  }
  if (press_status.value === 3) {
    getReportDetail()
  }
})

onDeactivated(() => {
  stopGetReportDetail()
  stopGetEngineChartData()
})
</script>

<template>
  <el-skeleton :loading="loading" :rows="12" animated>
    <div v-if="press_status === 1">
      <content-wrap>
        <el-descriptions :title="report_detail.plan_name">
          <el-descriptions-item label="持续时间">
            {{ report_detail.duration }} 分钟
          </el-descriptions-item>
          <el-descriptions-item label="压测模式">
            <el-tag size="small">
              {{ press_type_mp[report_detail.press_type] }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="预估机器数">
            {{ report_detail.engine_list.length }}
          </el-descriptions-item>
          <el-descriptions-item label="备注">
            {{ report_detail.remark }}
          </el-descriptions-item>
        </el-descriptions>
        <el-progress
          :text-inside="true"
          :stroke-width="24"
          :percentage="progress"
          status="success"
        />
      </content-wrap>
    </div>
    <div v-if="press_status === 2">
      <el-card>
        <span class="plan-name">{{ report_detail.plan_name }}</span>
        <div class="header">
          <div class="exec-schedule" v-show="report_detail.press_type !== 4">
            <div class="exec-time">
              <span class="status"> {{ pre_scene_end ? '运行中' : '前置场景执行中' }}... </span>
              <div>
                <span class="current-duration">
                  {{ pre_scene_end ? current_duration : prefix_scene_data_current_count }}
                </span>
                /
                <span class="total-duration">
                  {{ pre_scene_end ? total_duration : prefix_scene_data_count }}
                </span>
              </div>
            </div>
            <el-progress :percentage="time_percent" style="margin-top: 5px" />
          </div>
          <el-button
            type="danger"
            :disabled="authStore.getUserInfo?.id !== report_detail.create_user_id"
            :loading="stop_press_btn_loading"
            v-tc="{ name: '停止压测', report_id: report_id }"
            @click="stopTest"
          >
            停止压测
          </el-button>
        </div>
        <el-row :gutter="20" style="margin-top: 16px">
          <el-col :span="4">
            <el-statistic
              group-separator=","
              :value="report_data?.concurrency || 0"
              title="并发数"
            />
          </el-col>
          <el-col :span="4">
            <el-statistic
              group-separator=","
              :value="report_data?.stage_rps || 0"
              title="当前RPS(每秒)"
            />
          </el-col>
          <el-col :span="4">
            <el-statistic
              :precision="2"
              :value="report_data?.stage_rt || 0 + 'ms'"
              title="当前RT"
            />
          </el-col>
          <el-col :span="4">
            <el-statistic
              title="实时成功率"
              :value="(report_data?.stage_success_rate || 0) + '%'"
            />
          </el-col>
          <el-col :span="4">
            <el-statistic
              title="失败/总请求次数"
              :value="
                convertReqNum(report_data?.total_error_num || 0) +
                '/' +
                convertReqNum(report_data?.total_request_num || 0)
              "
            />
          </el-col>
          <el-col :span="4">
            <el-statistic
              title="发送/接受流量"
              :value="
                convertByte(report_data?.total_send_bytes) +
                '/' +
                convertByte(report_data?.total_received_bytes)
              "
            />
          </el-col>
        </el-row>
        <el-tabs
          v-model="active_tab"
          tab-position="right"
          style="height: 290px"
          @tab-change="changeChartType"
        >
          <el-tab-pane name="case">
            <template #label>
              <icon icon="ep:stopwatch" />
            </template>
            <plan-chart :plan-result="report_data.graph" :group-id="'report' + report_id" />
          </el-tab-pane>
          <el-tab-pane name="engine">
            <template #label>
              <icon icon="ep:data-line" />
            </template>
            <div style="width: 100%; height: 280px; display: flex">
              <cpu-chart :data="engine_data.cpu" style="width: 49%; height: 280px" />
              <mem-chart :data="engine_data.mem" style="width: 49%; height: 280px" />
            </div>
          </el-tab-pane>
          <el-tab-pane name="server_info">
            <template #label>
              <icon icon="ep:help-filled" />
            </template>
            <div style="width: 100%; height: 280px; flex-direction: column" class="flex">
              <div v-for="item in serverInfo" :key="item.index" class="flex" style="margin: 10px">
                <span>
                  NameSpace: <el-tag type="success">{{ item.namespace }}</el-tag>
                </span>
                <span style="margin-left: 20px">
                  容器: <el-tag type="success">{{ item.deployment_name }}</el-tag>
                </span>
                <span style="margin-left: 20px">
                  副本数:
                  <el-input-number
                    v-model="item.replicas"
                    :min="0"
                    controls-position="right"
                    @change="replicasChange(item)"
                  />
                </span>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-card>
      <!--错误集合信息-->
      <api-distribution :data="report_data" :case-name-map="caseNameMap" style="margin-top: 10px" />
      <url-distribution :data="report_data" style="margin-top: 10px" />
      <el-card
        v-if="report_detail.press_type === 3"
        style="margin-top: 10px"
        :body-style="{ padding: '10px 20px' }"
      >
        目标RPS：
        <el-popover placement="top" :width="185" @show="pre_target_rps = target_rps">
          <div style="display: flex">
            <el-input-number
              style="width: 90px"
              v-model="pre_target_rps"
              :min="1"
              :max="500000"
              :controls="false"
            />
            <el-button type="primary" @click="editPlanRps"> 保存 </el-button>
          </div>
          <template #reference>
            <span>
              {{ target_rps }}
              <el-button :icon="editIcon" link type="primary" />
            </span>
          </template>
        </el-popover>
      </el-card>
      <el-card
        v-for="(scene, scene_id_str) in report_data.scenes"
        :key="scene_id_str"
        class="scene-card"
      >
        <div class="scene-info">
          <div>
            <span class="scene-indicator"> {{ scene.scene_id }} </span>
            <span class="scene-title"> {{ sceneNameMap[scene.scene_id] || '--' }} </span>
            <span style="margin-left: 20px"> 并发数: {{ scene.concurrency }} </span>
          </div>

          <div>
            <el-button
              link
              type="primary"
              v-tc="{
                name: '采样日志',
                scene_id: scene.scene_id
              }"
              @click="getSamplingData(scene)"
            >
              采样日志
            </el-button>
            <el-button
              v-show="[0, 4].indexOf(report_detail.press_type) > -1"
              link
              type="primary"
              style="margin-left: 10px"
              v-tc="{
                name: '修改并发',
                scene_id: scene.scene_id
              }"
              @click="preChangeConcurrency(scene)"
            >
              修改并发
            </el-button>
            <el-popconfirm
              confirm-button-text="确定"
              cancel-button-text="取消"
              :icon="infoFillIcon"
              title="确定释放前置场景执行？"
              @confirm="releaseScene(scene.scene_id)"
            >
              <template #reference>
                <el-button
                  link
                  v-show="scene.scene_type === 1"
                  v-tc="{ name: '前置释放', scene_id: scene.scene_id }"
                  type="primary"
                >
                  前置释放
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
        <el-table :data="scene.cases" border stripe>
          <el-table-column label="ID" align="center">
            <template #default="{ row }">
              <el-button
                type="primary"
                link
                v-tc="{ name: '查看单接口', scene_id: scene.scene_id, case_id: row.case_id }"
                @click="getCaseGraph(row)"
              >
                {{ row.case_id }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="接口名" align="center">
            <template #default="{ row }">
              {{ caseNameMap[row.case_id] || '--' }}
            </template>
          </el-table-column>
          <el-table-column label="当前RPS" align="center">
            <template #default="{ row }">{{ row.stage_rps }}</template>
          </el-table-column>
          <el-table-column label="实时成功率" align="center">
            <template #default="{ row }">
              <span
                :style="
                  row.stage_success_rate <= 99 && scene.scene_type === 0
                    ? 'color:#F56C6C;font-weight:bold'
                    : ''
                "
              >
                {{ row.stage_success_rate }}%
              </span>
            </template>
          </el-table-column>
          <el-table-column label="实时RT" align="center">
            <template #default="{ row }">
              <span
                :style="
                  row.stage_avg_rt >= 500 && scene.scene_type === 0
                    ? 'color:#F56C6C;font-weight:bold'
                    : ''
                "
              >
                {{ row.stage_avg_rt }} ms
              </span>
            </template>
          </el-table-column>
          <el-table-column label="平均RPS" align="center">
            <template #default="{ row }"> {{ row.total_rps }}/s </template>
          </el-table-column>
          <el-table-column label="平均RT" align="center">
            <template #default="{ row }">
              <el-button
                type="primary"
                link
                @click="getTraces(row)"
                v-tc="{ name: '查看trace', scene_id: scene.scene_id, case_id: row.case_id }"
                :style="{ color: row.total_avg_rt >= 500 ? '#F56C6C' : '', fontWeight: 'bold' }"
              >
                {{ row.total_avg_rt }} ms
              </el-button>
              <!--              <span :style="row.total_avg_rt >= 500 ? 'color:#F56C6C;font-weight:bold' : ''">-->
              <!--                {{ row.total_avg_rt }} ms-->
              <!--              </span>-->
            </template>
          </el-table-column>
          <el-table-column label="成功率" align="center">
            <template #default="{ row }">
              <span :style="row.total_success_rate <= 99 ? 'color:#F56C6C;font-weight:bold' : ''">
                {{ row.total_success_rate }}%
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="two_xx_code_num" label="2xx个数" align="center" />
          <el-table-column prop="other_code_num" label="非2xx个数" align="center" />
          <!--          <el-table-column v-if="report_detail.press_type === 2" label="目标RPS">-->
          <!--            <template #default="scope">-->
          <!--              {{ caseRpsMap[scope.row.case_id] || 1 }}-->
          <!--              <el-popover placement="top" :width="185" @show="openEditRps(scope.row.case_id)">-->
          <!--                <div style="display: flex">-->
          <!--                  <el-input-number-->
          <!--                    style="width: 90px"-->
          <!--                    v-model="editTargetRps"-->
          <!--                    :min="1"-->
          <!--                    :max="100000"-->
          <!--                    :controls="false"-->
          <!--                  />-->
          <!--                  <el-button type="primary" @click="editRunningStepRps(scope.row.case_id)">-->
          <!--                    保存-->
          <!--                  </el-button>-->
          <!--                  &lt;!&ndash; <el-input v-model="editTargetRps" style="100%">-->
          <!--                              <template #append>-->
          <!--                                <el-button type="primary" @click="editRunningStepRps(row.case_id, 15)">-->
          <!--                                  保存-->
          <!--                                </el-button>-->
          <!--                              </template>-->
          <!--                            </el-input> &ndash;&gt;-->
          <!--                </div>-->
          <!--                <template #reference>-->
          <!--                  <el-button :icon="editIcon" link type="primary" />-->
          <!--                </template>-->
          <!--              </el-popover>-->
          <!--            </template>-->
          <!--          </el-table-column>-->
        </el-table>
      </el-card>
    </div>
    <div v-if="press_status === 3">压测结束，正在生成报告</div>
    <div v-if="press_status === 4">
      <el-card>
        <el-descriptions>
          <template #title>
            <div style="display: flex">
              <el-button
                link
                type="primary"
                style="font-size: 18px; font-weight: bold"
                v-saclink="{ name: '查看计划', plan_id: report_detail.plan_id }"
                @click="jumpPlanDetail(report_detail.plan_id)"
              >
                {{ report_detail.plan_name }}
              </el-button>
              <div style="margin-left: 10px">报告编号：{{ report_detail.report_id }}</div>
            </div>
          </template>
          <el-descriptions-item label="预设时长">
            {{ report_detail.duration }} 分钟
          </el-descriptions-item>
          <el-descriptions-item label="压测模式">
            {{ press_type_mp[report_detail.press_type] }}
          </el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ report_detail.start_time }}
          </el-descriptions-item>
          <el-descriptions-item label="结束时间">
            {{ report_detail.end_time }}
          </el-descriptions-item>
          <el-descriptions-item label="执行人">
            {{ report_detail.create_user_name }}
          </el-descriptions-item>
          <el-descriptions-item label="总请求数">
            {{ report_data?.total_request_num }}
          </el-descriptions-item>
          <el-descriptions-item label="成功率">
            {{ report_data?.total_success_rate }}%
          </el-descriptions-item>
          <el-descriptions-item label="总流量(请求/响应)">
            {{ convertByte(report_data?.total_send_bytes) }} /
            {{ convertByte(report_data?.total_received_bytes) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
      <plan-chart
        :plan-result="report_data.graph"
        :group-id="'report' + report_id"
        style="margin-top: 10px"
      />
      <api-distribution :data="report_data" :case-name-map="caseNameMap" style="margin-top: 10px" />
      <url-distribution :data="report_data" style="margin-top: 10px" />

      <div v-for="item in report_data.scenes" :key="item.scene_id" style="margin-top: 10px">
        <down-scene-detail
          :data="item"
          :report-id="report_id"
          :sceneNameMap="sceneNameMap"
          :caseNameMap="caseNameMap"
        />
      </div>
    </div>
  </el-skeleton>
  <el-dialog
    :title="sceneNameMap[change_param.scene_id]"
    v-model="chang_con_dialog_visible"
    width="700px"
  >
    <el-form ref="form" :model="change_param" label-width="80px">
      <el-form-item label="并发数">
        <el-input v-model.number="change_param.concurrency" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="chang_con_dialog_visible = false">取 消</el-button>
        <el-button type="primary" @click="changeConcurrency">确 定</el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="case_graph_dialog_visible" :title="case_graph_title" width="90%">
    <plan-chart
      :plan-result="case_data.graph"
      :group-id="'case' + case_graph_id"
      style="margin-top: 10px"
    />
  </el-dialog>
  <sampling-data ref="samplingDataRef" />
  <traces-list ref="tracesListRef" />
</template>

<style scoped lang="less">
.plan-name {
  font-size: 20px;
  color: var(--el-color-primary);
}

:deep(.el-col) {
  text-align: center;
}

:deep(.el-tabs__item) {
  padding: 0 0 0 5px;
}

.header {
  display: flex;
  justify-content: space-between;

  .exec-schedule {
    display: flex;
    flex-direction: column;
    width: 85%;

    .exec-time {
      display: flex;
      justify-content: space-between;
      width: 100%;

      .status {
        font-size: small;
        color: gray;
      }

      .current-duration {
        font-size: 20px;
        color: var(--el-color-primary);
      }

      .total-duration {
        font-size: 16px;
        color: gray;
      }
    }
  }
}

.scene-card {
  margin-top: 10px;
  margin-bottom: 10px;
  width: 100%;

  :deep(.el-card__body) {
    width: calc(99% - 20px);
    padding: 10px 20px;
  }

  .scene-info {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;

    .scene-indicator {
      font-size: 13px;
      margin-left: 10px;
    }

    .scene-title {
      margin-left: 10px;
      font-size: 15px;
      font-weight: bold;
    }
  }

  .statistic_content {
    font-size: 18px;
    color: var(--el-color-primary);
  }
}
.stop-button-right {
  text-align: right;
}
</style>

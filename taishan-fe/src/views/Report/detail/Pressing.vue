<script setup lang="ts">
import ReportService from '@/api/scene/report'
import SceneService from '@/api/scene/scene'
import { onActivated, onDeactivated, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PlanGraph from '@/views/Report/detail/PlanChart.vue'
import SamplingData from '@/views/Report/sampling/SamplingData.vue'
import { useIcon } from '@/hooks/web/useIcon'
import { isNumber } from '@/utils/is'
import SceneReportService from '@/api/report'

const props = defineProps({
  reportId: {
    type: Number,
    default: 0
  }
})

const editIcon = useIcon({ icon: 'eva:edit-2-outline' })
const emit = defineEmits(['preDownPress'])

const sceneReportService = new SceneReportService()
const reportService = new ReportService()
const sceneService = new SceneService()
const sceneNameMap = ref<any>({})
const caseNameMap = ref<any>({})

let plan_loading = ref(true)

let plan_id, plan_name, duration, start_time
const press_type = ref(0)
let current_duration = ref('00:00:00')
let total_duration = ref('preDownPress')
let time_percent = ref<number>(0.0)
let stop_press_btn_loading = ref(false)
let data_interval: any = null
let time_interval: any = null
let table_data = ref<any>({
  graph: {},
  scenes: []
})
const caseRpsMap = ref({})
const editTargetRps = ref(0)
const planRps = ref(1)
const editPlanRpsValue = ref(1)

function openEditRps(caseId) {
  editTargetRps.value = getCaseRps(caseId)
}

function openEditPlanRps() {
  editPlanRpsValue.value = planRps.value
}

function getReportRps() {
  sceneReportService
    .getReportRps({
      report_id: props.reportId
    })
    .then((res) => {
      const d = res.data || {}
      planRps.value = d[props.reportId]
    })
    .finally(() => {
      // caseRpsMap.value = {
      //   '155': 11
      // }
    })
}

function getReportRpsList() {
  sceneReportService
    .getReportRpsList({
      report_id: props.reportId
    })
    .then((res) => {
      caseRpsMap.value = res.data || {}
    })
    .finally(() => {
      // caseRpsMap.value = {
      //   '155': 11
      // }
    })
}

function getCaseRps(caseId) {
  return caseRpsMap.value[caseId] || 1
}

function editRunningStepRps(caseId) {
  // sceneReportService
  //   .editRunningStepRps({
  //     report_id: props.reportId,
  //     case_id: caseId,
  //     new_rps: editTargetRps.value
  //   })
  //   .then(() => {
  //     getReportRpsList()
  //     ElMessage.success('修改目标RPS成功')
  //   })
  //   .finally(() => {})
}

function editPlanRps() {
  sceneReportService
    .reportRpsModify({
      report_id: props.reportId,
      new_rps: editPlanRpsValue.value
    })
    .then(() => {
      getReportDetail()
      ElMessage.success('修改计划RPS成功')
    })
    .finally(() => {})
}

const getReportDetail = async () => {
  const res = await reportService.getReportDetail(props.reportId)
  plan_id = res.plan_id
  plan_name = res.plan_name
  duration = res.duration
  start_time = res.start_time
  press_type.value = res.press_type
  planRps.value = res.rps
  total_duration.value = getTotalDuration(res.duration)
  getPlanScenes(plan_id)
  await searchSceneList(plan_id)
  time_interval = setInterval(() => {
    calcSchedule()
  }, 1000)
  if (press_type.value === 2) {
    //RPS模式
    getReportRpsList()
  }
  startGetReportDetail()
}

const getPlanScenes = async (plan_id) => {
  sceneNameMap.value = {}
  caseNameMap.value = {}
  const scene_list = (await sceneService.getSceneList(plan_id)) || []
  for (let scene of scene_list) {
    sceneNameMap.value[scene.scene_id + ''] = scene.scene_name
    if (scene.case_tree) {
      buildCaseNameMap(scene.case_tree)
    }
  }
}

function buildCaseNameMap(case_tree) {
  if (!case_tree) return
  for (let case1 of case_tree) {
    caseNameMap.value[case1.case_id + ''] = case1.title
    if (case1.children) {
      buildCaseNameMap(case1.children)
    }
  }
}

const startGetReportDetail = () => {
  if (!data_interval) {
    data_interval = setInterval(getGraphData, 5000)
  }
}

const stopGetReportDetail = () => {
  if (data_interval) {
    clearInterval(data_interval)
    data_interval = null
  }
}

const getGraphData = async () => {
  const res = (await reportService.getReportData(props.reportId)) || {}
  if (res?.end) {
    if (data_interval !== null) {
      ElMessage.success('本次压测已结束！')
      clearInterval(data_interval)
    }
    emit('preDownPress')
  }
  table_data.value = res || {}
}
const getTotalDuration = (duration) => {
  const hours = String(Math.floor(duration / 60)).padStart(2, '0')
  const minutes = String(duration % 60).padStart(2, '0')
  return `${hours}:${minutes}:00`
}

const calcSchedule = () => {
  const startTime = new Date(start_time)
  const currentTime = new Date()
  //@ts-ignore
  const secondDiff = Math.floor(Math.abs(currentTime - startTime) / 1000)
  const seconds = String(Math.floor(secondDiff % 60)).padStart(2, '0')
  const minutes = String(Math.floor((secondDiff / 60) % 60)).padStart(2, '0')
  const hours = String(Math.floor(secondDiff / 3600)).padStart(2, '0')
  current_duration.value = `${hours}:${minutes}:${seconds}`
  const percent = parseFloat(((secondDiff * 100) / (duration * 60)).toFixed(1))
  if (percent >= 100) {
    emit('preDownPress')
  }
  time_percent.value = percent > 100 ? 100 : percent
}

const stopTest = () => {
  stop_press_btn_loading.value = true
  reportService.stopPress(props.reportId)
  if (data_interval !== null) {
    ElMessage.success('手动停止压测，本次压测已结束！')
    clearInterval(data_interval)
  }
  if (time_interval !== null) {
    clearInterval(time_interval)
  }
  emit('preDownPress')
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

let all_scene_list: Array<any> = []
const searchSceneList = async (plan_id) => {
  const res = await sceneService.getSceneList(plan_id)
  all_scene_list = [...res]
  plan_loading.value = false
}
const getSceneNameById = (id) => {
  const _scene_id = Number(id)
  if (all_scene_list.length === 0) {
    return
  }
  if (all_scene_list.find((e) => e.scene_id === _scene_id) === undefined) {
    return
  }
  return all_scene_list.find((e) => e.scene_id === _scene_id).scene_name
}

const apiTableFormat = (scene) => {
  if (scene === undefined) {
    return
  }
  const obj = scene.cases
  const res: any = []
  Object.values(obj).forEach((val: any) => {
    const _item = { ...val }
    _item.two_xx = val['2xx_code_num']
    _item.not_two_xx = val.other_code_num
    res.push(_item)
  })
  return res
}

let chang_con_dialog_visible = ref(false)
let change_param = reactive({
  scene_id: 0,
  concurrency: 0,
  report_id: 0
})
const preChangeConcurrency = (item) => {
  chang_con_dialog_visible.value = true
  change_param.scene_id = item.scene_id
  change_param.concurrency = item.concurrency
}

const changeConcurrency = () => {
  chang_con_dialog_visible.value = false
  change_param.report_id = props.reportId
  if (!isNumber(change_param.concurrency)) {
    ElMessage.warning('请输入正整数')
  }
  change_param.concurrency = Number(change_param.concurrency)
  reportService.updatePress(change_param)
}

const samplingDataRef = ref()
const getSamplingData = (scene_id) => {
  const param = {
    report_id: props.reportId,
    scene_id: scene_id
  }
  samplingDataRef.value.drawerOpen(param)
}

getReportDetail()

onActivated(() => {
  startGetReportDetail()
})

// 路由离开前停止定时任务
onDeactivated(() => {
  stopGetReportDetail()
})
</script>

<template>
  <el-card>
    <el-skeleton :loading="plan_loading">
      <span class="plan-name">{{ plan_name }}</span>
      <div class="header">
        <div class="exec-schedule">
          <div class="exec-time">
            <span class="status"> 运行中... </span>
            <div>
              <span class="current-duration">{{ current_duration }}</span>
              /
              <span class="total-duration">{{ total_duration }}</span>
            </div>
          </div>
          <el-progress :percentage="time_percent" style="margin-top: 5px" />
        </div>
        <el-button type="danger" :loading="stop_press_btn_loading" @click="stopTest">
          停止压测
        </el-button>
      </div>
      <el-row :gutter="20" style="margin-top: 16px">
        <el-col :span="4">
          <el-statistic
            group-separator=","
            :value="table_data?.current?.concurrency || 0"
            title="并发数"
          />
        </el-col>
        <el-col :span="4">
          <el-statistic
            group-separator=","
            :value="table_data?.current?.rps || 0"
            title="当前RPS(每秒)"
          />
        </el-col>
        <el-col :span="4">
          <el-statistic
            :precision="2"
            :value="table_data?.current?.rt || 0 + 'ms'"
            title="当前RT"
          />
        </el-col>
        <el-col :span="4">
          <el-statistic
            title="接口成功率"
            :value="(table_data?.current?.success_rate || 0) + '%'"
          />
        </el-col>
        <el-col :span="4">
          <el-statistic
            title="失败/总请求次数"
            :value="
              convertReqNum(table_data?.total?.error_num || 0) +
              '/' +
              convertReqNum(table_data?.total?.request_num || 0)
            "
          />
        </el-col>
        <el-col :span="4">
          <el-statistic
            title="发送/接受流量"
            :value="
              convertByte(table_data?.current?.send_bytes) +
              '/' +
              convertByte(table_data?.current?.received_bytes)
            "
          />
        </el-col>
      </el-row>
      <plan-graph :plan-result="table_data.graph" style="margin-top: 10px" />
    </el-skeleton>
  </el-card>
  <el-card v-if="press_type === 3" style="margin-top: 10px" :body-style="{ padding: '10px 20px' }">
    修改计划RPS：
    <el-popover placement="top" :width="185" @show="openEditPlanRps()">
      <div style="display: flex">
        <el-input-number
          style="width: 90px"
          v-model="editPlanRpsValue"
          :min="1"
          :max="100000"
          :controls="false"
        />
        <el-button type="primary" @click="editPlanRps()"> 保存 </el-button>
      </div>
      <template #reference>
        <span>
          {{ planRps }}
          <el-button :icon="editIcon" link type="primary" />
        </span>
      </template>
    </el-popover>
  </el-card>
  <el-card
    v-for="(scene, scene_id_str) in table_data.scenes"
    :key="scene_id_str"
    class="scene-card"
  >
    <div class="scene-info">
      <span>
        {{ sceneNameMap[scene.scene_id] || '--' }}
      </span>
      <div>
        <el-button link type="primary" @click="getSamplingData(scene.scene_id)">采样日志</el-button>
        <el-button
          v-if="press_type === 0"
          link
          type="primary"
          style="margin-left: 10px"
          @click="preChangeConcurrency(scene)"
        >
          修改并发
        </el-button>
      </div>
    </div>
    <!-- <span style="margin-top: 20px">
      共设置
      <span style="color: red">
        {{ scene.parameter_count }}
      </span>
      个参数
    </span> -->
    <el-table :data="apiTableFormat(scene)">
      <el-table-column prop="case_name" label="接口名">
        <template #default="{ row }">
          {{ caseNameMap[row.case_id] || '--' }}
        </template>
      </el-table-column>
      <el-table-column prop="success_rps" label="成功数">
        <template #default="{ row }"> {{ row.success_rps }}/s </template>
      </el-table-column>
      <el-table-column prop="error_rps" label="失败数">
        <template #default="{ row }"> {{ row.error_rps }}/s </template>
      </el-table-column>
      <el-table-column prop="avg_rt" label="平均响应时间">
        <template #default="{ row }"> {{ row.avg_rt }} ms</template>
      </el-table-column>
      <el-table-column prop="concurrency" label="并发数" />
      <el-table-column prop="two_xx_code_num" label="2xx个数" />
      <el-table-column prop="other_code_num" label="非2xx个数" />
      <el-table-column prop="rally_point" label="集合点" />
      <el-table-column prop="avg_rps" label="当前RPS" />
      <el-table-column v-if="press_type === 2" label="目标RPS">
        <template #default="scope">
          {{ getCaseRps(scope.row.case_id) }}
          <el-popover placement="top" :width="185" @show="openEditRps(scope.row.case_id)">
            <div style="display: flex">
              <el-input-number
                style="width: 90px"
                v-model="editTargetRps"
                :min="1"
                :max="100000"
                :controls="false"
              />
              <el-button type="primary" @click="editRunningStepRps(scope.row.case_id)">
                保存
              </el-button>
              <!-- <el-input v-model="editTargetRps" style="100%">
                <template #append>
                  <el-button type="primary" @click="editRunningStepRps(row.case_id, 15)">
                    保存
                  </el-button>
                </template>
              </el-input> -->
            </div>
            <template #reference>
              <el-button :icon="editIcon" link type="primary" />
            </template>
          </el-popover>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog
    :title="getSceneNameById(change_param.scene_id)"
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
  <sampling-data ref="samplingDataRef" />
</template>

<style scoped lang="less">
.plan-name {
  font-size: 20px;
  color: var(--el-color-primary);
}

:deep(.el-col) {
  text-align: center;
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
  //height: 600px;

  .scene-info {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
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

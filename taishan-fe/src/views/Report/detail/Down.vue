<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import ReportService from '@/api//scene/report'
import PlanGraph from '@/views/Report/detail/PlanChart.vue'
import DownSceneDetail from '@/views/Report/detail/DownSceneDetail.vue'
import SceneService from '@/api/scene/scene'

const props = defineProps({
  reportId: {
    type: Number,
    default: 0
  }
})

const sceneService = new SceneService()
const reportService = new ReportService()

const data_loading = ref(false)
const sceneNameMap = ref<any>({})
const caseNameMap = ref<any>({})

let report_detail = reactive<any>({})
let report_data = reactive<any>({})

const press_type_mp = {
  0: '并发模式',
  1: '阶梯式加压',
  2: 'RPS模式',
  3: 'RPS比例模式'
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
const getReportDetail = async () => {
  data_loading.value = true
  try {
    const detail = await reportService.getReportDetail(props.reportId)
    report_detail = Object.assign(report_detail, detail)
    const data = await reportService.getReportData(props.reportId)
    getPlanScenes(detail?.plan_id)
    report_data = Object.assign(report_data, data)
  } finally {
    data_loading.value = false
  }
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
onMounted(() => {
  getReportDetail()
})
</script>

<template>
  <el-skeleton :loading="data_loading" :rows="12" animated>
    <el-card>
      <el-descriptions :title="'报告编号：' + report_detail.report_id">
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
          {{ report_data?.current?.request_num }}
        </el-descriptions-item>
        <el-descriptions-item label="成功率">
          {{ report_data?.current?.success_rate }}%
        </el-descriptions-item>
        <el-descriptions-item label="总流量(请求/响应)">
          {{ convertByte(report_data?.current?.send_bytes) }} /
          {{ convertByte(report_data?.current?.received_bytes) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
    <plan-graph :plan-result="report_data.graph" style="margin-top: 10px" />
    <div v-for="item in report_data.scenes" :key="item.scene_id" style="margin-top: 10px">
      <down-scene-detail
        :data="item"
        :report-id="props.reportId"
        :sceneNameMap="sceneNameMap"
        :caseNameMap="caseNameMap"
      />
    </div>
  </el-skeleton>
</template>

<style scoped lang="less"></style>

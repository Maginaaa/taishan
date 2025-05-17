<script setup lang="ts">
import { reactive, ref } from 'vue'
import SamplingData from '@/views/Report/sampling/SamplingData.vue'
import PlanChart from '@/views/Report/detail/PlanChart.vue'
import ReportService from '@/api/scene/report'
import TracesList from '@/views/Report/detail/TracesList.vue'

const reportService = new ReportService()

const props = defineProps({
  data: {
    type: Object,
    default: () => {}
  },
  reportId: {
    type: Number,
    default: 0
  },
  sceneNameMap: {
    type: Object,
    default: () => {}
  },
  caseNameMap: {
    type: Object,
    default: () => {}
  }
})

const samplingDataRef = ref()
const getSamplingData = (scene) => {
  const param = {
    report_id: props.reportId,
    scene_id: scene.scene_id
  }
  let _case_list = []
  scene.cases.forEach((c) => {
    _case_list.push({ id: c.case_id, value: props.caseNameMap[c.case_id] })
  })
  samplingDataRef.value.drawerOpen(param, _case_list)
}

let case_graph_dialog_visible = ref(false)
let case_graph_title = ref('')
let case_graph_id = ref<number>(0)
let case_data = reactive<any>({})

const getCaseGraph = async (row) => {
  case_graph_dialog_visible.value = true
  case_graph_title.value = row.case_name
  case_graph_id.value = row.case_id
  let res = await reportService.getReportCaseData({
    report_id: props.reportId,
    case_id: row.case_id
  })
  case_data = Object.assign(case_data, res)
}

const tracesListRef = ref()
const getTraces = async (row) => {
  tracesListRef.value.drawerOpen(row.case_id, props.reportId)
  // const res = await dataService.getSlsData(row.case_id)
}
</script>

<template>
  <div>
    <el-card class="scene-info">
      <div style="width: 100%">
        <span class="scene-indicator">{{ props.data.scene_id }}</span>
        <span class="scene-title">{{ props.sceneNameMap[props.data.scene_id] }}</span>
        <el-button
          link
          type="primary"
          style="float: right"
          v-tc="{
            name: '采样日志',
            scene_id: props.data.scene_id
          }"
          @click="getSamplingData(props.data)"
        >
          采样日志
        </el-button>

        <!--      <el-button link type="primary" style="float: right">-->
        <!--        查看明细-->
        <!--        <el-icon class="el-icon&#45;&#45;right">-->
        <!--          <icon icon="ep:view" style="color: #409eff" />-->
        <!--        </el-icon>-->
        <!--      </el-button>-->
      </div>
      <el-table :data="props.data.cases" border stripe class="table">
        <el-table-column label="ID" align="center">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              v-tc="{
                name: '查看单接口',
                scene_id: props.data.scene_id,
                case_id: row.case_id
              }"
              @click="getCaseGraph(row)"
              >{{ row.case_id }}</el-button
            >
          </template>
        </el-table-column>
        <el-table-column prop="case_name" label="接口名" align="center">
          <template #default="scope">
            <span> {{ props.caseNameMap[scope.row.case_id] }}</span>
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
              v-tc="{
                name: '查看trace',
                scene_id: props.data.scene_id,
                case_id: row.case_id
              }"
              :style="{ color: row.total_avg_rt >= 500 ? '#F56C6C' : '', fontWeight: 'bold' }"
            >
              {{ row.total_avg_rt }} ms
            </el-button>
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
        <!--        <el-table-column label="最小RT" align="center">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.min_rt }} ms</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <!--        <el-table-column label="最大RT" align="center">-->
        <!--          <template #default="scope">-->
        <!--            <span>{{ scope.row.max_rt }} ms</span>-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column prop="ninety_request_time_line_value" label="90线RT" align="center">
          <template #default="scope">
            <span>{{ scope.row.ninety_request_time_line_value }} ms</span>
          </template>
        </el-table-column>
        <el-table-column prop="ninety_five_request_time_line_value" label="95线RT" align="center">
          <template #default="scope">
            <span>{{ scope.row.ninety_five_request_time_line_value }} ms</span>
          </template>
        </el-table-column>
        <el-table-column prop="ninety_nine_request_time_line_value" label="99线RT" align="center">
          <template #default="scope">
            <span>{{ scope.row.ninety_nine_request_time_line_value }} ms</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
  <sampling-data ref="samplingDataRef" />
  <traces-list ref="tracesListRef" />
  <el-dialog v-model="case_graph_dialog_visible" :title="case_graph_title" width="90%">
    <plan-chart
      :plan-result="case_data.graph"
      :group-id="'case' + case_graph_id"
      style="margin-top: 10px"
    />
  </el-dialog>
</template>

<style scoped lang="less">
:deep(.el-card__body) {
  width: 100%;
  padding: 10px 20px;
}

.scene-info {
  width: 100%;
  display: flex;
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

  .table {
    //width: calc(100% - 10px);
    margin-top: 10px;
  }
}
</style>

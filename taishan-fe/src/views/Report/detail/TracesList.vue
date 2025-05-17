<script setup lang="ts">
import { ref } from 'vue'
import DataService from '@/api/data'
import TraceDetail from '@/views/Report/component/TraceDetail.vue'

const dataService = new DataService()

let visible = ref(false)
let inner_visible = ref(false)
let table_loading = ref(false)
let report_id = ref(null)

const drawerOpen = async (caseId, reportId) => {
  visible.value = true
  report_id.value = reportId
  try {
    table_loading.value = true
    table_data.value = await dataService.getTracesList({
      case_id: caseId,
      report_id: report_id.value
    })
  } finally {
    table_loading.value = false
  }
}

const table_data = ref([])
const trace_detail = ref(null)
const current_trace_id = ref('')

const getTraceDetail = async (traceId) => {
  current_trace_id.value = traceId
  const res = await dataService.getTraceDetail({
    trace_id: traceId,
    report_id: report_id.value
  })
  trace_detail.value = res
  inner_visible.value = true
}

defineExpose({
  drawerOpen
})
</script>

<template>
  <el-dialog width="70%" v-model="visible">
    <el-table
      :data="table_data"
      v-loading="table_loading"
      border
      stripe
      style="width: 100%"
      max-height="660px"
    >
      <el-table-column align="center" header-align="center" prop="TraceID" label="TraceID">
        <template #default="{ row }">
          <el-button type="primary" link @click="getTraceDetail(row.TraceID)">
            {{ row.TraceID }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column align="center" header-align="center" prop="OperationName" label="path" />
      <el-table-column align="center" header-align="center" prop="Duration" label="响应时间">
        <template #default="{ row }">
          <span>{{ row.Duration }} ms</span>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="inner_visible" width="60%" :title="current_trace_id" append-to-body>
      <trace-detail :trace-detail="trace_detail" />
    </el-dialog>
  </el-dialog>
</template>

<style scoped lang="less"></style>

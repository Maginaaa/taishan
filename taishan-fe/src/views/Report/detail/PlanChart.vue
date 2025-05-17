<script setup lang="ts">
import TpsChart from '@/views/Report/chart/TpsChart.vue'
import ResponseTimeChart from '@/views/Report/chart/ResponseTimeChart.vue'
import SuccessRateChart from '@/views/Report/chart/SuccessRateChart.vue'
import { reactive, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  planResult: {
    type: Object,
    default: () => {
      return {}
    }
  },
  groupId: {
    type: String
  }
})

let mp = reactive({
  concurrency: [] as any,
  tps: [] as any,
  response_time: [] as any,
  success_rate: [] as any
})

const calcGraphData = () => {
  mp.tps = props.planResult?.rps || []
  mp.response_time = props.planResult?.request_time || []
  mp.success_rate = props.planResult?.success_rate || []
}

watch(
  () => props.planResult,
  () => {
    calcGraphData()
  },
  { immediate: true, deep: true }
)

watch(
  () => props.groupId,
  (val) => {
    if (val !== undefined) {
      echarts.connect(val)
    }
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <div class="graph-row">
    <el-card class="graph-item">
      <tps-chart :data="mp.tps" :group="props.groupId" />
    </el-card>
    <el-card class="graph-item">
      <response-time-chart :data="mp.response_time" :group="props.groupId" />
    </el-card>
    <el-card class="graph-item" style="margin-right: 0">
      <success-rate-chart :data="mp.success_rate" :group="props.groupId" />
    </el-card>
  </div>
</template>

<style scoped lang="less">
:deep(.el-card__body) {
  padding: 10px 0 0 0;
}
.graph-row {
  width: 100%;
  display: flex;
  flex-direction: row;

  .graph-item {
    width: 33%;
    height: 280px;
    margin: 0 5px;
  }
}
</style>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { baseChartOption } from '@/views/Report/chart/chart-data'
import { Echart } from '@/components/Echart'
import type { EChartsOption } from 'echarts/types/dist/echarts'
import { cloneDeep } from 'lodash-es'

let optionsData = reactive<EChartsOption>(cloneDeep(baseChartOption))

const props = defineProps({
  data: {
    type: Array,
    default: baseChartOption
  },
  group: {
    type: String
  }
})

const initGraph = () => {
  optionsData.title.text = '并发数(全场景)'
}

const refreshData = () => {
  const sources = [['time', '并发数(全场景)']]
  sources.push(...props.data)
  optionsData.dataset.source = sources
}

watch(
  () => props.data,
  () => {
    refreshData()
  },
  { immediate: true, deep: true }
)
initGraph()
</script>

<template>
  <echart :options="optionsData" :height="280" :group="props.group" />
</template>

<style scoped lang="less"></style>

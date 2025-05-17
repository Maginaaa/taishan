<script setup lang="ts">
import { reactive, watch } from 'vue'
import { baseChartOption } from '@/views/Report/chart/chart-data'
import { Echart } from '@/components/Echart'
import { cloneDeep } from 'lodash-es'

let optionsData = reactive<any>(cloneDeep(baseChartOption))

const props = defineProps({
  data: {
    type: Object,
    default: baseChartOption
  },
  group: {
    type: String
  }
})

const initGraphData = () => {
  optionsData.title.text = '接口成功率'
  optionsData.yAxis = Object.assign(optionsData.yAxis, {
    min: 0,
    max: 100,
    interval: 20,
    axisLabel: {
      show: true,
      interval: 'auto',
      formatter: function (value) {
        return value + '%'
      }
    }
  })
  optionsData.tooltip = {
    trigger: 'axis',
    formatter(params) {
      // 使用 map 函数构建每个场景的数据字符串
      return params
        .map((param) => {
          return param.marker + `${param.seriesName} ：${param.data[1]}%`
        })
        .join('<br/>')
    }
  }
  optionsData.visualMap = {
    show: false,
    pieces: [
      {
        gt: 0,
        lte: 99,
        color: '#F56C6C'
      },
      {
        gt: 99,
        color: '#5470c6'
      }
    ]
  }
}

const refreshData = () => {
  const sources: Array<any> = [['time', '接口成功率']]
  //@ts-ignore
  sources.push(...props.data)
  optionsData.dataset.source = sources
}

initGraphData()
watch(
  () => props.data,
  () => {
    refreshData()
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <echart :options="optionsData" :height="280" :group="props.group" />
</template>

<style scoped lang="less"></style>

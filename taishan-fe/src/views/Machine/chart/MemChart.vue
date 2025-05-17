<script setup lang="ts">
import { reactive, watch } from 'vue'
import { cloneDeep } from 'lodash-es'
import { baseChartOption } from '@/views/Report/chart/chart-data'
import { Echart } from '@/components/Echart'

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
  optionsData.title.text = '内存信息'
  optionsData.grid = {
    left: '40px',
    right: '0px',
    bottom: '50px'
  }
  optionsData.yAxis = {
    type: 'value',
    splitNumber: 4,
    splitLine: {
      show: true, // 显示分割线
      lineStyle: {
        color: ['rgba(0,0,0,0.1)'], // 分割线的颜色，可以设置为数组来显示不同颜色的虚线段
        width: 1, // 分割线的宽度
        type: 'dashed' // 分割线类型，'dashed' 表示虚线
      }
    },
    axisLabel: {
      show: true,
      interval: 'auto',
      formatter: function (value) {
        return value + '%'
      }
    }
  }
  optionsData.tooltip = {
    trigger: 'axis',
    confine: true,
    formatter(params) {
      return (
        params[0].axisValue +
        `<br/>` +
        params
          .map((param) => {
            return `<div style='width: 320px;height:6px;display: flex;justify-content:space-between'>
                      <div>${param.marker} ${param.seriesName}</div>
                      <div style='font-weight: bold'>${param.data[param.seriesIndex + 1]}%</div>
                    </div>`
          })
          .join('<br/>')
      )
    }
  }
}

const refreshData = () => {
  optionsData.dataset.source = props.data
  const _series: any = []
  if (props.data[0] !== undefined) {
    for (let i = 0; i < props.data[0].length - 1; i++) {
      _series.push({
        type: 'line',
        showSymbol: false,
        areaStyle: {}
      })
    }
  }
  optionsData.series = _series
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
  <echart :options="optionsData" :height="280" />
</template>

<style scoped lang="less"></style>

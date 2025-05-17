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

function formatNumber(num) {
  if (num >= 1000 && num < 10000) {
    return (num / 1000).toFixed(1) + 'K'
  } else if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  } else {
    return num.toString()
  }
}

const initGraphData = () => {
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
      formatter(value) {
        return formatNumber(value)
      }
    }
  }
  optionsData.title.text = '响应时间'
  optionsData.series = [
    {
      type: 'line',
      showSymbol: false,
      markLine: {
        symbol: ['circle', 'circle'],
        data: [
          {
            type: 'average',
            name: 'Avg',
            label: {
              position: 'insideEnd' // 表现内容展示的位置
            }
          }
        ]
      }
    }
  ]
  optionsData.tooltip = {
    trigger: 'axis',
    formatter(params) {
      // 使用 map 函数构建每个场景的数据字符串
      return params
        .map((param) => {
          return param.marker + `${param.seriesName} ：${param.data[1]}ms`
        })
        .join('<br/>')
    }
  }
  optionsData.visualMap = {
    show: false,
    pieces: [
      {
        gte: 0,
        lt: 500,
        color: '#5470c6'
      },
      {
        gte: 500,
        color: '#F56C6C'
      }
    ]
  }
}

const refreshData = () => {
  const sources = [['time', '响应时间']]
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

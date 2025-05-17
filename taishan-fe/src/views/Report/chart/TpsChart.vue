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
  optionsData.legend = {
    icon: 'circle',
    x: 'left',
    left: 16
  }
  optionsData.toolbox = {
    feature: {
      dataZoom: {
        yAxisIndex: false,
        title: {
          zoom: '缩放',
          back: '还原'
        }
      }
    }
  }
  optionsData.title.text = '并发数/RPS'
  optionsData.grid.right = '45px'
  optionsData.yAxis = [
    {
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
    },
    {
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
          return formatNumber(value)
        }
      }
    }
  ]
  optionsData.series = [
    {
      type: 'line',
      showSymbol: false
    },
    {
      type: 'line',
      showSymbol: false,
      yAxisIndex: 1,
      markPoint: {
        data: [
          {
            type: 'max',
            name: 'Max',
            animation: true
          }
        ]
      },
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
      return (
        params[0].axisValue +
        `<br/>` +
        params
          .map((param) => {
            return (
              param.marker +
              `${param.seriesName} ：${param.data[param.seriesIndex + 1]}` +
              `${param.seriesIndex === 0 ? '' : '/s'}`
            )
          })
          .join('<br/>')
      )
    }
  }
}

const refreshData = () => {
  const sources: Array<any> = [['time', '并发数', 'RPS']]
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

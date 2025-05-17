<script setup lang="ts">
import { reactive, watch } from 'vue'
import { cloneDeep } from 'lodash-es'
import { baseChartOption } from '@/views/Report/chart/chart-data'
import { Echart } from '@/components/Echart'
import UrlDetail from '@/views/Report/detail/UrlDetail.vue'
import { extractPath } from '@/utils/url'

const props = defineProps({
  data: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

let optionsData = reactive<any>(cloneDeep(baseChartOption))

const initGraphData = () => {
  optionsData.title.text = 'URL聚合信息'
  optionsData.grid = {
    top: '10px'
  }
  optionsData.xAxis = []
  optionsData.yAxis = []
  optionsData.series = [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      encode: {
        itemName: 'path', // 使用第一列作为名称
        value: ['count'] // 使用第二、三列作为值
      },
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10
      },
      label: {
        show: false
      },
      labelLine: {
        show: false
      }
    }
  ]
  optionsData.tooltip = {
    trigger: 'item',
    appendToBody: true,
    formatter(params) {
      // 获取原始数据
      var data = params.data
      if (!data || data.length === 0) return ''

      return `
        ${params.marker}${data[0]}<br/>
        总请求数: ${data[1]} 占比: ${params.percent}%<br/>
        成功数: ${data[1] - data[2]} 错误数: ${data[2]}
    `
    }
  }
}

const refreshData = () => {
  const sources = [['path', 'count', 'error_count']]
  const baseData = props.data?.error_url_array || []
  for (let item of baseData) {
    sources.push([extractPath(item.url), item.count, item.error_count])
  }
  optionsData.dataset.source = sources
}

initGraphData()

watch(
  () => props.data.error_url_array,
  () => {
    refreshData()
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <div class="table-row">
    <el-card class="table-item">
      <div class="flex items-center">
        <echart :options="optionsData" :height="280" style="width: 25%" />
        <url-detail style="width: 75%" :api-data="props.data.error_url_array" />
      </div>
    </el-card>
  </div>
</template>

<style scoped lang="less">
:deep(.el-card__body) {
  padding: 3px 0 0 0;
}
.table-row {
  width: 100%;
  display: flex;
  flex-direction: row;

  .table-item {
    width: 100%;
    height: 280px;
    margin: 0 5px;

    .center-text {
      text-align: center; /* 将文本水平居中 */
      font-size: 15px;
      font-weight: bold;
    }

    .param-table {
      width: calc(100% - 20px);
      margin-right: 10px;
    }
  }
}
</style>

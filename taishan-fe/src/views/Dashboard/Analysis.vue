<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { useAuthStore } from '@/store/modules/auth'
import { CountTo } from '@/components/CountTo'
import { ElDivider } from 'element-plus'
import { reactive, ref } from 'vue'
import { Echart } from '@/components/Echart'
import { cloneDeep } from 'lodash-es'
import { baseChartOption } from '@/views/Report/chart/chart-data'
import DataService from '@/api/data'

const authStore = useAuthStore()
const dataService = new DataService()

let username: string | undefined = ''
const getUserInfo = () => {
  const res = authStore.getUserInfo
  username = res?.username
}

const getTimeState = () => {
  // 获取当前时间
  let timeNow = new Date()
  // 获取当前小时
  let hours = timeNow.getHours()
  // 设置默认文字
  let state = ``
  // 判断当前时间段
  if (hours >= 0 && hours <= 10) {
    state = `早上好!`
  } else if (hours > 10 && hours <= 14) {
    state = `中午好!`
  } else if (hours > 14 && hours <= 18) {
    state = `下午好!`
  } else if (hours > 18 && hours <= 24) {
    state = `晚上好!`
  }
  return state
}

let data = ref({
  total_plan: 0,
  total_report: 0,
  total_cost: 0
})

let optionsData = reactive<any>(cloneDeep(baseChartOption))

const getData = async () => {
  const res = await dataService.getDashboardData()
  data.value.total_plan = res.total_plan
  data.value.total_report = res.total_report
  data.value.total_cost = res.total_cost
  refreshData(res.graph)
}

const initGraphData = () => {
  optionsData.title.text = '压测次数/费用'
  optionsData.grid = {
    left: '40px',
    right: '60px',
    bottom: '50px'
  }
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
        interval: 'auto'
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
        interval: 'auto'
      }
    }
  ]
  optionsData.series = [
    {
      type: 'line',
      showSymbol: false,
      zlevel: 3
    },
    {
      type: 'bar',
      showSymbol: false,
      yAxisIndex: 1,
      zlevel: 2
    }
  ]
}

const refreshData = (dt) => {
  const sources: Array<any> = []
  //@ts-ignore
  sources.push(dt.column)
  dt.data.forEach((item) => {
    sources.push(item)
  })
  optionsData.dataset.source = sources
}

getUserInfo()
getData()
initGraphData()
</script>

<template>
  <content-wrap style="height: 120px">
    <div class="head">
      <div class="flex items-center">
        <img
          :src="authStore.getUserInfo?.avatar"
          alt=""
          class="w-70px h-70px rounded-[50%] mr-20px"
        />
        <div> {{ getTimeState() }} {{ username }}, 欢迎使用泰山全链路压测平台 </div>
      </div>
      <div class="flex h-70px items-center">
        <div class="px-8px text-right">
          <div class="text-14px text-gray-400 mb-20px">计划总数</div>
          <count-to class="text-20px" :start-val="0" :end-val="data.total_plan" :duration="1500" />
        </div>
        <ElDivider direction="vertical" />
        <div class="px-8px text-right">
          <div class="text-14px text-gray-400 mb-20px">压测次数</div>
          <count-to
            class="text-20px"
            :start-val="0"
            :end-val="data.total_report"
            :duration="1500"
          />
        </div>
        <ElDivider direction="vertical" />
        <div class="px-8px text-right">
          <div class="text-14px text-gray-400 mb-20px">节省费用</div>
          <count-to class="text-20px" :start-val="0" :end-val="data.total_cost" :duration="1500" />
        </div>
      </div>
    </div>
  </content-wrap>
  <content-wrap style="height: 320px; margin-top: 20px">
    <echart :options="optionsData" :height="320" />
  </content-wrap>
</template>

<style scoped lang="less">
.head {
  font-size: 24px;
  font-weight: normal;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>

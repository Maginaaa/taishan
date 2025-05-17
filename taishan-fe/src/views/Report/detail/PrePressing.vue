<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { reactive, ref } from 'vue'
import { ReportDetail } from '@/views/Report/report-data'
import ReportService from '@/api/scene/report'

const reportService = new ReportService()
const props = defineProps({
  reportId: {
    type: Number,
    default: 0
  }
})

let loading = ref(false)
let report_detail = reactive<ReportDetail | any>({})

const press_type_mp = {
  0: '并发模式',
  1: '阶梯式加压',
  2: 'RPS模式',
  3: 'RPS比例模式'
}
let progress = ref(0)

const getReportDetail = async () => {
  loading.value = true
  try {
    const res = await reportService.getReportDetail(props.reportId)
    report_detail = Object.assign(report_detail, res)
  } finally {
    loading.value = false
    getProgress()
  }
}

const getProgress = () => {
  const interval = setInterval(() => {
    const increase = Math.floor(Math.random() * 10) + 15
    const _val = progress.value + increase
    if (_val > 99) {
      clearInterval(interval)
      progress.value = 99
      return
    }
    progress.value = _val
  }, 1000)
}

getReportDetail()
</script>

<template>
  <content-wrap>
    <el-skeleton :loading="loading" :rows="12" animated>
      <el-descriptions :title="report_detail.plan_name">
        <el-descriptions-item label="持续时间">
          {{ report_detail.duration }} 分钟
        </el-descriptions-item>
        <el-descriptions-item label="压测模式">
          <el-tag size="small">
            {{ press_type_mp[report_detail.press_type] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="预估机器数">
          {{ report_detail.machine_num }}
        </el-descriptions-item>
        <el-descriptions-item label="起始并发数">
          {{ report_detail.concurrency }}
        </el-descriptions-item>
        <el-descriptions-item label="备注">
          {{ report_detail.remark }}
        </el-descriptions-item>
      </el-descriptions>
      <el-progress :text-inside="true" :stroke-width="24" :percentage="progress" status="success" />
    </el-skeleton>
  </content-wrap>
</template>

<style scoped lang="less"></style>

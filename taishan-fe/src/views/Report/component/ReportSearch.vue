<script setup lang="ts">
import type { ReportSearchReq } from '@/views/Report/report-data'
import { ref, unref, watch } from 'vue'
import { useRouter } from 'vue-router'
import ReportService from '@/api/scene/report'
import { Icon } from '@/components/Icon'

const router = useRouter()

const reportService = new ReportService()

const search = defineModel<ReportSearchReq | any>({})

const table_loading = ref(false)
let report_list = ref([])
const total = ref(0)
const press_type_list = [
  {
    value: 0,
    label: '并发模式'
  },
  {
    value: 1,
    label: '阶梯式加压'
  },
  {
    value: 2,
    label: 'RPS模式'
  },
  {
    value: 3,
    label: 'RPS比例模式'
  },
  {
    value: 4,
    label: '固定次数'
  }
]
const searchReportList = async () => {
  table_loading.value = true
  try {
    const res = await reportService.getReportList(unref(search))
    report_list.value = res.list
    total.value = res.total
  } finally {
    table_loading.value = false
  }
}

const updateReportName = (row) => {
  row.is_editing = false
  reportService.updateReportName({
    report_id: row.report_id,
    report_name: row.report_name
  })
}

const jumpReportDetail = (reportId) => {
  router.push({
    name: 'ReportDetail',
    params: {
      id: reportId
    }
  })
}

const jumpPlanDetail = (planId) => {
  router.push({
    name: 'PlanDetail',
    params: {
      id: planId
    }
  })
}

watch(
  () => search,
  () => {
    searchReportList()
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <el-table
    v-loading="table_loading"
    :data="report_list"
    resizable
    border
    stripe
    class="search-table"
  >
    <el-table-column prop="report_id" label="报告ID" align="center" header-align="center" />
    <el-table-column label="压测时长" align="center" header-align="center">
      <template #default="scope">
        <el-tag v-if="scope.row.status" type="success" effect="plain"> 运行中 </el-tag>
        <span v-else>{{ scope.row.actual_duration }} 分钟</span>
      </template>
    </el-table-column>
    <el-table-column prop="report_name" label="报告名称" align="center" header-align="center">
      <template #default="scope">
        <el-input
          v-if="scope.row.is_editing"
          v-model="scope.row.report_name"
          maxlength="50"
          show-word-limit
          autofocus
          @blur="updateReportName(scope.row)"
        />
        <div v-else class="report-name">
          <span @click="scope.row.is_editing = true" class="remark">
            {{ scope.row.report_name }}
          </span>
          <Icon icon="ep:edit" class="icon" />
        </div>
      </template>
    </el-table-column>
    <el-table-column prop="plan_name" label="计划" align="center" header-align="center">
      <template #default="{ row }">
        <el-button
          link
          type="primary"
          v-tc="{ name: '查看计划', plan_id: row.plan_id }"
          @click="jumpPlanDetail(row.plan_id)"
        >
          {{ row.plan_name }}
        </el-button>
      </template>
    </el-table-column>
    <el-table-column label="压测模式" align="center" header-align="center">
      <template #default="scope">
        {{ press_type_list.find((item) => item.value === scope.row.press_type)?.label }}
      </template>
    </el-table-column>
    <el-table-column prop="create_user_name" label="创建人" align="center" header-align="center" />
    <el-table-column prop="start_time" label="开始时间" align="center" header-align="center" />
    <el-table-column fixed="right" align="center" header-align="center" label="操作" width="120">
      <template #default="{ row }">
        <el-button
          link
          type="primary"
          @click="jumpReportDetail(row.report_id)"
          v-tc="{ name: '查看报告', report_id: row.report_id }"
        >
          查看
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <div class="pagination-div">
    <el-pagination
      class="report-pag"
      small
      background
      v-model:current-page="search.page"
      v-model:page-size="search.page_size"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next, jumper "
      :total="total"
      @change="searchReportList"
    />
  </div>
</template>

<style scoped lang="less">
.search-table {
  width: 100%;
}
.button-spacing {
  display: inline-block;
  width: 10px;
}

.pagination-div {
  margin-top: 10px;

  .report-pag {
    margin: 10px 0;
    float: right;
  }
}

.report-name {
  .icon {
    display: none;
  }
}

.report-name:hover .icon {
  display: inline;
  color: #409eff;
}
</style>

<script setup lang="ts">
import { ref } from 'vue'
import PlanService from '@/api/scene/plan'
import DebugResult from '@/views/Plan/plan/DebugResult.vue'

const planService = new PlanService()

const props = defineProps({
  planId: Number
})

let table_loading = ref(false)
let debug_record_list = ref([])

let search = ref({
  page: 1,
  page_size: 10
})
let total = ref(0)

const searchRecordList = async () => {
  table_loading.value = true
  try {
    const res = await planService.getPlanDebugRecord({
      plan_id: props.planId,
      page: search.value.page,
      page_size: search.value.page_size
    })
    debug_record_list.value = res.list
    total.value = res.total
  } finally {
    table_loading.value = false
  }
}
const debugResultRef = ref()
const getDebugDetail = (dt) => {
  debugResultRef.value.drawerOpen(dt)
}

searchRecordList()
</script>

<template>
  <el-table
    v-loading="table_loading"
    :data="debug_record_list"
    resizable
    border
    stripe
    class="search-table"
  >
    <el-table-column prop="time" label="调试时间" align="center" header-align="center" />
    <el-table-column label="状态" align="center" header-align="center">
      <template #default="{ row }">
        <el-tag v-if="row.passed" type="success" size="small" effect="plain"> 成功 </el-tag>
        <el-tag v-else type="danger" size="small" effect="plain">失败</el-tag>
      </template>
    </el-table-column>
    <el-table-column fixed="right" align="center" header-align="center" label="操作">
      <template #default="{ row }">
        <el-button link type="primary" @click="getDebugDetail(row.result)"> 查看 </el-button>
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
      @change="searchRecordList"
    />
  </div>
  <debug-result ref="debugResultRef" />
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
</style>

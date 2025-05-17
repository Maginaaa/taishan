<script setup lang="ts">
import { ref } from 'vue'
import OperationLogService from '@/api/scene/operationLog'
import { OperationTypeMap, SourceNameMap } from '@/constants'

const operationLogService = new OperationLogService()

const props = defineProps({
  sourceId: {
    type: Number,
    default: 0
  },
  sourceName: {
    type: String,
    default: ''
  }
})

let table_loading = ref(false)
let log_list = ref([])

let search = ref({
  page: 1,
  page_size: 10
})
let total = ref(0)

const searchOperationLogList = async () => {
  table_loading.value = true
  try {
    const res = await operationLogService.getOperationLog({
      page: search.value.page,
      page_size: search.value.page_size,
      source_id: props.sourceId,
      source_name: props.sourceName
    })
    log_list.value = res.list
    total.value = res.total
  } finally {
    table_loading.value = false
  }
}

searchOperationLogList()
</script>

<template>
  <el-table v-loading="table_loading" :data="log_list" resizable border stripe class="search-table">
    <el-table-column prop="create_time" label="操作时间" align="center" header-align="center" />
    <el-table-column prop="operator_name" label="操作人" align="center" header-align="center" />
    <el-table-column label="内容" align="center" header-align="center">
      <template #default="{ row }">
        {{ SourceNameMap[row.source_name] }}
      </template>
    </el-table-column>
    <el-table-column label="类型" align="center" header-align="center">
      <template #default="{ row }">
        {{ OperationTypeMap[row.operation_type] }}
      </template>
    </el-table-column>
    <el-table-column label="变更字段" align="center" header-align="center">
      <template #default="{ row }">
        <el-popover
          placement="left"
          :width="600"
          trigger="hover"
          v-for="(v, k) in row.value_diff"
          :key="k"
        >
          <template #reference>
            <el-tag type="primary">{{ k }}</el-tag>
          </template>
          <div>
            {{ v }}
          </div>
        </el-popover>
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
      @change="searchOperationLogList"
    />
  </div>
</template>

<style scoped lang="less">
.pagination-div {
  margin-top: 10px;

  .report-pag {
    margin: 10px 0;
    float: right;
  }
}
</style>

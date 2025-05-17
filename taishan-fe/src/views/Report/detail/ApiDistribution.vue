<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  data: {
    type: Object,
    default: () => {
      return {}
    }
  },
  caseNameMap: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const needShow = computed(() => {
  return (
    props.data.error_code_array?.length > 0 ||
    props.data.request_error_array?.length > 0 ||
    props.data.assert_error_array?.length > 0
  )
})
</script>

<template>
  <div class="table-row" v-show="needShow">
    <el-card class="table-item">
      <div class="center-text">错误码</div>
      <el-table :data="props.data.error_code_array" max-height="240" border class="param-table">
        <el-table-column align="center" header-align="center" label="接口名" prop="case_id">
          <template #default="{ row }">
            {{ caseNameMap[row.case_id] || '--' }}
          </template>
        </el-table-column>
        <el-table-column align="center" header-align="center" label="响应码" prop="error_code" />
        <el-table-column align="center" header-align="center" label="次数" prop="count" />
      </el-table>
    </el-card>
    <el-card class="table-item">
      <div class="center-text">请求失败</div>
      <el-table :data="props.data.request_error_array" max-height="240" border class="param-table">
        <el-table-column align="center" header-align="center" label="接口名" prop="case_id">
          <template #default="{ row }">
            {{ caseNameMap[row.case_id] || '--' }}
          </template>
        </el-table-column>
        <el-table-column align="center" header-align="center" label="次数" prop="count" />
      </el-table>
    </el-card>
    <el-card class="table-item">
      <div class="center-text">断言失败</div>
      <el-table :data="props.data.assert_error_array" max-height="240" border class="param-table">
        <el-table-column align="center" header-align="center" label="接口名" prop="case_id">
          <template #default="{ row }">
            {{ caseNameMap[row.case_id] || '--' }}
          </template>
        </el-table-column>
        <el-table-column align="center" header-align="center" label="次数" prop="count" />
      </el-table>
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
    width: 33%;
    height: 280px;
    margin: 0 5px;

    .center-text {
      text-align: center; /* 将文本水平居中 */
      font-size: 15px;
      font-weight: bold;
    }

    .param-table {
      width: calc(100% - 20px);
      margin-top: 3px;
      margin-left: 10px;
      margin-right: 10px;
    }
  }
}
</style>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { PlanDebugResult } from '@/views/Plan/scene/case-data'
import 'splitpanes/dist/splitpanes.css'
import HttpResp from '@/views/Plan/component/HttpResp.vue'
import { cloneDeep } from 'lodash-es'

const draw_visible = ref(false)
const result_drawer_visible = ref(false)
let results = ref<PlanDebugResult[]>([])
let expand_row_keys = ref([])
const drawerOpen = (val: any) => {
  draw_visible.value = true
  results.value = val
  expand_row_keys.value = val.filter((item) => !item.passed).map((item) => item.scene_id)
}

let current_result = reactive({})

const getSamplingDetail = (dt) => {
  result_drawer_visible.value = true
  current_result = reactive(cloneDeep(dt))
}

defineExpose({
  drawerOpen
})
</script>

<template>
  <div>
    <el-drawer direction="rtl" size="85%" v-model="draw_visible" :destroy-on-close="true">
      <template #header>
        <span>执行结果</span>
      </template>
      <el-table
        :data="results"
        style="width: 100%"
        default-expand-all
        row-style="background-color: var(--el-fill-color-light);"
        row-key="scene_id"
        :expand-row-keys="expand_row_keys"
        :border="true"
      >
        <el-table-column type="expand">
          <template #default="props">
            <el-table :data="props.row.case_result_list" style="margin-left: 50px">
              <el-table-column label="状态">
                <template #default="{ row }">
                  <el-tag
                    v-if="row.request_success && row.assert_success"
                    type="success"
                    size="small"
                    effect="plain"
                  >
                    成功
                  </el-tag>
                  <el-tag v-else type="danger" size="small" effect="plain">失败</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="接口名" prop="case_name" />
              <el-table-column label="响应码" prop="status_code" />
              <el-table-column label="操作">
                <template #default="{ row }">
                  <el-button link type="primary" @click="getSamplingDetail(row)">
                    查看详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="{ row }">
            <el-tag v-if="row.passed" type="success" size="small" effect="plain"> 成功 </el-tag>
            <el-tag v-else type="danger" size="small" effect="plain">失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scene_name" label="场景名" />
        <el-table-column prop="total_requests" label="接口数" />
        <el-table-column prop="total_time" label="执行时长">
          <template #default="{ row }">
            <div>{{ row.total_time.toFixed(2) }} s</div>
          </template>
        </el-table-column>
      </el-table>
      <el-drawer size="65%" :with-header="false" v-model="result_drawer_visible">
        <http-resp :result="current_result" />
      </el-drawer>
    </el-drawer>
  </div>
</template>

<style scoped lang="less">
:deep(.el-drawer__header) {
  padding: 10px;
  margin: 0;
  border: 1px solid #cecece;
  color: #0f2438;
}

:deep(.el-drawer__body) {
  padding: 0;
}

:deep(.el-card__body) {
  padding: 10px;
  margin-right: 10px;
}

:deep(.el-tabs__header) {
  margin-bottom: 0;
}

:deep(.el-tab-pane) {
  height: calc(100vh - 80px);
}

:deep(.el-tabs__item) {
  padding: 0 8px;
}

.tabs {
  height: 100%;
  padding: 0 10px;

  :deep(.splitpanes__splitter) {
    width: 0;
  }
}

.result-card-container {
  padding: 10px 10px 0 0;
  overflow-y: auto;

  .result-card {
    margin-bottom: 6px;
    width: 100%;
    height: 40px;
    cursor: pointer;
    margin-right: 10px;

    .result-row {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .case-info {
        padding: 0;
        width: calc(100% - 20px);
        display: flex;
        align-items: center;

        .case-name {
          padding-left: 6px;
          color: #1890ff;
          font-weight: 800;
          font-size: 14px;

          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }

      .assert-success-btn {
        color: #67c23a;
      }

      .assert-fail-btn {
        color: #f56c6c;
      }
    }
  }
}
</style>

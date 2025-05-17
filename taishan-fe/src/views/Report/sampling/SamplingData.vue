<script setup lang="ts">
import { reactive, ref, unref } from 'vue'
import { cloneDeep } from 'lodash-es'
import HttpResp from '@/views/Plan/component/HttpResp.vue'
import { ElMessage } from 'element-plus'
import ReportService from '@/api/report'
import { isNumber } from '@/utils/is'
import DataService from '@/api/data'
import TraceDetail from '@/views/Report/component/TraceDetail.vue'

const reportService = new ReportService()
const dataService = new DataService()

let visible = ref(false)
let result_drawer_visible = ref(false)
const base_search = {
  report_id: 0,
  scene_id: 0,
  case_id_list: [],
  page: 1,
  page_size: 10,
  start_time: undefined,
  end_time: undefined,
  min_rt: undefined,
  max_rt: undefined,
  code_list: [],
  status_list: [1, 2]
}
let search = reactive(cloneDeep(base_search))
let total = ref(0)
let apiList = ref<any>([])
let sampling_data = ref([])
let current_result = reactive({})
let table_loading = ref(false)

let chooseAllCode = ref<boolean>(false)
const codeList = [200, 301, 302, 400, 401, 405, 415, 429, 500, 502, 503, 504]

const statusList = [
  { label: '请求失败', value: 1 },
  { label: '断言失败', value: 2 }
]

const drawerOpen = (data, cs_list) => {
  apiList.value = cs_list
  visible.value = true
  search = reactive(cloneDeep(base_search))
  search.report_id = data.report_id
  search.scene_id = data.scene_id
  searchSamplingData()
}

const chooseCodeChange = (val) => {
  chooseAllCode.value = val.length === codeList.length
}
const handleChooseAllCodeList = (val) => {
  if (val) {
    search.code_list = [...codeList]
  } else {
    search.code_list = []
  }
}

const searchSamplingData = async () => {
  const _search = unref(search)
  if (_search.min_rt === '') {
    delete _search.min_rt
  }
  if (_search.max_rt === '') {
    delete _search.max_rt
  }
  if (_search.min_rt !== undefined && _search.max_rt !== undefined) {
    if (_search.min_rt >= _search.max_rt) {
      ElMessage.warning('请输入正确的RT范围')
      return
    }
    if (!isNumber(_search.min_rt) || !isNumber(_search.max_rt)) {
      ElMessage.warning('请输入正确的RT范围')
      return
    }
  }
  table_loading.value = true
  try {
    const data = await reportService.getSamplingData(_search)
    sampling_data.value = data.list
    total.value = data.total
  } finally {
    table_loading.value = false
  }
}

const getSamplingDetail = (dt) => {
  result_drawer_visible.value = true
  current_result = reactive(cloneDeep(dt))
}
const drawerClose = () => {
  visible.value = false
  search = reactive(cloneDeep(base_search))
}

const getBffResponseTraceID = (row) => {
  return row.http_response_data.response_header.find((item) => item.key === 'Traceid').value || ''
}

const trace_detail = ref(null)
const current_trace_id = ref('')
const inner_visible = ref(false)
const getTraceDetail = async (traceId) => {
  current_trace_id.value = traceId
  const res = await dataService.getTraceDetail({
    trace_id: traceId,
    report_id: search.report_id
  })
  trace_detail.value = res
  inner_visible.value = true
}

defineExpose({
  drawerOpen
})
</script>

<template>
  <div class="container">
    <el-drawer direction="rtl" size="85%" v-model="visible" :before-close="drawerClose">
      <template #header>
        <span>查看日志</span>
      </template>
      <div style="padding: 0 20px 0 20px">
        <el-form :inline="true" :model="search" class="demo-form-inline">
          <el-form-item label="压测API">
            <el-select
              v-model="search.case_id_list"
              multiple
              collapse-tags
              collapse-tags-tooltip
              filterable
              placeholder="请选择接口"
              style="width: 240px"
            >
              <el-option
                v-for="item in apiList"
                :key="item.id"
                :label="item.value"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="RT范围">
            <div style="display: flex">
              <el-input
                v-model.number="search.min_rt"
                placeholder="最小RT"
                onchange="if(value<0)value=0;if(value>5000)value=5000"
                style="width: 150px"
              />
              <span style="padding: 0 8px 0 8px">-</span>
              <el-input
                v-model.number="search.max_rt"
                placeholder="最大RT"
                onchange="if(value<0)value=0;if(value>5000)value=5000"
                style="width: 150px"
              />
            </div>
          </el-form-item>
          <el-form-item label="时间范围">
            <div style="display: flex">
              <el-date-picker
                v-model="search.start_time"
                value-format="YYYY-MM-DD HH:mm:ss"
                type="datetime"
                placeholder="开始时间"
              />
              <span style="padding: 0 8px 0 8px">-</span>
              <el-date-picker
                v-model="search.end_time"
                value-format="YYYY-MM-DD HH:mm:ss"
                type="datetime"
                placeholder="结束时间"
              />
            </div>
          </el-form-item>
          <el-form-item label="响应码">
            <el-select
              v-model="search.code_list"
              multiple
              clearable
              collapse-tags
              :max-collapse-tags="1"
              style="width: 240px"
              @change="chooseCodeChange"
            >
              <template #header>
                <el-checkbox v-model="chooseAllCode" @change="handleChooseAllCodeList">
                  全选
                </el-checkbox>
              </template>
              <el-option v-for="item in codeList" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-checkbox-group v-model="search.status_list">
              <el-checkbox
                v-for="item in statusList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-checkbox-group>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="searchSamplingData">查询</el-button>
          </el-form-item>
        </el-form>
        <el-table :data="sampling_data" v-loading="table_loading" style="width: 100%">
          <el-table-column prop="result" label="结果" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.success" type="success" size="small">成功</el-tag>
              <el-tag v-else type="danger" size="small">失败</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="TraceID" width="300">
            <template #default="{ row }">
              <el-button
                type="primary"
                link
                @click="getTraceDetail(getBffResponseTraceID(row))"
                style="font-weight: bold"
              >
                {{ getBffResponseTraceID(row) }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column prop="start_time" label="时间" />
          <el-table-column prop="case_name" label="接口名" />
          <el-table-column prop="response_time" label="RT(ms)" />
          <el-table-column prop="status_code" label="响应码" />
          <el-table-column fixed="right" label="操作" width="200">
            <template #default="scope">
              <el-button
                link
                type="primary"
                @click="getSamplingDetail(scope.row.http_response_data)"
              >
                查看详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-div">
          <el-pagination
            background
            v-model:current-page="search.page"
            v-model:page-size="search.page_size"
            :page-sizes="[10, 20, 50]"
            layout="total, prev, pager, next, jumper "
            :total="total"
            style="margin-top: 10px; margin-right: 20px; float: right"
            @change="searchSamplingData"
          />
        </div>
      </div>
      <el-drawer size="65%" :with-header="false" v-model="result_drawer_visible">
        <http-resp :result="current_result" />
      </el-drawer>
    </el-drawer>
    <el-dialog v-model="inner_visible" width="60%" :title="current_trace_id">
      <trace-detail :trace-detail="trace_detail" />
    </el-dialog>
  </div>
</template>

<style scoped lang="less">
:deep(.el-drawer__body) {
  padding: 10px 0 0 0;
}

.container {
  :deep(.el-drawer__header) {
    height: 20px;
    margin: 0;
  }

  .search-form {
    padding: 0 20px 0 20px;
  }

  .pagination-div {
    margin-top: 10px;
  }
}
</style>

<script setup lang="ts">
import { JsonEditor } from '@/components/JsonEditor'
import { ref, watch } from 'vue'

const props = defineProps({
  result: {
    type: Object,
    default: () => ({})
  },
  readOnly: {
    type: Boolean,
    default: true
  }
})

const active_tab = 'body'
const checking_rule_list: Array<any> = [
  {
    value: 1,
    label: '等于',
    disabled: false
  },
  {
    value: 2,
    label: '不等于'
  },
  {
    value: 3,
    label: '包含'
  },
  {
    value: 4,
    label: '不包含'
  },
  {
    value: 5,
    label: '大于'
  },
  {
    value: 6,
    label: '小于'
  },
  {
    value: 7,
    label: '大于等于'
  },
  {
    value: 8,
    label: '小于等于'
  }
]

const convertToSeconds = (val) => {
  if (val >= 1000) {
    return (val / 1000).toFixed(2) + 's'
  } else {
    return val + 'ms'
  }
}

const convertByte = (val) => {
  if (!val) {
    return 0
  }
  if (val >= 1024 && val < 1024 ** 2) {
    return (val / 1024).toFixed(2) + 'MB'
  } else if (val >= 1024 ** 2 && val < 1024 ** 3) {
    return (val / 1024 ** 2).toFixed(2) + 'GB'
  } else if (val >= 1024 ** 3) {
    return (val / 1024 ** 3).toFixed(2) + 'TB'
  } else {
    return val.toFixed(2) + 'KB'
  }
}

const getStatusCodeClass = () => {
  const _code = props.result.status_code
  if (_code >= 200 && _code < 300) {
    return 'green'
  } else if (_code >= 300 && _code < 400) {
    return 'orange'
  } else {
    return 'red'
  }
}
const getResponseTimeClass = () => {
  const _ms = props.result.response_time
  if (_ms < 1000) {
    return 'green'
  } else if (_ms >= 1000 && _ms < 10000) {
    return 'orange'
  } else {
    return 'red'
  }
}
let request_json = ref<Object | null>()
let response_json = ref<Object | string>()
watch(
  () => props.result,
  (newVal) => {
    try {
      response_json.value = JSON.parse(newVal.response_content)
    } catch (e) {
      response_json.value = newVal.response_content
    }
    try {
      request_json.value = JSON.parse(newVal.request_content)
    } catch (e) {
      request_json.value = newVal.request_content
    }
  },
  { immediate: true, deep: true }
)

const valueExtract = (type, node) => {
  emits('value-extract', type, node)
}
const emits = defineEmits(['value-extract'])
</script>

<template>
  <div v-if="Object.keys(props.result).length > 0" class="http-resp">
    <el-tabs v-model="active_tab" class="tabs">
      <el-tab-pane label="响应Body" name="body">
        <div class="body-container">
          <json-editor
            v-model="response_json"
            :read-only="readOnly"
            @value-extract="valueExtract"
          />
        </div>
      </el-tab-pane>
      <el-tab-pane
        v-if="props.result && props.result.response_header"
        label="响应Header"
        name="response_header"
      >
        <template #label>
          <span> 响应header </span>
        </template>
        <el-table :key="active_tab" border stripe :data="props.result.response_header">
          <el-table-column prop="key" label="Key" />
          <el-table-column prop="value" label="Value" />
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="请求信息" name="request">
        <el-row>
          <el-col :span="2">
            <div class="method">
              {{ props.result.method_type }}
            </div>
          </el-col>
          <el-col :span="20">
            <div class="url">
              {{ props.result.url }}
            </div>
          </el-col>
        </el-row>
        <el-row>
          <json-editor v-show="request_json !== ''" :read-only="true" v-model="request_json" />
        </el-row>
      </el-tab-pane>
      <el-tab-pane
        v-if="props.result && props.result.request_header"
        label="请求Header"
        name="request_header"
      >
        <template #label>
          <span> 请求header </span>
        </template>
        <el-table :key="active_tab" border stripe :data="props.result.request_header">
          <el-table-column prop="key" label="Key" />
          <el-table-column prop="value" label="Value" />
        </el-table>
      </el-tab-pane>
      <el-tab-pane name="assert">
        <template #label>
          <span>
            断言结果
            <el-badge
              v-show="props.result.assert_res.length > 0"
              class="badge"
              :type="props.result.assert_success ? 'primary' : 'danger'"
              :value="props.result.assert_res.length"
            />
          </span>
        </template>
        <el-table :key="active_tab" border stripe :data="props.result.assert_res">
          <el-table-column prop="extract_rule" label="提取规则" />
          <el-table-column prop="extract_value" label="实际值" />
          <el-table-column prop="expected_value" label="预期值">
            <template #default="scope">
              <span>
                {{
                  checking_rule_list.find((item) => item.value === scope.row.checking_rule).label
                }}
              </span>
              <span>{{ scope.row.expected_value }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="result" label="断言结果">
            <template #default="scope">
              <el-tag v-if="scope.row.assert_pass" type="success" size="small">成功</el-tag>
              <el-tag v-else type="danger" size="small">失败</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane name="variable">
        <template #label>
          <span>
            参数提取
            <el-badge
              v-show="props.result.variable_res.length > 0"
              class="badge"
              :type="props.result.extract_all_success ? 'primary' : 'danger'"
              :value="props.result.variable_res.length"
            />
          </span>
        </template>
        <el-table :key="active_tab" border stripe :data="props.result.variable_res">
          <el-table-column prop="extract_rule" label="提取规则" />
          <el-table-column prop="variable_name" label="变量名" />
          <el-table-column prop="actual_res" label="提取值">
            <template #default="scope">
              <el-tag v-if="!scope.row.extract_success" type="danger" size="small">提取失败</el-tag>
              <span v-else>{{ scope.row.actual_res }}<span></span> </span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
    <div class="response-info">
      <div class="info-col">
        状态码：
        <span :class="getStatusCodeClass()">
          {{ props.result.status_code }}
        </span>
      </div>
      <div class="info-col">
        耗时：
        <span :class="getResponseTimeClass()">
          {{ convertToSeconds(props.result.response_time) }}
        </span>
      </div>
      <div class="info-col">
        大小：
        <el-dropdown :hide-on-click="false" trigger="click">
          <span class="green el-dropdown-link response-total-size">
            {{
              convertByte(
                props.result.response_size.body_size + props.result.response_size.header_size
              )
            }}
          </span>
          <template #dropdown>
            <el-dropdown-menu style="width: 200px">
              <el-dropdown-item>
                <div class="response-size">
                  <div>Body</div>
                  <span style="color: #13ce66">
                    {{ convertByte(props.result.response_size.body_size) }}
                  </span>
                </div>
              </el-dropdown-item>
              <el-dropdown-item>
                <div class="response-size">
                  <div>Header</div>
                  <span style="color: #13ce66">
                    {{ convertByte(props.result.response_size.header_size) }}
                  </span>
                </div>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.http-resp {
  display: flex;
  flex-direction: column;
  height: 100%;

  .tabs {
    height: 100%;
    padding: 0 10px 0 10px;

    :deep(.el-tabs__header) {
      margin-bottom: 5px;
    }
    :deep(.el-tabs__item) {
      padding: 0 12px;
      font-weight: normal;
    }
    :deep(.el-tabs__content) {
      height: calc(100% - 40px);
      overflow-y: auto;
    }
    :deep(.el-radio) {
      margin-right: 16px;
    }
    :deep(.el-radio__label) {
      font-weight: normal;
    }
    :deep(.is-active) {
      font-weight: bold;
    }

    .body-container {
      //display: flex;
      overflow-y: auto;
    }

    .badge {
      margin-top: 3px;
    }

    .kv-row {
      margin-top: 6px;
      .kv-checkbox {
        width: 20px;
        margin-right: 10px;
      }

      .kv-delete {
        width: 60px;
        margin-right: 10px;

        .icon-button {
          color: #409eff;
          cursor: pointer;
        }
      }

      .kv-select {
        width: 500px;
      }
    }

    .method {
      display: flex;
      color: #e6a23c;
      font-weight: 500;
      vertical-align: middle;
      text-align: center;
      justify-content: center;
      align-items: center;
      border-radius: 4px;
      min-height: 36px;
      font-size: 16px;
    }

    .url {
      display: flex;
      color: #383737;
      align-items: center;
      border-radius: 4px;
      min-height: 36px;
      font-size: 15px;
    }
  }

  .response-info {
    position: absolute;
    right: 20px;
    height: 40px;
    display: flex;
    align-items: center;

    .info-col {
      margin-right: 10px;
      font-weight: normal;
      font-size: 14px;

      .response-total-size {
        padding: 3px 0 3px 0;
      }

      .response-size {
        display: flex;
        justify-content: space-between;
        width: 100%;
      }

      .green {
        color: #13ce66;
      }

      .orange {
        color: #e6a23c;
      }

      .red {
        color: #f56c6c;
      }
    }
  }
}
</style>

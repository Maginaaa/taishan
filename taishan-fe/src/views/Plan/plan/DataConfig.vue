<script setup lang="ts">
import { ParamsForm, PlanInfo } from '@/views/Plan/plan-data'
import { ref, unref } from 'vue'
import PlanService from '@/api/scene/plan'
import DataSourceManage from '@/views/Plan/plan/DataSourceManage.vue'

const plan_info = defineModel<PlanInfo | any>()
const planService = new PlanService()

// 高级配置
let active_tab = ref('variable_file')

const eleIsDisable = (index: number, data: any[]) => {
  return data.length - 1 === index
}
const formChange = (val: ParamsForm[]) => {
  if (notBlankLines(val) >= val.length) {
    val.push({ enable: true, key: '', value: '', desc: '' })
  }
  planService.updatePlan(unref(plan_info))
}
const notBlankLines = (array: ParamsForm[]) => {
  return array.filter((item) => item.key !== '' || item.value !== '').length | 0
}
const removeParam = (index: number) => {
  plan_info.value.global_variable.splice(index, 1)
  planService.updatePlan(unref(plan_info))
}

const removeHeader = (index: number) => {
  plan_info.value.default_header.splice(index, 1)
  planService.updatePlan(unref(plan_info))
}

const header_suggestions = [
  { value: 'Authorization' },
  { value: 'Accept' },
  { value: 'Accept-Encoding' },
  { value: 'Accept-Language' },
  { value: 'Content-Type' },
  { value: 'Connection' },
  { value: 'Cookie' },
  { value: 'Host' },
  { value: 'User-Agent' }
]

const headerSuggestions = (queryString, cb) => {
  const res = queryString
    ? header_suggestions.filter((item) => item.value.toLowerCase().match(queryString.toLowerCase()))
    : header_suggestions
  cb(res)
}
</script>

<template>
  <el-tabs v-model="active_tab">
    <el-tab-pane label="参数文件" name="variable_file" class="variable-file-container">
      <data-source-manage :plan-id="plan_info.plan_id" />
    </el-tab-pane>
    <el-tab-pane label="全局变量" name="global_variable">
      <el-row
        v-for="(item, index) in plan_info.global_variable"
        class="kv-row"
        :key="index"
        :gutter="6"
      >
        <el-col :span="0.5" class="kv-checkbox">
          <el-checkbox
            v-show="!eleIsDisable(index, plan_info.global_variable)"
            v-model="item.enable"
          />
        </el-col>
        <el-col :span="4">
          <el-input
            v-model="item.key"
            :disabled="!eleIsDisable(index, plan_info?.global_variable) && !item.enable"
            clearable
            placeholder="变量名"
            style="width: 100%"
            @change="formChange(plan_info.global_variable)"
          />
        </el-col>
        <el-col :span="4">
          <el-input
            v-model="item.value"
            :disabled="!eleIsDisable(index, plan_info.global_variable) && !item.enable"
            clearable
            placeholder="变量值"
            @change="formChange(plan_info.global_variable)"
          />
        </el-col>
        <el-col :span="4">
          <el-input
            v-model="item.desc"
            :disabled="!eleIsDisable(index, plan_info.global_variable) && !item.enable"
            clearable
            placeholder="备注"
            @change="formChange(plan_info.global_variable)"
          />
        </el-col>
        <el-col :span="0.5" class="kv-delete">
          <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
            <Icon
              v-show="!eleIsDisable(index, plan_info.global_variable)"
              icon="ep:delete"
              class="icon-button"
              @click.stop="removeParam(index)"
            />
          </el-tooltip>
        </el-col>
      </el-row>
    </el-tab-pane>
    <el-tab-pane label="默认请求头" name="default_header">
      <el-row
        v-for="(item, index) in plan_info.default_header"
        class="kv-row"
        :key="index"
        :gutter="6"
      >
        <el-col :span="0.5" class="kv-checkbox">
          <el-checkbox
            v-show="!eleIsDisable(index, plan_info.default_header)"
            v-model="item.enable"
          />
        </el-col>
        <el-col :span="4">
          <el-autocomplete
            v-model="item.key"
            :disabled="!eleIsDisable(index, plan_info.default_header) && !item.enable"
            clearable
            :fetch-suggestions="headerSuggestions"
            placeholder="键"
            style="width: 100%"
            @select="formChange(plan_info.headers_form)"
            @change="formChange(plan_info.headers_form)"
          />
        </el-col>
        <el-col :span="4">
          <el-input
            v-model="item.value"
            :disabled="!eleIsDisable(index, plan_info.default_header) && !item.enable"
            clearable
            placeholder="值"
            @change="formChange(plan_info.default_header)"
          />
        </el-col>
        <el-col :span="4">
          <el-input
            v-model="item.desc"
            :disabled="!eleIsDisable(index, plan_info.default_header) && !item.enable"
            clearable
            placeholder="备注"
            @change="formChange(plan_info.default_header)"
          />
        </el-col>
        <el-col :span="0.5" class="kv-delete">
          <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
            <Icon
              v-show="!eleIsDisable(index, plan_info.default_header)"
              icon="ep:delete"
              class="icon-button"
              @click.stop="removeHeader(index)"
            />
          </el-tooltip>
        </el-col>
      </el-row>
    </el-tab-pane>
  </el-tabs>
</template>

<style scoped lang="less">
.variable-file-container {
  height: calc(100% - 60px);
  // y轴滚动
  overflow-y: auto;
}

.kv-row {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  padding: 3px 0 5px 0;
  width: calc(100% - 1px);

  .kv-checkbox {
    min-width: 20px;
    margin: 0 5px 0 5px;
  }

  .kv-delete {
    min-width: 20px;

    .icon-button {
      color: #409eff;
      cursor: pointer;
      margin-left: 3px;
    }
  }
}
</style>

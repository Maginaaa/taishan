<script setup lang="ts">
import { nextTick, reactive, ref, unref } from 'vue'
import { CaseType } from '@/constants'
import { CaseData, HttpCaseExtend } from '@/views/Plan/scene/case-data'
import { ElMessage } from 'element-plus'
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'
import * as monaco from 'monaco-editor'
import MonacoEditor from '@/views/Components/Editor/MonacoEditor.vue'
import CaseService from '@/api/scene/case'
import HttpResp from '@/views/Plan/component/HttpResp.vue'
import { cloneDeep } from 'lodash-es'
import { curlParse } from '@/utils/curlParse'

const caseService = new CaseService()

const emit = defineEmits(['refresh'])
function refreshPlanScene() {
  emit('refresh')
}

const drawer_visible = ref(false)
const name_edit = ref(false)
let edit_type = ref(-1)
let case_info = reactive<HttpCaseExtend>({})
let scene_id = ref(0)
let case_id = ref<number>(0)

const drawerOpen = (action: number, case_dt: CaseData) => {
  drawer_visible.value = true
  edit_type.value = action
  scene_id.value = case_dt.scene_id
  case_id.value = case_dt.case_id
  case_info = reactive<HttpCaseExtend>(dataCompatibility(case_dt.extend))
}

const drawerClose = () => {
  drawer_visible.value = false
  debug_result = {}
  // done()
}

const popover_visible = ref(false)
const curl = ref('')
const parse = () => {
  const res = curlParse(curl.value)
  if (res) {
    case_info.url = res.url
    case_info.method_type = res.method_type
    case_info.body = res.body
    case_info.headers_form = res.headers_form
    popover_visible.value = false
  } else {
    ElMessage.error('curl格式有误')
  }
}

// header
const caseNameRef = ref()
const caseNameChange = () => {
  if (case_info.case_name.trim() === '') {
    ElMessage.error('接口名称不能为空')
    return
  }
  name_edit.value = false
}
const editCaseName = () => {
  name_edit.value = true
  nextTick(() => {
    caseNameRef.value.focus()
  })
}

// request
const updateFromUrl = () => {
  const parts = case_info.url.split('?')
  if (parts.length <= 1) {
    case_info.params_form = [
      {
        enable: true,
        key: '',
        value: '',
        desc: ''
      }
    ]
    return
  }
  const params = parts[1].split('&')
  case_info.params_form = params.map((param) => {
    const [key, value] = param.split('=')
    return { enable: true, key: key || '', value: value || '' }
  })
  formChange(case_info.params_form)
}
const formChange = (val: any[]) => {
  if (notBlankLines(val) >= val.length) {
    val.push({ ...blank_variable })
  }
}
const method_type_list = ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH', 'HEAD', 'TRACE']
const extract_type_list = [
  {
    value: 1,
    label: 'JsonPath',
    disabled: false
  },
  {
    value: 2,
    label: '正则表达式',
    disabled: true
  },
  {
    value: 3,
    label: 'XPath',
    disabled: true
  }
]
const checking_rule_list = [
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
const caseDebug = async () => {
  request_loading.value = true
  try {
    case_info.scene_id = unref(scene_id)
    const res = await caseService.caseDebug(case_info)
    debug_result = reactive(cloneDeep(res))
  } finally {
    request_loading.value = false
  }
}
const creat_btn_loading = ref(false)
const operateCase = async () => {
  if (case_info.case_name.trim() === '') {
    ElMessage.error('接口名称不能为空')
    return
  }
  case_info.url = case_info.url.trim()
  if (case_info.case_name === '新的HTTP请求' && case_info.url.length > 0) {
    case_info.case_name = case_info.url.slice(0, 40)
  }
  if (case_info.body === undefined || case_info.body === null) {
    case_info.body = {
      body_type: 'json',
      body_value: ''
    }
  }
  creat_btn_loading.value = true
  try {
    edit_type.value > 0
      ? await caseService.updateSceneCase({
          case_id: unref(case_id),
          scene_id: unref(scene_id),
          extend: case_info
        })
      : await caseService.createSceneCase({
          scene_id: unref(scene_id),
          case_name: case_info.case_name,
          type: CaseType.HttpCase,
          extend: case_info
        })
  } finally {
    creat_btn_loading.value = false
  }

  drawerClose()
  refreshPlanScene()
}

let active_tab = ref('body')
const notBlankLines = (array) => {
  if (array === undefined) {
    return
  }
  return array.filter((item) => item.key !== '' || item.value !== '').length
}
const eleIsDisable = (index: number, data: any[]) => {
  return data.length - 1 === index
}
const updateFromParams = () => {
  const params = case_info.params_form
    .filter((param) => param.enable && (param.key || param.value))
    .map((param) => `${param.key}=${param.value}`)
    .join('&')
  case_info.url = case_info.url.split('?')[0] + (params ? '?' + params : '')
}
const removeParam = (index: number) => {
  case_info.params_form.splice(index, 1)
  updateFromParams()
}
const header_suggestions = [
  { value: 'Accesstoken' },
  { value: 'appVersion' },
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
const assertFormChange = (form) => {
  if (assertNotBlankLines(form) >= form.length) {
    form.push({ ...blank_assert })
  }
}
const assertNotBlankLines = (array) => {
  if (array === undefined) {
    return
  }
  return array.filter((item) => item.extract_express !== '' || item.expect_value !== '').length
}
const body_type_list = ['none', 'json', 'file']
const editorMounted = (editor: monaco.editor.IStandaloneCodeEditor) => {
  console.log('editor实例加载完成')
}

// 新字段兼容/容错
const dataCompatibility = (data: HttpCaseExtend) => {
  if (data.overtime_config === undefined) {
    data.overtime_config = {
      enable: false,
      timeout_period: 0
    }
  }
  return cloneDeep(data)
}

// response
const request_loading = ref(false)
let debug_result = reactive({})

const blank_assert = {
  enable: true,
  extract_type: 1,
  extract_express: '',
  checking_rule: 1,
  expect_value: '',
  desc: ''
}
const blank_variable = {
  enable: true,
  extract_type: 1,
  key: '',
  value: ''
}
const valueExtract = (type, node) => {
  if (type === 1) {
    active_tab.value = 'assert'
    const _assert = { ...blank_assert }
    _assert.extract_express = node.path
    _assert.expect_value = String(node.content)
    case_info.assert_form.splice(case_info.assert_form.length - 1, 1)
    case_info.assert_form.push(_assert)
    case_info.assert_form.push({ ...blank_assert })
  }
  if (type === 2) {
    active_tab.value = 'variable'
    const _variable = { ...blank_variable }
    _variable.key = node.path
    const arr = node.path.split('.')
    const _value = arr[arr.length - 1].replace(/([A-Z](?!ID))/g, '_$1').toUpperCase()
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    let suffix = ''
    for (let i = 0; i < 6; i++) {
      suffix += characters.charAt(Math.floor(Math.random() * characters.length))
    }
    _variable.value = _value + '_' + suffix
    case_info.variable_form.splice(case_info.variable_form.length - 1, 1)
    case_info.variable_form.push(_variable)
    case_info.variable_form.push({ ...blank_variable })
  }
}

defineExpose({
  drawerOpen
})
</script>

<template>
  <div class="container">
    <el-drawer
      direction="rtl"
      size="70%"
      :destroy-on-close="true"
      v-model="drawer_visible"
      :before-close="drawerClose"
    >
      <template #header>
        <div>
          <el-input
            v-if="name_edit"
            ref="caseNameRef"
            v-model="case_info.case_name"
            maxlength="40"
            show-word-limit
            clearable
            class="name-edit-input"
            @blur="caseNameChange"
          />
          <span v-else class="name-span" @click="editCaseName">{{ case_info.case_name }}</span>
          <el-popover :visible="popover_visible" placement="bottom-start" width="500">
            <div>
              <el-input v-model="curl" type="textarea" :rows="10" placeholder="请输入curl内容" />
              <el-button
                size="small"
                style="margin-top: 6px"
                type="primary"
                :disabled="!curl"
                @click="parse"
              >
                解析
              </el-button>
              <el-button
                size="small"
                text
                style="margin-top: 6px; margin-left: 6px"
                @click="popover_visible = false"
              >
                取消
              </el-button>
            </div>
            <template #reference>
              <el-tag
                class="operate-col"
                type="success"
                @click="popover_visible = true"
                style="margin-left: 10px"
              >
                Curl自动解析
              </el-tag>
            </template>
          </el-popover>
        </div>
      </template>
      <div class="request-detail">
        <div class="full-url">
          <el-input
            v-model="case_info.url"
            placeholder="请输入url"
            class="url-input"
            size="small"
            @input="updateFromUrl"
          >
            <template #prepend>
              <el-select v-model="case_info.method_type" class="method-select">
                <el-option
                  v-for="item in method_type_list"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </template>
          </el-input>
          <div class="url-btns">
            <el-button
              type="primary"
              @click="caseDebug"
              v-tc="{ name: 'case调试', case_id }"
              class="send-btn"
            >
              发 送
            </el-button>
            <el-button
              :loading="creat_btn_loading"
              @click="operateCase"
              v-tc="{ name: edit_type > 0 ? 'case修改' : 'case创建', case_id }"
            >
              {{ edit_type > 0 ? '修 改' : '创 建' }}
            </el-button>
          </div>
        </div>
        <splitpanes class="default-theme" horizontal>
          <pane min-size="20" size="40" style="background-color: #fff">
            <el-tabs v-model="active_tab" class="tabs">
              <el-tab-pane name="params">
                <template #label>
                  <span>
                    Params
                    <el-tooltip class="item-tabs" effect="dark" placement="top">
                      <template #content>
                        <div>url携带的参数，将以?key1=value1&key2=value2的形式拼接在url上</div>
                      </template>
                      <el-badge
                        v-show="notBlankLines(case_info.params_form) > 0"
                        class="badge"
                        type="primary"
                        :value="notBlankLines(case_info.params_form)"
                      />
                    </el-tooltip>
                  </span>
                </template>
                <el-row
                  v-for="(item, index) in case_info.params_form"
                  class="kv-row"
                  :key="index"
                  :gutter="6"
                >
                  <el-col :span="0.5" class="kv-checkbox">
                    <el-checkbox
                      v-show="!eleIsDisable(index, case_info.params_form)"
                      v-model="item.enable"
                      @change="updateFromParams"
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-input
                      v-model="item.key"
                      :disabled="!eleIsDisable(index, case_info.params_form) && !item.enable"
                      clearable
                      placeholder="键"
                      @change="formChange(case_info.params_form)"
                      @input="updateFromParams"
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-input
                      v-model="item.value"
                      :disabled="!eleIsDisable(index, case_info.params_form) && !item.enable"
                      clearable
                      placeholder="值"
                      @change="formChange(case_info.params_form)"
                      @input="updateFromParams"
                    />
                  </el-col>
                  <el-col :span="7">
                    <el-input
                      v-model="item.desc"
                      :disabled="!eleIsDisable(index, case_info.params_form) && !item.enable"
                      clearable
                      placeholder="备注"
                      @change="formChange(case_info.params_form)"
                    />
                  </el-col>
                  <el-col :span="0.5" class="kv-delete">
                    <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
                      <Icon
                        v-show="!eleIsDisable(index, case_info.params_form)"
                        icon="ep:delete"
                        class="icon-button"
                        @click.stop="removeParam(index)"
                      />
                    </el-tooltip>
                  </el-col>
                </el-row>
              </el-tab-pane>
              <el-tab-pane v-if="case_info.body" label="Body" name="body" style="height: 100%">
                <el-radio-group v-model="case_info.body.body_type" style="margin-left: 10px">
                  <el-radio
                    v-for="item in body_type_list"
                    :key="item"
                    :value="item"
                    :disabled="item !== 'json'"
                  >
                    {{ item }}
                  </el-radio>
                </el-radio-group>
                <monaco-editor
                  class="monaco-container"
                  v-model="case_info.body.body_value"
                  language="json"
                  @editor-mounted="editorMounted"
                />
              </el-tab-pane>
              <el-tab-pane name="head">
                <template #label>
                  <span>
                    Headers
                    <el-badge
                      v-show="notBlankLines(case_info.headers_form) > 0"
                      class="badge"
                      type="primary"
                      :value="notBlankLines(case_info.headers_form)"
                    />
                  </span>
                </template>
                <el-row
                  v-for="(item, index) in case_info.headers_form"
                  class="kv-row"
                  :key="index"
                  :gutter="6"
                >
                  <el-col :span="0.5" class="kv-checkbox">
                    <el-checkbox
                      v-show="!eleIsDisable(index, case_info.headers_form)"
                      v-model="item.enable"
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-autocomplete
                      v-model="item.key"
                      :disabled="!eleIsDisable(index, case_info.headers_form) && !item.enable"
                      clearable
                      :fetch-suggestions="headerSuggestions"
                      placeholder="键"
                      style="width: 100%"
                      @select="formChange(case_info.headers_form)"
                      @change="formChange(case_info.headers_form)"
                    />
                  </el-col>
                  <el-col :span="8">
                    <el-input
                      v-model="item.value"
                      :disabled="!eleIsDisable(index, case_info.headers_form) && !item.enable"
                      clearable
                      placeholder="值"
                      @change="formChange(case_info.headers_form)"
                    />
                  </el-col>
                  <el-col :span="7">
                    <el-input
                      v-model="item.desc"
                      :disabled="!eleIsDisable(index, case_info.headers_form) && !item.enable"
                      clearable
                      placeholder="备注"
                      @change="formChange(case_info.headers_form)"
                    />
                  </el-col>
                  <el-col :span="0.5" class="kv-delete">
                    <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
                      <Icon
                        v-show="!eleIsDisable(index, case_info.headers_form)"
                        icon="ep:delete"
                        class="icon-button"
                        @click.stop="case_info.headers_form.splice(index, 1)"
                      />
                    </el-tooltip>
                  </el-col>
                </el-row>
              </el-tab-pane>
              <el-tab-pane name="assert">
                <template #label>
                  <span>
                    断言
                    <el-badge
                      v-show="assertNotBlankLines(case_info.assert_form) > 0"
                      class="badge"
                      type="primary"
                      :value="assertNotBlankLines(case_info.assert_form)"
                    />
                  </span>
                </template>
                <el-row
                  v-for="(item, index) in case_info.assert_form"
                  class="kv-row"
                  :key="index"
                  :gutter="6"
                >
                  <el-col :span="0.5" class="kv-checkbox">
                    <el-checkbox
                      v-show="!eleIsDisable(index, case_info.assert_form)"
                      v-model="item.enable"
                    />
                  </el-col>
                  <el-col :span="4">
                    <el-select
                      v-model="item.extract_type"
                      :disabled="!eleIsDisable(index, case_info.assert_form) && !item.enable"
                      class="method-select"
                    >
                      <el-option
                        v-for="it in extract_type_list"
                        :key="it.value"
                        :disabled="it.disabled"
                        :label="it.label"
                        :value="it.value"
                      />
                    </el-select>
                  </el-col>
                  <el-col :span="5">
                    <el-input
                      v-model="item.extract_express"
                      :disabled="!eleIsDisable(index, case_info.assert_form) && !item.enable"
                      clearable
                      placeholder="提取表达式"
                      @change="assertFormChange(case_info.assert_form)"
                    />
                  </el-col>
                  <el-col :span="4">
                    <el-select
                      v-model="item.checking_rule"
                      :disabled="!eleIsDisable(index, case_info.assert_form) && !item.enable"
                      class="method-select"
                    >
                      <el-option
                        v-for="it in checking_rule_list"
                        :key="it.value"
                        :disabled="it.disabled"
                        :label="it.label"
                        :value="it.value"
                      />
                    </el-select>
                  </el-col>
                  <el-col :span="5">
                    <el-input
                      v-model="item.expect_value"
                      :disabled="!eleIsDisable(index, case_info.assert_form) && !item.enable"
                      clearable
                      placeholder="预期值"
                      @change="assertFormChange(case_info.assert_form)"
                    />
                  </el-col>
                  <el-col :span="5">
                    <el-input
                      v-model="item.desc"
                      :disabled="!eleIsDisable(index, case_info.assert_form) && !item.enable"
                      clearable
                      placeholder="备注"
                    />
                  </el-col>
                  <el-col :span="0.5" class="kv-delete">
                    <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
                      <Icon
                        v-show="!eleIsDisable(index, case_info.assert_form)"
                        icon="ep:delete"
                        class="icon-button"
                        @click.stop="case_info.assert_form.splice(index, 1)"
                      />
                    </el-tooltip>
                  </el-col>
                </el-row>
              </el-tab-pane>
              <el-tab-pane label="变量提取" name="variable">
                <template #label>
                  <span>
                    变量提取
                    <el-badge
                      v-show="notBlankLines(case_info.variable_form) > 0"
                      class="badge"
                      type="primary"
                      :value="notBlankLines(case_info.variable_form)"
                    />
                  </span>
                </template>
                <el-row
                  v-for="(item, index) in case_info.variable_form"
                  class="kv-row"
                  :key="index"
                  :gutter="6"
                >
                  <el-col :span="0.5" class="kv-checkbox">
                    <el-checkbox
                      v-show="!eleIsDisable(index, case_info.variable_form)"
                      v-model="item.enable"
                    />
                  </el-col>
                  <el-col :span="4">
                    <el-select
                      v-model="item.extract_type"
                      :disabled="!eleIsDisable(index, case_info.variable_form) && !item.enable"
                    >
                      <el-option
                        v-for="it in extract_type_list"
                        :key="it.value"
                        :label="it.label"
                        :value="it.value"
                      />
                    </el-select>
                  </el-col>
                  <el-col :span="10">
                    <el-input
                      v-model="item.key"
                      :disabled="!eleIsDisable(index, case_info.variable_form) && !item.enable"
                      clearable
                      placeholder="提取方式"
                      @change="formChange(case_info.variable_form)"
                    />
                  </el-col>
                  <el-col :span="9">
                    <el-input
                      v-model="item.value"
                      :disabled="!eleIsDisable(index, case_info.variable_form) && !item.enable"
                      clearable
                      placeholder="变量名"
                      @change="formChange(case_info.variable_form)"
                    />
                  </el-col>
                  <el-col :span="0.5" class="kv-delete">
                    <el-tooltip
                      v-if="!eleIsDisable(index, case_info.variable_form)"
                      effect="dark"
                      :open-delay="200"
                      content="删除"
                      placement="top"
                    >
                      <Icon
                        v-show="!eleIsDisable(index, case_info.variable_form)"
                        icon="ep:delete"
                        class="icon-button"
                        @click.stop="case_info.variable_form.splice(index, 1)"
                      />
                    </el-tooltip>
                    <div v-else style="width: 16px"></div>
                  </el-col>
                </el-row>
              </el-tab-pane>
              <el-tab-pane v-if="case_info.rally_point" label="集合点">
                <el-form
                  ref="form"
                  :model="case_info.rally_point"
                  label-width="80px"
                  style="width: 300px"
                >
                  <el-form-item label="启用">
                    <el-switch
                      v-model="case_info.rally_point.enable"
                      style="--el-switch-on-color: #13ce66"
                    />
                  </el-form-item>
                  <el-form-item label="并发数量">
                    <el-input
                      v-model.number="case_info.rally_point.concurrency"
                      :disabled="!case_info.rally_point.enable"
                      size="small"
                    />
                  </el-form-item>
                  <el-form-item label="超时时间">
                    <el-input
                      v-model.number="case_info.rally_point.timeout_period"
                      :disabled="!case_info.rally_point.enable"
                      size="small"
                    />
                  </el-form-item>
                </el-form>
              </el-tab-pane>
              <el-tab-pane v-if="case_info.waiting_config" label="等待时间">
                <el-form :inline="true" :model="case_info.waiting_config" label-width="80px">
                  <el-form-item label="前置等待">
                    <el-switch
                      v-model="case_info.waiting_config.pre_waiting_switch"
                      style="--el-switch-on-color: #13ce66"
                    />
                  </el-form-item>
                  <el-form-item v-show="case_info.waiting_config.pre_waiting_switch" label="">
                    <el-input
                      v-model.number="case_info.waiting_config.pre_waiting_time"
                      style="width: 120px"
                      size="small"
                    >
                      <template #suffix>
                        <div>ms</div>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-form>
                <el-form :inline="true" :model="case_info.waiting_config" label-width="80px">
                  <el-form-item label="后置等待">
                    <el-switch
                      v-model="case_info.waiting_config.post_waiting_switch"
                      style="--el-switch-on-color: #13ce66"
                    />
                  </el-form-item>
                  <el-form-item v-show="case_info.waiting_config.post_waiting_switch" label="">
                    <el-input
                      v-model.number="case_info.waiting_config.post_waiting_time"
                      style="width: 120px"
                      size="small"
                    >
                      <template #suffix>
                        <div>ms</div>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-form>
                <el-form :inline="true" :model="case_info.overtime_config" label-width="80px">
                  <el-form-item label="超时时间">
                    <el-switch
                      v-model="case_info.overtime_config.enable"
                      style="--el-switch-on-color: #13ce66"
                    />
                  </el-form-item>
                  <el-form-item v-show="case_info.overtime_config.enable" label="">
                    <el-input
                      v-model.number="case_info.overtime_config.timeout_period"
                      style="width: 120px"
                      size="small"
                    >
                      <template #suffix>
                        <div>ms</div>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-form>
              </el-tab-pane>
            </el-tabs>
          </pane>
          <pane min-size="20" size="60" style="background-color: #fff">
            <div v-loading="request_loading" style="height: 100%">
              <el-empty v-if="Object.keys(debug_result).length === 0" description="暂无数据" />
              <http-resp
                v-else
                :result="debug_result"
                :read-only="false"
                @value-extract="valueExtract"
              />
            </div>
          </pane>
        </splitpanes>
      </div>
    </el-drawer>
  </div>
</template>

<style scoped lang="less">
.container {
  :deep(.el-drawer__header) {
    height: 20px;
    margin: 0;
  }
  :deep(.el-drawer__body) {
    padding: 0;
  }
  .name-edit-input {
    width: 400px;
    cursor: text;
  }
  .name-span {
    font-size: 16px;
    font-weight: 500;
    color: #024e9b;
  }

  .request-detail {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding-top: 5px;

    :deep(.splitpanes__splitter) {
      height: 0;
    }

    .full-url {
      height: 40px;
      display: flex;
      margin-left: 20px;
      align-items: center;

      .url-input {
        width: calc(100% - 200px);
      }

      .method-select {
        width: 100px;
      }

      .url-btns {
        margin-left: 10px;
      }
      .send-btn {
        margin-left: 10px;
      }
    }

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
        height: calc(100% - 42px);
        overflow: auto;
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
      :deep(.el-form-item) {
        margin-bottom: 0;
      }

      .badge {
        margin-top: 3px;
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

        .kv-select {
          width: 500px;
        }
      }

      .monaco-container {
        height: calc(100% - 40px);
        border: 1px solid #eee;
      }
    }

    .variable-checkbox {
      position: absolute;
      right: 72px;
      top: 20px;
    }

    .response-div {
      width: 100%;
      height: 100%;
    }
  }
}
</style>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import SceneService from '@/api/scene/scene'
import { ElMessage } from 'element-plus'
import FileService from '@/api/scene/file'

const sceneService = new SceneService()
const fileService = new FileService()

const props = defineProps({
  sceneData: {
    type: Object
  },
  planId: {
    type: Number
  }
})

interface caseInfo {
  label: string
  options: option[]
}

interface option {
  label: string
  value: string
}

const param = reactive({
  variable_list: [],
  export_times: 1,
  concurrency: 1,
  has_cache: false,
  disable_cache: false
})
let all_variable = ref<caseInfo[]>([])
let download_loading = ref(false)

const save = () => {
  const base_data = { ...props.sceneData }
  const _param = { ...param }
  if (!_param.export_times) {
    _param.export_times = 0
  }
  if (!_param.concurrency) {
    _param.concurrency = 0
  }
  if (_param.export_times <= _param.concurrency) {
    ElMessage.warning('导出次数必须大于并发数')
    return
  }
  base_data.export_data_info = _param
  base_data.plan_id = props.planId
  sceneService.updateScene(base_data)
}

const fileDownload = async () => {
  download_loading.value = true
  try {
    await fileService.downloadFile({
      path: 'export/scene',
      file_name: props.sceneData?.scene_id + '.csv',
      alias: props.sceneData?.scene_name + '.csv'
    })
  } finally {
    download_loading.value = false
  }
}

watch(
  () => props.sceneData,
  (newVal) => {
    const _variable_list: caseInfo[] = []
    if (!newVal?.disabled) {
      newVal?.case_tree.forEach((cs) => {
        if (!cs.disabled) {
          const _case_info: caseInfo = {
            label: cs.title,
            options: []
          }
          cs.extend.variable_form.forEach((vrb) => {
            if (vrb.enable && vrb.value !== '') {
              _case_info.options.push(vrb.value)
            }
          })
          _variable_list.push(_case_info)
        }
      })
    }
    all_variable.value = _variable_list

    param.variable_list = props.sceneData?.export_data_info?.variable_list || []
    param.export_times = props.sceneData?.export_data_info?.export_times || null
    param.concurrency = props.sceneData?.export_data_info?.concurrency || null
    param.has_cache = props.sceneData?.export_data_info?.has_cache
    param.disable_cache = props.sceneData?.export_data_info?.disable_cache
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <el-card class="case-card">
    <el-tag size="small">数据导出</el-tag>
    <el-switch
      v-if="param.has_cache"
      v-model="param.disable_cache"
      inline-prompt
      active-text="重新执行"
      inactive-text="使用缓存"
      active-color="#13ce66"
      style="margin-left: 28px; --el-switch-off-color: #13ce66"
      @change="save"
    />
    <el-select
      v-model="param.variable_list"
      multiple
      collapse-tags
      collapse-tags-tooltip
      :max-collapse-tags="1"
      v-show="param.disable_cache || !param.has_cache"
      :disabled="sceneData?.disabled"
      placeholder="选择变量"
      style="width: 240px; margin-left: 28px"
      @change="save"
    >
      <el-option-group v-for="group in all_variable" :key="group.label" :label="group.label">
        <el-option v-for="item in group.options" :key="item" :label="item" :value="item" />
      </el-option-group>
    </el-select>
    <el-input
      placeholder="导出量级"
      v-show="param.disable_cache || !param.has_cache"
      style="width: 240px; margin-left: 20px"
      :disabled="sceneData?.disabled"
      v-model.number="param.export_times"
      @change="save"
    />
    <el-input
      v-show="param.disable_cache || !param.has_cache"
      placeholder="并发数"
      style="width: 240px; margin-left: 20px"
      :disabled="sceneData?.disabled"
      v-model.number="param.concurrency"
      @change="save"
    />
    <el-button
      :loading="download_loading"
      text
      type="primary"
      style="margin-left: 20px"
      v-show="!param.disable_cache"
      @click="fileDownload()"
    >
      下载缓存
    </el-button>
  </el-card>
</template>

<style scoped lang="less">
.case-card {
  flex: 1 1 auto;
  display: flex;
  align-items: center;
  font-size: 14px;
  height: 46px;
  margin: 2px 8px 0 6px;
}
</style>

<script setup lang="ts">
import { reactive, ref, unref } from 'vue'
import type { SceneInfo } from '@/views/Plan/scene/case-data'
import { ElMessage } from 'element-plus'
import SceneService from '@/api/scene/scene'
import { OperationType } from '@/constants'

const sceneService = new SceneService()

const dialog_visible = ref(false)
const edit_btn_loading = ref(false)
let edit_type = ref(-1)
let scene_form = reactive<SceneInfo>({})

const emit = defineEmits(['refresh'])
function refreshPlanScene() {
  emit('refresh')
}

const preEditScene = (edit_tp, val) => {
  dialog_visible.value = true
  edit_type.value = edit_tp
  scene_form = Object.assign(scene_form, val)
}

const editScene = async () => {
  if (scene_form.scene_name === '') {
    ElMessage.error('场景名不能为空')
  }
  edit_btn_loading.value = true
  try {
    unref(edit_type) === OperationType.CREATE
      ? await sceneService.createScene(scene_form)
      : await sceneService.updateScene(scene_form)
    refreshPlanScene()
    dialog_visible.value = false
  } finally {
    edit_btn_loading.value = false
  }
}

const scene_type_list = [
  {
    label: '普通场景',
    value: 0
  },
  {
    label: '前置场景',
    value: 1
  }
]

defineExpose({
  preEditScene
})
</script>

<template>
  <el-dialog
    v-model="dialog_visible"
    :title="edit_type === 0 ? '创建场景' : '编辑场景'"
    @closed="scene_form = reactive({})"
    width="30%"
  >
    <el-form :model="scene_form" label-position="right">
      <el-form-item label="场景名" label-width="80px">
        <el-input
          v-model="scene_form.scene_name"
          maxlength="40"
          show-word-limit
          style="width: 90%"
        />
      </el-form-item>
      <el-form-item label="场景类型" label-width="80px">
        <el-select v-model="scene_form.scene_type" placeholder="请选择场景类型" style="width: 90%">
          <el-option
            v-for="item in scene_type_list"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialog_visible = false">取 消</el-button>
        <el-button
          :loading="edit_btn_loading"
          type="primary"
          v-tc="{
            name: edit_type === 0 ? '创建场景' : '编辑场景',
            scene_id: scene_form.scene_id
          }"
          @click="editScene"
        >
          确 定
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="less"></style>

<script setup lang="ts">
import { CaseType } from '@/constants'
import CommonCase from '@/views/Plan/scene/CommonCase.vue'
import type { SceneInfo } from '@/views/Plan/scene/case-data'
import CaseService from '@/api/scene/case'

const caseService = new CaseService()

const emit = defineEmits(['operateCase'])
const operateCase = (opType, data) => {
  emit('operateCase', opType, data)
}
// const scene_data = defineModel<SceneInfo>()
const props = defineProps({
  sceneData: {
    type: Object as SceneInfo
  }
})
const expand_node = []

const allowDrop = (draggingNode, dropNode, type) => {
  if (type === 'inner') {
    if (
      draggingNode.data.type !== CaseType.HttpCase ||
      dropNode.data.type !== CaseType.LogicControl
    ) {
      return false
    }
  }
  return true
}
// 节点展开
const nodeExpand = () => {}
// 节点关闭
const nodeCollapse = () => {}
const nodeDrop = (before, after, inner) => {
  const inner_val = { before: 0, after: 1, inner: 2 }
  const param = {
    before: {
      scene_id: before.data.scene_id,
      case_id: before.data.case_id,
      sort: before.data.sort,
      parent_id: before.data.parent_id
    },
    after: {
      scene_id: after.data.scene_id,
      case_id: after.data.case_id,
      sort: after.data.sort,
      parent_id: after.data.parent_id
    },
    position: inner_val[inner]
  }
  caseService.resortSceneCase(param)
}
</script>

<template>
  <el-tree
    ref="treeRef"
    class="ms-tree"
    :draggable="true"
    highlight-current
    node-key="case_id"
    :indent="15"
    :data="props.sceneData.case_tree"
    :default-expanded-keys="expand_node"
    :allow-drop="allowDrop"
    :expand-on-click-node="false"
    @node-expand="nodeExpand"
    @node-collapse="nodeCollapse"
    @node-drop="nodeDrop"
  >
    <template #default="{ data }">
      <common-case
        :case-data="data"
        :scene-disabled="props.sceneData.disabled"
        @operate-case="operateCase"
      />
    </template>
  </el-tree>
</template>

<style scoped lang="less">
.ms-tree {
  height: 100%;
  margin-top: 2px;
  overflow-y: auto;
  padding-right: 8px;
  margin-left: -6px;

  :deep(.el-tree-node__content) {
    height: 52px;
  }

  :deep(.el-tree-node__expand-icon) {
    padding: 0;
  }
}
</style>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import HttpCase from '@/views/Plan/scene/HttpCase.vue'
import LoginController from '@/views/Plan/scene/LogicController.vue'

import { CaseType } from '@/constants'

const props = defineProps({
  caseData: {
    type: Object
  },
  sceneDisabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['operateCase'])
const operateCase = (opType, data) => {
  emit('operateCase', opType, data)
}

let current_node_data = reactive({})

let component_type = computed(() => {
  switch (current_node_data.type) {
    case CaseType.HttpCase:
      return HttpCase
    case CaseType.LogicControl:
      return LoginController
    default:
      return HttpCase
  }
})
// onMounted(() => {
// })
watch(
  () => props.caseData,
  (newVal) => {
    current_node_data = Object.assign(current_node_data, newVal)
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <el-card
    :class="[props.sceneDisabled || current_node_data.disabled ? 'not-valid' : '', 'case-card']"
  >
    <component :is="component_type" v-model="current_node_data" @operate-case="operateCase" />
  </el-card>
</template>

<style scoped lang="less">
.case-card {
  flex: 1 1 auto;
  display: flex;
  align-items: center;
  font-size: 14px;
  width: 100%;
  height: 46px;

  :deep(.el-card__body) {
    width: 100%;
  }
}

.not-valid {
  background-color: #f5f5f5;
  border: #e9e9eb solid 1px;
}
</style>

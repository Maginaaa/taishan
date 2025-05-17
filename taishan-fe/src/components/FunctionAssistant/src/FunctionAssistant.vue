<script setup lang="ts">
import { Icon } from '@/components/Icon'
import { ref } from 'vue'
import FunctionDialog from '@/components/FunctionAssistant/src/FunctionDialog.vue'

const dialogRef = ref()
</script>

<template>
  <div class="icon" v-tc="{ name: '函数助手' }" @click="dialogRef.openDialog()">
    <Icon :size="18" icon="ant-design:experiment-outlined" color="var(--top-header-text-color)" />
  </div>
  <function-dialog ref="dialogRef" />
</template>

<style scoped lang="less">
.icon {
  display: flex;
  height: 100%;
  padding: 1px 10px 0;
  cursor: pointer;
  align-items: center;
  transition: background var(--transition-time-02);
}
</style>

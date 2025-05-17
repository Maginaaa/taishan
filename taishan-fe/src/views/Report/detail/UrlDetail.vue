<script setup>
import { computed } from 'vue'
import { extractPath } from '@/utils/url'

const props = defineProps({
  apiData: {
    type: Array,
    required: true,
    default: () => []
  },
  height: {
    type: String,
    default: '280px' // 默认高度
  }
})

const processedData = computed(() => {
  // 找到最大总数
  const maxTotal = Math.max(...props.apiData.map((item) => item.count))

  return props.apiData.map((item) => {
    const successCount = item.count - item.error_count
    const errorCount = item.error_count

    // 计算宽度百分比（基于最大总数）
    const successWidth = (successCount / maxTotal) * 100
    const errorWidth = (errorCount / maxTotal) * 100

    return {
      path: extractPath(item.url),
      successCount,
      errorCount,
      successWidth,
      errorWidth
    }
  })
})
</script>

<template>
  <div class="data-container" :style="{ maxHeight: props.height, marginTop: '5px' }">
    <div v-for="(item, index) in processedData" :key="index" class="data-row">
      <!-- URL 文本 -->
      <div class="url-text">{{ item.path }}</div>
      <!-- 柱状图容器 -->
      <div class="bar-container">
        <!-- 成功部分 -->
        <div class="success-bar" :style="{ width: `${item.successWidth}%` }"></div>
        <!-- 失败部分 -->
        <div class="error-bar" :style="{ width: `${item.errorWidth}%` }"></div>
        <div class="success-number">{{ item.successCount }}</div> /
        <div class="error-number">{{ item.errorCount }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.data-container {
  overflow-y: auto;

  @success-color: #91cc75;
  @error-color: #ee6666;
  @warning-color: #fac858;

  .data-row {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    margin-right: 20px;

    /* URL 文本样式 */
    .url-text {
      width: 20%;
      font-size: 14px;
      color: #333;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    /* 柱状图容器 */
    .bar-container {
      display: flex;
      align-items: center;
      width: 80%;
      height: 20px;
      position: relative;
      overflow: hidden;
    }

    .success-bar {
      background-color: @success-color;
      height: 80%;
      border-radius: 4px;
    }

    .error-bar {
      background-color: @error-color;
      height: 80%;
      border-radius: 4px;
    }

    .success-number {
      margin-left: 10px;
      margin-right: 4px;
      color: @success-color;
    }

    .error-number {
      margin-left: 4px;
      color: @error-color;
    }
  }
}
/* 每行数据 */
</style>

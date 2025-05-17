<script setup lang="ts">
import { ref, watch, computed, onMounted, unref } from 'vue'
import { useAppStore } from '@/store/modules/app'

const appStore = useAppStore()

const show = ref(true)

const title = computed(() => '泰山')

const layout = computed(() => appStore.getLayout)

const collapse = computed(() => appStore.getCollapse)

onMounted(() => {
  if (unref(collapse)) show.value = false
})

watch(
  () => collapse.value,
  (collapse: boolean) => {
    if (unref(layout) === 'topLeft' || unref(layout) === 'cutMenu') {
      show.value = true
      return
    }
    show.value = !collapse
  }
)

watch(
  () => layout.value,
  (layout) => {
    if (layout === 'top' || layout === 'cutMenu') {
      show.value = true
    } else {
      if (unref(collapse)) {
        show.value = false
      } else {
        show.value = true
      }
    }
  }
)
</script>

<template>
  <div class="logo-container">
    <router-link v-tc="{ name: 'icon' }" to="/" class="link">
      <img
        src="@/assets/imgs/logo.png"
        class="w-[calc(var(--logo-height)-10px)] h-[calc(var(--logo-height)-10px)]"
      />
      <div v-if="show" class="title">
        {{ title }}
      </div>
    </router-link>
  </div>
</template>

<style scoped lang="less">
.logo-container {
  display: flex;
  // 水平布局
  justify-content: center;

  .link {
    text-decoration: none;
    display: flex;
    align-items: center;

    .title {
      margin-left: 12px;
      font-size: 22px;
      font-weight: 700;
      color: #1c64a8;
    }
  }
}
</style>

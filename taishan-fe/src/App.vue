<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { ConfigGlobal } from '@/components/ConfigGlobal'
import { useDesign } from '@/hooks/web/useDesign'
import { useAuthStore } from '@/store/modules/auth'
import { getToken } from '@/utils/auth'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('app')

const appStore = useAppStore()

const currentSize = computed(() => appStore.getCurrentSize)

const greyMode = computed(() => appStore.getGreyMode)

const authStore = useAuthStore()

appStore.initTheme()

const initializeToken = () => {
  const token = getToken()
  if (token) {
    authStore.setToken(token)
  } else {
    window.location.href = `https://oauth.xxx.com/login?url=${window.location.href}`
  }
}

initializeToken()

// ElNotification({
//   title: '提示',
//   type: 'warning',
//   duration: 0,
//   dangerouslyUseHTMLString: true,
//   message:
//     '<div><p><strong>遇事不决，请先查阅常见问题，说不定你能找到相关解答</strong></p><p><a href="https://element-plus-admin-doc.cn/guide/fqa.html" target="_blank">链接地址</a></p></div>'
// })
</script>

<template>
  <ConfigGlobal :size="currentSize">
    <RouterView :class="greyMode ? `${prefixCls}-grey-mode` : ''" />
  </ConfigGlobal>
</template>

<style lang="less">
@prefix-cls: ~'@{namespace}-app';

.size {
  width: 100%;
  height: 100%;
}

html,
body {
  padding: 0 !important;
  margin: 0;
  overflow: hidden;
  .size;

  #app {
    .size;
  }
}

.@{prefix-cls}-grey-mode {
  filter: grayscale(100%);
}
</style>

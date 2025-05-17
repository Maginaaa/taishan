<script setup lang="ts">
import { ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'
import { useDesign } from '@/hooks/web/useDesign'
import { useAuthStore } from '@/store/modules/auth'
import AccountService from '@/api/account'
// const userStore = useUserStore()

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('user-info')

const authStore = useAuthStore()
const accountService = new AccountService()

const initializeUserInfo = async () => {
  const res = await accountService.getUserInfo()
  const userInfo = { ...res }
  authStore.setUserInfo({
    id: userInfo.id,
    username: userInfo.name,
    avatar: userInfo.avatarUrl
  })
}

initializeUserInfo()

const loginOut = () => {
  // userStore.logoutConfirm()
  authStore.logout()
}
</script>

<template>
  <ElDropdown class="custom-hover" :class="prefixCls" trigger="click">
    <div v-tc="{ name: '个人信息' }" class="flex items-center">
      <img
        :src="authStore.getUserInfo?.avatar"
        alt=""
        class="w-[calc(var(--logo-height)-25px)] rounded-[49%]"
      />
      <span class="<lg:hidden text-14px pl-[5px] text-[var(--top-header-text-color)]">
        {{ authStore.getUserInfo?.username }}
      </span>
    </div>
    <template #dropdown>
      <ElDropdownMenu>
        <ElDropdownItem>
          <div v-tc="{ name: '退出登录' }" @click="loginOut">退出登录</div>
        </ElDropdownItem>
      </ElDropdownMenu>
    </template>
  </ElDropdown>
</template>

<style scoped lang="less">
.fade-bottom-enter-active,
.fade-bottom-leave-active {
  transition:
    opacity 0.25s,
    transform 0.3s;
}

.fade-bottom-enter-from {
  opacity: 0;
  transform: translateY(-10%);
}

.fade-bottom-leave-to {
  opacity: 0;
  transform: translateY(10%);
}
</style>

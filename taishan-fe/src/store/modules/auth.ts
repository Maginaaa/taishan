import { defineStore } from 'pinia'
import { UserType } from '@/api/login/types'
import { store } from '@/store'
import { useTagsViewStore } from '@/store/modules/tagsView'
import { ElMessageBox } from 'element-plus'
import { getToken } from '@/utils/auth'
import AccountService from '@/api/account'

// import { loginApi } from '@/api/account'

interface AuthState {
  token: string
  userInfo?: UserType
}

const accountService = new AccountService()

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => {
    return {
      token: '',
      userInfo: undefined
    }
  },
  getters: {
    getToken(): string {
      return 'token'
      // return this.token
    },
    getUserInfo(): UserType | undefined | any {
      return {
        username: '简单随风',
        avatar: 'https://avatars.githubusercontent.com/u/19279437'
      }
      // return this.userInfo
    }
  },
  actions: {
    jumpToSso() {
      window.location.href = `https://oauth.xxxx.com/login?url=${window.location.href}`
    },
    initToken() {
      const token = getToken
      this.setToken(token)
    },
    logout() {
      ElMessageBox.confirm('确定退出吗', '温馨提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          const res = await accountService.userLogout()
          if (res) {
            this.reset()
            this.jumpToSso()
          }
        })
        .catch(() => {})
      // router.replace(`/login?redirect=${window.location.href}`)
    },
    setToken(token) {
      this.token = token
    },
    setUserInfo(userInfo?: UserType) {
      this.userInfo = userInfo
    },
    reset() {
      const tagsViewStore = useTagsViewStore()
      tagsViewStore.delAllViews()
      this.setToken('')
      this.setUserInfo(undefined)
    }
  },
  persist: true
})

export const useAuthStoreWithOut = () => {
  return useAuthStore(store)
}

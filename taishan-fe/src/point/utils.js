import sa from './sensors.js'
import { useAuthStore } from '@/store/modules/auth'

function buildParam(param = {}) {
  const authStore = useAuthStore()
  const baseParam = {
    $event_session_id: authStore.getUserInfo?.id,
    $element_name: param.name,
    $element_content: param.id || param.search || ''
  }
  const _param = param || {}
  Object.assign(baseParam, _param)
  return baseParam
}

const trackClick = (param = {}) => {
  sa.track('$buttonclick', buildParam(param))
}

const pageView = (param = {}) => {
  sa.track('$pageview', buildParam())
}

export { trackClick, pageView }

import request from '@/axios'

const domain = 'scene'
class UtilsService {
  async Debug(key) {
    const data = { key }
    const res = await request.post({ url: `/${domain}/utils/function_assistant/debug`, data })
    return res.data
  }
}

export default UtilsService

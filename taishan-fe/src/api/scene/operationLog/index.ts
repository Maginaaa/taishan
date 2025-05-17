import request from '@/axios'

const domain = 'scene'
class OperationLogService {
  async getOperationLog(data) {
    const res = await request.post({ url: `/${domain}/log`, data })
    return res.data
  }
}

export default OperationLogService

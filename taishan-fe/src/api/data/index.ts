import request from '@/axios'

const domain = 'data'
class DataService {
  async getDashboardData() {
    const res = await request.get({ url: `/${domain}/dashboard` })
    return res.data
  }

  async getTracesList(data) {
    const res = await request.post({ url: `/${domain}/arms/trace/list`, data })
    return res.data
  }

  async getTraceDetail(data) {
    const res = await request.post({ url: `/${domain}/arms/trace/detail`, data })
    return res.data
  }
}

export default DataService

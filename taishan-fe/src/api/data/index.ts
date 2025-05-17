import request from '@/axios'

const domain = 'data'
class DataService {
  async getDashboardData() {
    const res = await request.get({ url: `/${domain}/dashboard` })
    return res.data
  }

  async getSlsData(data) {
    const res = await request.post({ url: `/${domain}/sls/customer`, data })
    return res.data
  }

  async saveSnapshot(data) {
    const res = await request.post({ url: `/${domain}/snapshot/save`, data })
    return res.data
  }

  async getSnapshotList(data) {
    const res = await request.post({ url: `/${domain}/snapshot/list`, data })
    return res.data
  }

  async getSnapshotDetail(id) {
    const res = await request.get({ url: `/${domain}/snapshot/detail/${id}` })
    return res.data
  }

  async updateSnapshot(data) {
    const res = await request.post({ url: `/${domain}/snapshot/update`, data })
    return res.data
  }

  async getServiceDashboard(data) {
    const res = await request.post({ url: `/${domain}/sw/dashboard`, data })
    return res.data
  }

  async getServiceFlowData(data) {
    const res = await request.post({ url: `/${domain}/flow`, data })
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

  async getLogSize(data) {
    const res = await request.post({ url: `/${domain}/sls/log/size`, data })
    return res.data
  }
}

export default DataService

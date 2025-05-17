import request from '@/axios'

const domain = 'scene'
class NormalizedService {
  async getNormalizedList(data) {
    const res = await request.post({ url: `/${domain}/normalized/list`, data })
    return res.data
  }

  async createApiBinding(data) {
    const res = await request.post({ url: `/${domain}/normalized/create`, data })
    return res.data
  }

  async updateApiBinding(data) {
    const res = await request.post({ url: `/${domain}/normalized/update`, data })
    return res.data
  }

  async deleteApiBinding(id) {
    const res = await request.get({ url: `/${domain}/normalized/delete/${id}` })
    return res.data
  }

  async getNormalizedSwitch() {
    const res = await request.get({ url: `/${domain}/normalized/switch` })
    return res.data
  }

  async updateNormalizedSwitch(data) {
    const res = await request.post({ url: `/${domain}/normalized/switch/update`, data })
    return res.data
  }

  async getResultList(data) {
    const res = await request.post({ url: `/${domain}/normalized/result/list`, data })
    return res.data
  }

  async getResultHistogramData(data) {
    const res = await request.post({ url: `/${domain}/normalized/result/histogram`, data })
    return res.data
  }

  async getResultAnalysis(data) {
    const res = await request.post({ url: `/${domain}/normalized/result/analysis`, data })
    return res.data
  }

  async disposeTask(data) {
    const res = await request.post({ url: `/${domain}/normalized/result/task/dispose`, data })
    return res.data
  }

  async getTaskList(data) {
    const res = await request.post({ url: `/${domain}/normalized/result/task/list`, data })
    return res.data
  }
}

export default NormalizedService

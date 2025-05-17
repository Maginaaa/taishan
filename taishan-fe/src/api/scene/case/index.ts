import request from '@/axios'

const domain = 'scene'
class CaseService {
  async getSceneList(id) {
    const res = await request.get({ url: `/${domain}/scene/list/${id}` })
    return res.data
  }

  async createSceneCase(data) {
    const res = await request.post({ url: `/${domain}/case/create`, data })
    return res.data
  }

  async importCase(data) {
    const res = await request.post({ url: `/${domain}/case/import`, data })
    return res.data
  }

  async updateSceneCase(data) {
    const res = await request.post({ url: `/${domain}/case/update`, data })
    return res?.data
  }

  async deleteSceneCase(id) {
    const res = await request.get({ url: `/${domain}/case/delete/${id}` })
    return res.data
  }
  async getSceneCaseTree(id) {
    const res = await request.get({ url: `/${domain}/case/tree/${id}` })
    return res.data
  }

  async resortSceneCase(data) {
    const res = await request.post({ url: `/${domain}/case/sort`, data })
    return res.data
  }

  async caseDebug(data) {
    const res = await request.post({ url: `/${domain}/case/debug`, data })
    return res.data
  }
}

export default CaseService

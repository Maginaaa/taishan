import request from '@/axios'

const domain = 'scene'
class VariableService {
  async createSceneVariable(data) {
    const res = await request.post({ url: `/${domain}/scene/variable/create`, data })
    return res.data
  }

  async updateSceneVariable(data) {
    const res = await request.post({ url: `/${domain}/scene/variable/update`, data })
    return res.data
  }

  async deleteSceneVariable(id) {
    const res = await request.get({ url: `/${domain}/scene/variable/delete/${id}` })
    return res.data
  }
}

export default VariableService

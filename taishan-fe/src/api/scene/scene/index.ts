import request from '@/axios'

const domain = 'scene'
class SceneService {
  async getSceneList(id) {
    const res = await request.get({ url: `/${domain}/scene/list/${id}` })
    return res.data
  }

  async createScene(data) {
    const res = await request.post({ url: `/${domain}/scene/create`, data })
    return res.data
  }

  async updateScene(data) {
    const res = await request.post({ url: `/${domain}/scene/update`, data })
    return res.data
  }

  async copyScene(id) {
    const res = await request.get({ url: `/${domain}/scene/copy/${id}` })
    return res.data
  }

  async deleteScene(id) {
    const res = await request.get({ url: `/${domain}/scene/delete/${id}` })
    return res.data
  }

  async updateSceneSort(data) {
    const res = await request.post({ url: `/${domain}/scene/sort`, data })
    return res.data
  }
}

export default SceneService

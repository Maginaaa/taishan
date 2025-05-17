import request from '@/axios'

const domain = 'env'
class EnvService {
  async getConfigSyncHistory(data) {
    const res = await request.post({ url: `/${domain}/conf/sync/list`, data })
    return res.data
  }

  async getConfigSyncDetail(id) {
    const res = await request.post({ url: `/${domain}/conf/sync/${id}` })
    return res.data
  }

  async syncVersion(data) {
    const res = await request.post({ url: `/${domain}/version/sync`, data })
    return res.data
  }

  async syncTdConfig() {
    const res = await request.get({ url: `/${domain}/conf/sync` })
    return res.data
  }

  async getAllAppID() {
    const res = await request.get({ url: `/${domain}/conf/appid/list` })
    return res.data
  }

  async getNamespace(data) {
    const res = await request.post({ url: `/${domain}/conf/server/namespace`, data })
    return res.data
  }

  async getConfigDetail(data) {
    const res = await request.post({ url: `/${domain}/conf/server/config`, data })
    return res.data
  }

  async saveConfigDetail(data) {
    const res = await request.post({ url: `/${domain}/conf/server/config/save`, data })
    return res.data
  }

  async updateConfigDetail(data) {
    const res = await request.post({ url: `/${domain}/conf/server/config/update`, data })
    return res.data
  }

  async getConfigList(data) {
    const res = await request.post({ url: `/${domain}/conf/server/config/list`, data })
    return res.data
  }

  async configDoCheck() {
    const res = await request.get({ url: `/${domain}/conf/check` })
    return res.data
  }

  async getConfigCheckHistory(data) {
    const res = await request.post({ url: `/${domain}/conf/check/history`, data })
    return res.data
  }

  async getConfigCheckDetail(id) {
    const res = await request.get({ url: `/${domain}/conf/check/${id}` })
    return res.data
  }
}

export default EnvService

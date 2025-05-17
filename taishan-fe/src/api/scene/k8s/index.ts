import request from '@/axios'

const domain = 'scene'
class K8sService {
  async getNamespaces() {
    const res = await request.get({ url: `/${domain}/k8s/namespace/list` })
    return res.data
  }

  async getDeployments(data) {
    const res = await request.post({ url: `/${domain}/k8s/deployment/list`, data })
    return res.data
  }

  async getPlanServerInfo(id) {
    const res = await request.get({ url: `/${domain}/k8s/server/list/${id}` })
    return res.data
  }

  async updateServerReplicas(data) {
    const res = await request.post({ url: `/${domain}/k8s/server/update`, data })
    return res.data
  }
}

export default K8sService

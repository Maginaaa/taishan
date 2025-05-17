import request from '@/axios'

const domain = 'scene'
class PlanService {
  async createPlan(data) {
    const res = await request.post({ url: `/${domain}/plan/create`, data })
    return res.data
  }

  async copyPlan(data) {
    const res = await request.post({ url: `/${domain}/plan/copy`, data })
    return res.data
  }

  async getPlanList(data) {
    const res = await request.post({ url: `/${domain}/plan/list`, data })
    return res.data
  }

  async getAllPlanList() {
    const res = await request.get({ url: `/${domain}/plan/list/all` })
    return res.data
  }

  async updatePlan(data) {
    const res = await request.post({ url: `/${domain}/plan/update`, data })
    return res
  }

  async deletePlan(id) {
    const res = await request.get({ url: `/${domain}/plan/delete/${id}` })
    return res.data
  }

  async updatePlanSimple(data) {
    const res = await request.post({ url: `/${domain}/plan/update/simple`, data })
    return res.data
  }

  async debugPlan(id) {
    const res = await request.get({ url: `/${domain}/plan/debug/${id}` })
    return res.data
  }

  async executePlan(id) {
    const res = await request.get({ url: `/${domain}/plan/execute/${id}` })
    return res.data
  }

  async getPlanDetail(id) {
    const res = await request.get({ url: `/${domain}/plan/detail/${id}` })
    return res.data
  }

  async getPlanDebugRecord(data) {
    const res = await request.post({ url: `/${domain}/plan/record`, data })
    return res.data
  }
}

export default PlanService

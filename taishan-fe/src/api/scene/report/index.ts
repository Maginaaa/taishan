import request from '@/axios'

const domain = 'report'
class ReportService {
  async getReportDetail(id) {
    const res = await request.get({ url: `/${domain}/report/detail/${id}` })
    return res.data
  }

  async getReportData(id) {
    const res = await request.get({ url: `/${domain}/report/data/${id}` })
    return res.data
  }

  async getReportCaseData(data) {
    const res = await request.post({ url: `/${domain}/report/data/case`, data })
    return res.data
  }

  async getReportList(data) {
    const res = await request.post({ url: `/${domain}/report/list`, data })
    return res.data
  }

  async updatePress(data) {
    const res = await request.post({ url: `/${domain}/report/update/currency`, data })
    return res.data
  }

  async releasePreScene(data) {
    const res = await request.post({ url: `/${domain}/report/update/release`, data })
    return res.data
  }

  async stopPress(id) {
    const res = await request.get({ url: `/${domain}/report/stop/${id}` })
    return res.data
  }

  async updateReportName(data) {
    const res = await request.post({ url: `/${domain}/report/update/name`, data })
    return res.data
  }

  async getReportTargetRps(id) {
    const res = await request.get({ url: `/${domain}/report/rps/target/${id}` })
    return res.data
  }
}

export default ReportService

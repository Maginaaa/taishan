import request from '@/axios'

const domain = 'scene'
class SceneReportService {
  async getSamplingData(data) {
    const res = await request.post({ url: `/${domain}/report/sampling`, data })
    return res.data
  }

  reportRpsModify(data) {
    return request.post({ url: `/${domain}/report/rps/modify`, data })
  }
}

export default SceneReportService

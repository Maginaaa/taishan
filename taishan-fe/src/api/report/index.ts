import request from '@/axios'

const domain = 'scene'
class SceneReportService {
  async getSamplingData(data) {
    const res = await request.post({ url: `/${domain}/report/sampling`, data })
    return res.data
  }

  getReportRpsList(data) {
    return request.post({ url: `/${domain}/report/rpsList`, data })
  }

  // editRunningStepRps(data) {
  //   return request.post({ url: `/${domain}/report/rpsModify`, data })
  // }

  reportRpsModify(data) {
    return request.post({ url: `/${domain}/report/reportRpsModify`, data })
  }
}

export default SceneReportService

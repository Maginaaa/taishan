import request from '@/axios'

const domain = 'machine'
class MachineService {
  async getMachineList() {
    const res = await request.get({ url: `/${domain}/machine/available/list` })
    return res.data
  }

  async getMachineUseInfo(data) {
    const res = await request.post({ url: `/${domain}/machine/info/list`, data })
    return res.data
  }

  async getAllMachineUseInfo() {
    const res = await request.post({ url: `/${domain}/machine/info/all` })
    return res.data
  }
}

export default MachineService

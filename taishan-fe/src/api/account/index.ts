import request from '@/axios'
import type { User } from '@/views/Plan/user-data'

const domain = 'account'
class AccountService {
  async getUserList() {
    const res: IResponse<Array<User>> = await request.get({ url: `/${domain}/user/list` })
    return res.data
  }

  async getUserInfo() {
    const res = await request.get({ url: `/${domain}/user/info` })
    return res.data
  }

  async userLogout() {
    const res = await request.post({ url: `/${domain}/user/logout` })
    return res.data
  }
}

export default AccountService

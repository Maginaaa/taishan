import request from '@/axios'

const domain = 'account'
class AccountService {
  async getUserList() {
    // const res: IResponse<Array<User>> = await request.get({ url: `/${domain}/user/list` })
    return [
      {
        id: 1,
        name: '简单随风',
        nickname: '简单随风'
      }
    ]
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

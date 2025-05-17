import request from '@/axios'

const domain = 'scene'
class TagService {
  async GetTagList(type) {
    const data = { type }
    const res = await request.post({ url: `/${domain}/tag/list`, data })
    return res.data
  }
}

export default TagService

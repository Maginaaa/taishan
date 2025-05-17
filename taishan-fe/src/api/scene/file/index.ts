import request from '@/axios'
import { ElMessage } from 'element-plus'

const domain = 'scene'

class FileService {
  async fileUpload(dt) {
    const res = await request.post({
      data: dt,
      url: `/${domain}/file/upload`,
      headers: {
        'Content-type': 'multipart/form-data'
      }
    })
    return res.data
  }

  async getFileList(id) {
    const res = await request.get({ url: `/${domain}/file/list/${id}` })
    return res.data
  }

  async deleteFile(id) {
    const res = await request.get({ url: `/${domain}/file/delete/${id}` })
    return res.data
  }

  async fileColumnUpdate(data) {
    const res = await request.post({ url: `/${domain}/file/column/update`, data })
    return res.data
  }

  async downloadFile(dt) {
    const res = await request.post({
      data: dt,
      url: `/${domain}/file/download`,
      responseType: 'blob'
    })
    if (res.code === 0) {
      ElMessage.error('文件下载失败')
      return
    }
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', dt.alias || dt.file_name)
    document.body.appendChild(link)
    link.click()

    // 清理
    window.URL.revokeObjectURL(url)
    link.remove()
  }
}

export default FileService

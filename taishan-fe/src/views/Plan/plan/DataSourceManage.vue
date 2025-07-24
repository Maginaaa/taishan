<script setup lang="ts">
import { ref } from 'vue'
import FileService from '@/api/scene/file'
import { ElMessage } from 'element-plus'

const fileService = new FileService()

let props = defineProps({
  planId: {
    type: Number,
    default: 0
  }
})

let file_data = ref([])
let columns = ref<any[]>([])
const getFileList = async () => {
  const res = await fileService.getFileList(props.planId)
  let _columns: any[] = []
  res.forEach((file) => {
    _columns.push(...file.column)
  })
  file_data.value = res
  columns.value = _columns
  refreshAllAlias()
}

const read_type_list = [
  { label: '顺序读取', value: 0 },
  { label: '随机读取', value: 1 },
  { label: '单轮遍历', value: 2 }
]

const getReadTypeLabel = (val: number) => {
  const item = read_type_list.find((item) => item.value === val)
  return item ? item.label : ''
}

const upload_loading = ref(false)

const doUpload = async (options) => {
  upload_loading.value = true
  try {
    const { file } = options
    // 验证文件类型和大小
    const isCSV = file.type === 'text/csv'
    // 判断上传的文件大小是否小于 50MB
    if (!isCSV) {
      ElMessage.error('只能上传 CSV 文件')
      return false
    }
    const fileSize = file.size < 50 * 1024 * 1024
    if (!fileSize) {
      ElMessage.error('文件大小超过限制')
      return false
    }
    await fileService.fileUpload({
      file: file,
      plan_id: props.planId
    })
    await getFileList()
  } finally {
    upload_loading.value = false
  }
}

const fileDownload = async (file_name: string) => {
  await fileService.downloadFile({
    path: 'plan/' + props.planId,
    file_name: file_name
  })
}

const fileDelete = async (id: number) => {
  await fileService.deleteFile(id)
  await getFileList()
}

const startEditing = (row) => {
  if (row.alias === '') {
    row.alias = row.col
  }
  row.editing = true
}

const finishEditing = (row) => {
  const trimmedValue = row.alias.trim()
  if (!trimmedValue) {
    ElMessage.warning('参数名不能为空')
  }
  if (trimmedValue === row.col) {
    row.alias = ''
  }
  row.editing = false
  updateColumn(row)
}

const convertByte = (val) => {
  if (!val) {
    return 0
  }
  if (val >= 1024 && val < 1024 ** 2) {
    return (val / 1024).toFixed(2) + 'KB'
  } else if (val >= 1024 ** 2 && val < 1024 ** 3) {
    return (val / 1024 ** 2).toFixed(2) + 'MB'
  } else if (val >= 1024 ** 3) {
    return (val / 1024 ** 3).toFixed(2) + 'GB'
  } else {
    return val + 'Byte'
  }
}

const changeVariableReadType = (cmd, row) => {
  row.read_type = cmd
  updateColumn(row)
}

const updateColumn = (row) => {
  const _row = { ...row }
  _row.plan_id = props.planId
  fileService.fileColumnUpdate(_row)
  refreshAllAlias()
}

const repeatedVariable = ref<Array<string>>([])
const refreshAllAlias = () => {
  repeatedVariable.value = []
  let allVariable: Array<string> = []
  file_data.value.forEach((file: any) => {
    file.column.forEach((column) => {
      const _vrb = column.alias ? column.alias : column.col
      if (allVariable.includes(_vrb)) {
        repeatedVariable.value.push(_vrb)
      }
      allVariable.push(_vrb)
    })
  })
}

getFileList()
</script>

<template>
  <el-upload :http-request="doUpload" :show-file-list="false" :accept="'.csv'">
    <el-button
      type="primary"
      v-tc="{
        name: '上传文件',
        plan_id: props.planId
      }"
      :loading="upload_loading"
      >点击上传</el-button
    >
    <template #tip>
      <div class="el-upload__tip">仅支持CSV文件(限50M)</div>
    </template>
  </el-upload>
  <el-table
    :data="file_data"
    :border="true"
    style="width: 100%; margin-top: 20px"
    default-expand-all
    row-style="background-color: var(--el-fill-color-light);"
  >
    <el-table-column type="expand">
      <template #default="{ row }">
        <el-table :data="row.column">
          <el-table-column label="参数名">
            <template #default="{ row: childRow }">
              <div class="cell-wrapper">
                <el-input
                  v-if="childRow.editing"
                  v-model="childRow.alias"
                  size="small"
                  @blur="finishEditing(childRow)"
                />
                <el-tooltip v-else :content="`初始参数名: ${childRow.col}`" placement="top">
                  <span
                    v-if="childRow.alias"
                    :style="{ color: repeatedVariable.includes(childRow.alias) ? 'red' : '' }"
                    >{{ childRow.alias }}</span
                  >
                  <span
                    v-else
                    :style="{ color: repeatedVariable.includes(childRow.col) ? 'red' : '' }"
                    >{{ childRow.col }}</span
                  >
                </el-tooltip>
                <Icon
                  icon="ep:edit"
                  v-show="!childRow.editing"
                  class="cell-edit-icon"
                  v-tc="{
                    name: '修改参数名',
                    row: childRow
                  }"
                  @click="startEditing(childRow)"
                />
                <!-- 编辑图标，只有在悬停时才显示 -->
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="file_name" label="数据来源" align="center" />
          <el-table-column label="读取方式" align="center">
            <template #default="{ row: childRow }">
              <el-dropdown @command="(cmd) => changeVariableReadType(cmd, childRow)">
                <el-tag class="el-dropdown-link">
                  <div class="flex items-center">
                    <span>{{ getReadTypeLabel(childRow.read_type) }}</span>
                    <Icon icon="ep:arrow-down" style="margin-left: 4px" />
                  </div>
                </el-tag>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item
                      v-for="item in read_type_list"
                      :key="item.value"
                      :command="item.value"
                    >
                      {{ item.label }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
          <el-table-column label="列索引" align="center">
            <template #default="{ row: childRow }"> 第 {{ childRow.col_index }} 列 </template>
          </el-table-column>
        </el-table>
      </template>
    </el-table-column>
    <el-table-column label="文件名" prop="name" />
    <el-table-column label="上传时间" prop="updated_time" align="center" />
    <el-table-column label="文件大小" align="center">
      <template #default="{ row }">
        {{ convertByte(row.size) }}
      </template>
    </el-table-column>
    <el-table-column label="行数" prop="rows" align="center" />
    <el-table-column align="center" header-align="center" label="操作" width="150">
      <template #default="{ row }">
        <el-button
          link
          type="primary"
          v-tc="{
            name: '下载文件',
            file_id: row.id
          }"
          @click="fileDownload(row.name)"
        >
          下载
        </el-button>
        <el-button
          link
          type="primary"
          v-tc="{
            name: '删除文件',
            file_id: row.id
          }"
          @click="fileDelete(row.id)"
        >
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped lang="less">
.cell-wrapper {
  display: flex;
  align-items: center;
}

.cell-edit-icon {
  display: none;
}

.cell-wrapper:hover .cell-edit-icon {
  display: inline-block;
}
</style>

<script setup lang="ts">
import { Variable } from '@/views/Plan/scene/case-data'
import VariableService from '@/api/scene/variable'

const variableService = new VariableService()

let variable_list = defineModel<Array<Variable>>([])
const props = defineProps({
  sceneId: {
    type: Number,
    default: 0
  }
})

const initVariable = () => {
  variable_list.value.push({
    variable_name: '',
    variable_val: '',
    remark: ''
  })
}

const updateVariable = async (row: Variable) => {
  const _row = { ...row }
  _row.scene_id = props.sceneId
  if (_row.variable_id === undefined) {
    row.variable_id = await variableService.createSceneVariable(_row)
    return
  }
  if (_row.variable_name === '' && _row.variable_val === '' && _row.remark === '') {
    variableService.deleteSceneVariable(row.variable_id)
    return
  }
  variableService.updateSceneVariable(_row)
}

const deleteVariable = (row, index) => {
  variable_list.value.splice(index, 1)
  if (row.variable_id !== undefined) {
    variableService.deleteSceneVariable(row.variable_id)
  }
}
</script>

<template>
  <el-table :data="variable_list" border style="width: 100%">
    <el-table-column header-align="center" label="变量名" width="200">
      <template #default="scope">
        <el-input
          v-model="scope.row.variable_name"
          size="small"
          @change="updateVariable(scope.row)"
        />
      </template>
    </el-table-column>
    <el-table-column header-align="center" label="变量值" width="200">
      <template #default="scope">
        <el-input
          v-model="scope.row.variable_val"
          size="small"
          @change="updateVariable(scope.row)"
        />
      </template>
    </el-table-column>
    <el-table-column prop="remark" header-align="center" label="备注" width="200">
      <template #default="scope">
        <el-input v-model="scope.row.remark" size="small" @change="updateVariable(scope.row)" />
      </template>
    </el-table-column>
    <el-table-column header-align="center" align="center" width="50">
      <template #header>
        <Icon icon="ep:circle-plus-filled" class="icon-button" @click="initVariable" />
      </template>
      <template #default="scope">
        <Icon
          icon="ep:delete"
          class="icon-button"
          @click="deleteVariable(scope.row, scope.$index)"
        />
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped lang="less">
:deep(.el-table__cell) {
  padding: 3px;
}

:deep(.el-table__cell:first-child .cell) {
  padding-left: 0;
}

:deep(.cell) {
  padding: 0;
}

.icon-button {
  color: #409eff;
}
</style>

<script setup lang="ts">
import { defineModel } from 'vue'
import { OperationType } from '@/constants'
import type { CaseData } from '@/views/Plan/scene/case-data'
import { useIcon } from '@/hooks/web/useIcon'
import CaseService from '@/api/scene/case'

const caseService = new CaseService()

const case_data = defineModel<CaseData>({})

const emit = defineEmits(['operateCase'])

const infoFillIcon = useIcon({ icon: 'ep:info-filled', color: 'red' })

const updateCaseSwitch = (data) => {
  caseService.updateSceneCase(data).then(() => emit('operateCase', OperationType.DELETE, null))
}

const editCase = (data) => {
  emit('operateCase', OperationType.EDIT, { ...data })
}
const copyCase = (data) => {
  const _data = { ...data }
  delete _data.case_id
  delete _data.sort
  _data.extend.case_name = `${_data.title}_COPY`
  _data.extend.variable_form.forEach((vrb) => {
    if (vrb.value !== '') {
      vrb.value = `${vrb.value}_COPY`
    }
  })
  caseService.createSceneCase(_data).then(() => emit('operateCase', OperationType.COPY, null))
}
const deleteCase = (case_id) => {
  caseService.deleteSceneCase(case_id).then(() => emit('operateCase', OperationType.DELETE, null))
}
</script>

<template>
  <div class="http-div">
    <div class="case-info" @click="editCase(case_data)">
      <el-tag size="small">HTTP请求</el-tag>
      <el-tooltip placement="top-start" effect="light">
        <template #content> {{ case_data?.title }} </template>
        <div class="name-button">{{ case_data?.title }}</div>
      </el-tooltip>
      <div class="method">{{ case_data?.extend.method_type }}</div>
      <div class="url">{{ case_data?.extend.url }}</div>
    </div>
    <div class="operate-div">
      <el-tooltip effect="dark" :open-delay="200" content="禁用" placement="top">
        <div class="operate-box">
          <el-switch
            v-model="case_data.disabled"
            class="switch"
            @change="updateCaseSwitch(case_data)"
          />
        </div>
      </el-tooltip>
      <el-tooltip effect="dark" :open-delay="200" content="复制" placement="top">
        <Icon icon="ep:copy-document" class="icon-button" @click="copyCase(case_data)" />
      </el-tooltip>
      <el-popconfirm
        confirm-button-text="好的"
        cancel-button-text="不用了"
        :icon="infoFillIcon"
        title="确定删除用例？"
        @confirm="deleteCase(case_data?.case_id)"
      >
        <template #reference>
          <div class="operate-box">
            <el-tooltip effect="dark" :open-delay="200" content="删除" placement="top">
              <Icon icon="ep:delete" class="icon-button" />
            </el-tooltip>
          </div>
        </template>
      </el-popconfirm>
    </div>
  </div>
</template>

<style scoped lang="less">
.http-div {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .case-info {
    width: calc(100% - 100px);
    display: flex;
    align-items: center;

    .name-button {
      margin-left: 20px;
      width: 18%;
      color: #409eff;

      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      text-align: left;
      font-weight: 800;
    }

    .method {
      margin-left: 20px;
      width: 40px;
      color: #e6a23c;
      font-weight: 500;
      font-size: 16px;
    }

    .url {
      margin-left: 20px;
      width: 40%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }

  .operate-div {
    width: 100px;
    display: flex;
    align-items: center;

    .operate-box {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 16px;

      .switch {
        zoom: 0.8;
        --el-switch-on-color: #edf2fc;
        --el-switch-off-color: #13ce66;
      }
    }

    .icon-button {
      color: #409eff;
      margin-left: 15px;
    }

    .icon-button-not-valid {
      margin-left: 15px;
      color: #c8c9cc;
      cursor: not-allowed;
    }
  }
}
</style>

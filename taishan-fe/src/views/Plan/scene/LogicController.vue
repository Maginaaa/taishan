<script setup lang="ts">
import { defineModel, unref } from 'vue'
import CaseService from '@/api/scene/case'
import type { CaseData } from '@/views/Plan/scene/case-data'
import { useIcon } from '@/hooks/web/useIcon'
import { OperationType } from '@/constants'

const caseService = new CaseService()
const case_data = defineModel<any>({})
const infoFillIcon = useIcon({ icon: 'ep:info-filled', color: 'red' })

const emit = defineEmits(['operateCase'])

const logic_type_list = [
  { label: '次数循环', value: 1, disabled: false },
  { label: 'If判断', value: 2, disabled: false }
]
const checking_rule_list = [
  {
    value: 1,
    label: '等于',
    disabled: false
  },
  {
    value: 2,
    label: '不等于'
  },
  {
    value: 3,
    label: '包含'
  },
  {
    value: 4,
    label: '不包含'
  },
  {
    value: 5,
    label: '大于'
  },
  {
    value: 6,
    label: '小于'
  },
  {
    value: 7,
    label: '大于等于'
  },
  {
    value: 8,
    label: '小于等于'
  }
]

const updateCase = () => {
  caseService.updateSceneCase(unref(case_data))
}
const copyItem = () => {}
const deleteCase = async (id: number) => {
  await caseService.deleteSceneCase(id)
  await emit('operateCase', OperationType.DELETE, null)
}
</script>

<template>
  <div class="controller-div">
    <div class="case-info">
      <el-tag size="small" type="warning">逻辑控制器</el-tag>
      <el-select
        v-model="case_data.extend.control_type"
        :disabled="case_data.disabled"
        class="type-select"
        placeholder="请选择"
        @change="updateCase"
      >
        <el-option
          v-for="option in logic_type_list"
          :key="option.value"
          :label="option.label"
          :value="option.value"
          :disabled="option.disabled"
        />
      </el-select>
      <el-input
        v-if="case_data.extend.control_type === 1"
        v-model="case_data.extend.control_val"
        :disabled="case_data.disabled"
        class="value-input"
        placeholder="请输入"
        @change="updateCase"
      />
      <div v-if="case_data.extend.control_type === 2" class="if-control-div">
        <el-input
          v-model="case_data.extend.param_one"
          clearable
          class="param-input"
          placeholder="param1"
          @change="updateCase"
        />
        <el-select
          v-model="case_data.extend.checking_rule"
          class="param-input"
          @change="updateCase"
        >
          <el-option
            v-for="it in checking_rule_list"
            :key="it.value"
            :disabled="it.disabled"
            :label="it.label"
            :value="it.value"
          />
        </el-select>
        <el-input
          v-model="case_data.extend.param_two"
          clearable
          class="param-input"
          placeholder="param2"
          @change="updateCase"
        />
      </div>
    </div>
    <div class="operate-div">
      <el-tooltip effect="dark" :open-delay="200" content="禁用" placement="top">
        <div class="operate-box">
          <el-switch v-model="case_data.disabled" class="switch" @change="updateCase" />
        </div>
      </el-tooltip>
      <el-tooltip effect="dark" :open-delay="200" content="复制" placement="top">
        <Icon icon="ep:copy-document" class="icon-button" @click.stop="copyItem(case_data)" />
      </el-tooltip>
      <el-popconfirm
        confirm-button-text="好的"
        cancel-button-text="不用了"
        :icon="infoFillIcon"
        title="确定删除？"
        @confirm="deleteCase(case_data.case_id)"
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
.controller-div {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .case-info {
    display: flex;
    align-items: center;

    .type-select {
      margin-left: 12px;
      width: 120px;
    }

    .value-input {
      margin-left: 8px;
      width: 120px;
    }

    .if-control-div {
      display: flex;
      align-items: center;

      .param-input {
        margin-left: 8px;
        width: 160px;
      }
    }
  }

  .operate-div {
    width: 100px;
    display: flex;
    align-items: center;

    .switch {
      zoom: 0.8;
      --el-switch-on-color: #edf2fc;
      --el-switch-off-color: #13ce66;
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

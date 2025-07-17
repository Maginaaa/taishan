<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { onActivated, reactive, ref } from 'vue'
import { baseReportSearch } from '@/views/Report/report-data'
import type { ReportSearchReq } from '@/views/Report/report-data'
import ReportSearch from '@/views/Report/component/ReportSearch.vue'
import { cloneDeep } from 'lodash-es'
import type { User } from '@/views/Plan/user-data'
import Pinyin from 'pinyin-match'
import AccountService from '@/api/account'
import PlanService from '@/api/scene/plan'
import { SimplePlan } from '@/views/Plan/plan-data'

const accountService = new AccountService()
const planService = new PlanService()

let search = reactive<ReportSearchReq>(cloneDeep(baseReportSearch))

let user_list: User[] = []
let user_list_options = ref<User[]>([])
const getUserList = async () => {
  const res = await accountService.getUserList()
  user_list = [...res]
  user_list_options.value = [...res]
}
const creatorMatchFun = (val) => {
  if (val) {
    // 定义一个空数组用来存储过滤后的数据
    const result: User[] = []
    // 开始循环过滤内容
    user_list.forEach((i) => {
      // 调用 PinyinMatch.match 方法进行拼音与汉字匹配
      const m = Pinyin.match(i.name, val)
      if (m) {
        // 匹配成功则push到result数组中
        result.push(i)
      }
    })
    // 将过滤后的数组重新赋给下拉列表数据
    user_list_options.value = result
  } else {
    // 如果输入框为空, 则将下拉列表数据还原
    user_list_options.value = user_list
  }
}

let plan_list = ref<SimplePlan[]>([])
const getAllPlan = async () => {
  const res = await planService.getAllPlanList()
  plan_list.value = res
}

const reportSearchRef = ref()

getUserList()

onActivated(() => {
  getAllPlan()
  search.time = new Date().getTime()
})
</script>

<template>
  <ContentWrap>
    <el-form :inline="true" :model="search" class="inline-form">
      <el-form-item label="创建人" class="form">
        <el-select
          v-model="search.create_user_id"
          class="create-user-sel"
          clearable
          filterable
          placeholder="请选择"
          :filter-method="creatorMatchFun"
        >
          <el-option
            v-for="item in user_list_options"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="计划" class="form">
        <el-select
          v-model="search.plan_id"
          class="create-user-sel"
          clearable
          filterable
          placeholder="请选择"
        >
          <el-option
            v-for="item in plan_list"
            :key="item.plan_id"
            :label="item.plan_name"
            :value="item.plan_id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="时间范围" class="form">
        <div style="display: flex">
          <el-date-picker
            v-model="search.start_time"
            value-format="YYYY-MM-DD hh:mm:ss"
            type="datetime"
            placeholder="开始时间"
          />
          <span style="padding: 0 8px 0 8px">-</span>
          <el-date-picker
            v-model="search.end_time"
            value-format="YYYY-MM-DD HH:mm:ss"
            type="datetime"
            placeholder="结束时间"
          />
        </div>
      </el-form-item>
    </el-form>
  </ContentWrap>
  <content-wrap class="report-list">
    <report-search ref="reportSearchRef" v-model="search" />
  </content-wrap>
</template>

<style scoped lang="less">
.inline-form {
  display: flex;
  align-items: center;

  .form {
    margin: 0 0 0 20px;
  }

  .create-user-sel {
    width: 168px;
  }
}

.report-list {
  :deep(.el-card__body) {
    padding: 10px;
    width: calc(100% - 20px);
  }

  margin-top: 20px;
}
</style>

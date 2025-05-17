<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { onActivated, reactive, ref, unref } from 'vue'
import { useIcon } from '@/hooks/web/useIcon'
import PlanService from '@/api/scene/plan'
import ReportService from '@/api/scene/report'
import AccountService from '@/api/account'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import Pinyin from 'pinyin-match'
import type { User } from '@/views/Plan/user-data'
import { OperationType } from '@/constants'
import { trackClick } from '@/point/utils.js'
import TagService from '@/api/scene/tag'

const router = useRouter()
const planService = new PlanService()
const reportService = new ReportService()
const accountService = new AccountService()
const tagService = new TagService()

let table_loading = ref(false)
let plan_list = ref([])
let total = ref(0)
const search = reactive({
  plan_info: '',
  create_user_id: null,
  case_info: '',
  page: 1,
  page_size: 10,
  tag_list: null
})
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

let tag_list_options = ref<any>([])
const getPlanTagList = async () => {
  const res = await tagService.GetTagList(1) // 1 为plan的tag类型
  tag_list_options.value = [...res]
}

const plus = useIcon({ icon: 'ep:plus', color: '#fff' })
const searchPlanList = async () => {
  table_loading.value = true
  try {
    const res = await planService.getPlanList(search)
    plan_list.value = res.list
    total.value = res.total
  } finally {
    table_loading.value = false
  }
}

const handleEdit = (planId) => {
  router.push({
    name: 'PlanDetail',
    params: {
      id: planId
    }
  })
}

const handleCopy = (row) => {
  create_plan_form = reactive({
    plan_id: row.plan_id,
    plan_name: row.plan_name,
    remark: row.remark
  })
  create_plan_type.value = OperationType.COPY
  create_dialog_visible.value = true
}

const handleDelete = async (id) => {
  await planService.deletePlan(id)
  await searchPlanList()
}

const planExecute = async (plan_id) => {
  const id = await planService.executePlan(unref(plan_id))
  if (id !== undefined) {
    await router.push({
      name: 'ReportDetail',
      params: {
        id: id
      }
    })
  } else {
    ElMessage.error('执行失败')
  }
}

const stopPress = async (report_id) => {
  await reportService.stopPress(report_id)
  await searchPlanList()
}

const updatePlan = (row) => {
  row.is_editing = false
  planService.updatePlanSimple(row)
}

let create_btn_loading = ref(false)
const create_dialog_visible = ref(false)
let create_plan_type = ref<number>(0)
let create_plan_form = reactive({
  plan_id: null,
  plan_name: '',
  remark: ''
})
const formRef = ref()
const resetForm = () => {
  formRef.value.resetFields()
}
const preCreatePlan = () => {
  create_plan_form = reactive({
    plan_id: null,
    plan_name: '',
    remark: ''
  })
  create_plan_type.value = OperationType.CREATE
  create_dialog_visible.value = true
}

const CreatePlan = async () => {
  if (create_plan_form.plan_name.trim() === '') {
    ElMessage.error('计划名称不能为空')
    return
  }
  create_btn_loading.value = true
  try {
    const planId =
      create_plan_type.value === OperationType.COPY
        ? await planService.copyPlan(create_plan_form)
        : await planService.createPlan(create_plan_form)
    if (planId) {
      create_dialog_visible.value = false
      await router.push({
        name: 'PlanDetail',
        params: {
          id: planId
        }
      })
    } else {
      ElMessage.error('创建失败')
    }
  } finally {
    create_btn_loading.value = false
  }
}

getUserList()
getPlanTagList()
searchPlanList()

onActivated(() => {
  searchPlanList()
})
</script>

<template>
  <ContentWrap>
    <el-form :inline="true" :model="search" class="inline-form">
      <div>
        <el-form-item label="创建人" class="form">
          <el-select
            v-model="search.create_user_id"
            class="create-user-sel"
            clearable
            filterable
            placeholder="请选择"
            :filter-method="creatorMatchFun"
            @change="searchPlanList"
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
          <el-input
            v-model="search.plan_info"
            placeholder="计划名/ID"
            clearable
            @change="searchPlanList"
          />
        </el-form-item>
        <el-form-item label="接口" class="form">
          <el-input
            v-model="search.case_info"
            placeholder="接口名/路径"
            clearable
            @change="searchPlanList"
          />
        </el-form-item>
        <el-form-item label="标签" class="form">
          <el-select
            v-model="search.tag_list"
            class="create-user-sel"
            clearable
            filterable
            multiple
            placeholder="请选择"
            @change="searchPlanList"
          >
            <el-option
              v-for="item in tag_list_options"
              :key="item.id"
              :label="item.label"
              :value="item.id"
            >
              <el-tag>{{ item.label }}</el-tag>
            </el-option>
            <template #tag>
              <el-tag v-for="tag in search.tag_list" :key="tag">
                {{ tag_list_options.find((opt) => opt.id === tag)?.label || tag }}
              </el-tag>
            </template>
          </el-select>
        </el-form-item>
      </div>
      <el-form-item class="form" style="float: right">
        <el-button type="primary" v-tc="{ name: '计划查询' }" @click="searchPlanList">
          查询
        </el-button>
      </el-form-item>
    </el-form>
  </ContentWrap>
  <ContentWrap class="plan-list">
    <div class="operate-row">
      <el-button type="primary" :icon="plus" v-tc="{ name: '预创建' }" @click="preCreatePlan">
        创建
      </el-button>
    </div>
    <el-table v-loading="table_loading" :data="plan_list" resizable border stripe class="table">
      <el-table-column prop="plan_id" label="测试计划ID" align="center" header-align="center" />
      <el-table-column prop="plan_name" label="测试计划名称" align="center" header-align="center" />
      <el-table-column label="状态" align="center" header-align="center">
        <template #default="scope">
          <el-tag v-if="scope.row.is_running" type="success"> 运行中 </el-tag>
          <el-tag v-else type="info"> 未运行 </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="create_user_name"
        label="创建人"
        align="center"
        header-align="center"
      />
      <el-table-column label="标签" align="center" header-align="center">
        <template #default="{ row }">
          <el-tag v-for="tag in row.tag_list" :key="tag" style="margin-left: 4px">
            {{ tag_list_options.find((opt) => opt.id === tag)?.label || tag }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="update_time"
        label="最后更新时间"
        align="center"
        header-align="center"
      />
      <el-table-column prop="remark" label="备注" align="center" header-align="center">
        <template #default="scope">
          <el-input
            v-if="scope.row.is_editing"
            v-model="scope.row.remark"
            maxlength="200"
            show-word-limit
            autofocus
            @blur="updatePlan(scope.row)"
          />
          <span v-else @click="scope.row.is_editing = true" class="remark">
            {{ scope.row.remark ? scope.row.remark : '--' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column fixed="right" align="center" header-align="center" label="操作" width="220">
        <template #default="{ row }">
          <el-button
            v-if="!row.is_running"
            :disabled="!row.debug_success"
            type="primary"
            link
            v-tc="{ name: '启动', plan_id: row.plan_id }"
            @click="planExecute(row.plan_id)"
          >
            启动
          </el-button>
          <el-button
            v-else
            link
            type="primary"
            v-tc="{ name: '停止', plan_id: row.plan_id }"
            @click="stopPress(row.report_id)"
          >
            停止
          </el-button>
          <el-button
            link
            type="primary"
            v-tc="{ name: '编辑', plan_id: row.plan_id }"
            @click="handleEdit(row.plan_id)"
            >编辑</el-button
          >
          <el-button
            link
            type="primary"
            v-tc="{ name: '复制', plan_id: row.plan_id }"
            @click="handleCopy(row)"
            >复制</el-button
          >
          <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.plan_id)">
            <template #reference>
              <el-button
                :disabled="row.is_running"
                v-tc="{ name: '删除', plan_id: row.plan_id }"
                link
                type="primary"
                >删除</el-button
              >
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-div">
      <el-pagination
        class="plan-pag"
        small
        background
        v-model:current-page="search.page"
        v-model:page-size="search.page_size"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper "
        :total="total"
        @change="
          searchPlanList(),
            trackClick({
              name: '翻页',
              search
            })
        "
      />
    </div>
  </ContentWrap>
  <el-dialog
    v-model="create_dialog_visible"
    :title="create_plan_type === 0 ? '创建测试计划' : '复制计划'"
    draggable
    @close="resetForm"
    width="30%"
  >
    <el-form ref="formRef" :model="create_plan_form" label-position="right">
      <el-form-item label="计划名" prop="plan_name" label-width="60px">
        <el-input
          v-model="create_plan_form.plan_name"
          maxlength="40"
          show-word-limit
          autocomplete="off"
          clearable
        />
      </el-form-item>
      <el-form-item label="备 注" prop="remark" label-width="60px">
        <el-input v-model="create_plan_form.remark" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div>
        <el-button v-tc="{ name: '取消创建' }" @click="create_dialog_visible = false">
          取 消
        </el-button>
        <el-button
          type="primary"
          v-tc="{ name: '创建' }"
          :loading="create_btn_loading"
          @click="CreatePlan"
        >
          创 建
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="less">
.inline-form {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .form {
    margin: 0 0 0 20px;
  }
  .search-button {
    margin-right: 20px;
  }

  .create-user-sel {
    width: 168px;
  }
}

.plan-list {
  :deep(.el-card__body) {
    padding: 10px;
    width: 100%;
  }

  display: flex;

  margin-top: 20px;
  width: 100%;
  height: calc(100% - 100px);

  .operate-row {
    height: 35px;
    display: flex;
    justify-content: flex-end;
  }

  .table {
    width: 100%;
    margin-top: 7px;

    .button-spacing {
      display: inline-block;
      width: 6px;
    }

    .remark {
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
}

.pagination-div {
  margin-top: 10px;

  .plan-pag {
    margin: 10px 0;
    float: right;
  }
}
</style>

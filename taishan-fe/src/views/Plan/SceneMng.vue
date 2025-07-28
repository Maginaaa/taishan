<script setup lang="ts">
import { nextTick, ref, unref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useIcon } from '@/hooks/web/useIcon'
import CaseTree from '@/views/Plan/scene/CaseTree.vue'
import { HttpCaseExtend, SceneInfo } from '@/views/Plan/scene/case-data'
import { CaseType, OperationType } from '@/constants'
import { ElMessage } from 'element-plus'
import SceneEdit from '@/views/Plan/component/SceneEdit.vue'
import SceneService from '@/api/scene/scene'
import CaseService from '@/api/scene/case'
import HttpCase from '@/views/Plan/cases/HttpCase.vue'
import VariableTable from '@/views/Plan/component/VariableTable.vue'
import DataExport from '@/views/Plan/scene/DataExport.vue'
import { TreeNode } from '@/views/Plan/plan-data'
import PlanService from '@/api/scene/plan'

const sceneService = new SceneService()
const caseService = new CaseService()
const planService = new PlanService()

const props = defineProps({
  planId: Number
})
let scene_list = defineModel<Array<SceneInfo> | any>()
const emit = defineEmits(['refresh'])
const refreshSceneList = () => {
  emit('refresh')
}

// scene-header
const moreFillIcon = useIcon({ icon: 'ep:more-filled' })
const emptyScene = {
  plan_id: props.planId,
  scene_name: '',
  scene_type: 0
}

const operateScene = async (tp, scene) => {
  if (scene) {
    scene.plan_id = props.planId
  }
  switch (tp) {
    case OperationType.CREATE:
      await sceneEditRef.value.preEditScene(tp, scene)
      break
    case OperationType.EDIT:
      await sceneEditRef.value.preEditScene(tp, scene)
      break
    case OperationType.COPY:
      await sceneService.copyScene(scene.scene_id)
      break
    case OperationType.ENABLED:
      scene.disabled = false
      await sceneService.updateScene(scene)
      break
    case OperationType.DISABLED:
      scene.disabled = true
      await sceneService.updateScene(scene)
      break
    case OperationType.DELETE:
      await sceneService.deleteScene(scene.scene_id)
      break
    default:
      ElMessage.error('错误操作')
      break
  }
  await refreshSceneList()
}

// scene-body
const empty_case_extend: HttpCaseExtend = {
  case_name: '新的HTTP请求',
  url: 'https://',
  method_type: 'POST',
  use_scene_variable: false,
  body: {
    body_type: 'json',
    body_value: ''
  },
  params_form: [
    {
      enable: true,
      key: '',
      value: '',
      desc: ''
    }
  ],
  headers_form: [
    {
      enable: true,
      key: '',
      value: '',
      desc: ''
    }
  ],
  assert_form: [
    {
      enable: true,
      extract_type: 1,
      extract_express: '',
      checking_rule: 1,
      expect_value: '',
      desc: ''
    }
  ],
  variable_form: [
    {
      enable: true,
      extract_type: 1,
      key: '',
      value: ''
    }
  ],
  rally_point: {
    enable: false,
    concurrency: 0,
    timeout_period: 0
  },
  waiting_config: {
    pre_waiting_switch: false,
    pre_waiting_time: 0,
    post_waiting_switch: false,
    post_waiting_time: 0
  },
  overtime_config: {
    enable: false,
    timeout_period: 0
  }
}

const preCreateCase = async (scene_id: number, type: number) => {
  if (type === CaseType.HttpCase) {
    httpCaseRef.value.drawerOpen(OperationType.CREATE, {
      scene_id: scene_id,
      type: CaseType.HttpCase,
      extend: empty_case_extend
    })
    return
  }
  if (type === CaseType.LogicControl) {
    await caseService.createSceneCase({
      scene_id: scene_id,
      type: CaseType.LogicControl,
      extend: {
        control_val: '',
        control_type: ''
      }
    })
    await refreshSceneList()
    return
  }
}
const operateCase = (opType: number, data) => {
  if (
    opType === OperationType.DELETE ||
    opType === OperationType.CREATE ||
    opType === OperationType.COPY
  ) {
    refreshSceneList()
    return
  }
  httpCaseRef.value.drawerOpen(opType, data)
}

let plan_list = ref<TreeNode[]>([])

const getAllPlan = async () => {
  const _res = await planService.getAllPlanList()
  let _plan_list: any = []
  _res.forEach((item) => {
    _plan_list.push({ label: item.plan_name, value: item.plan_id })
  })
  plan_list.value = _plan_list
  importRef.value++
}

const importRef = ref(0)
const sceneCaseMap = {}
const planSceneMap = {}
const selectCase = {
  emitPath: false,
  lazy: true,
  async lazyLoad(node, resolve) {
    const { level, data } = node
    let res: TreeNode[] = []
    switch (level) {
      case 0:
        res = plan_list.value
        break
      case 1:
        const plan_id = data.value
        if (!planSceneMap.hasOwnProperty(plan_id)) {
          const scenes_list = await sceneService.getSceneList(plan_id)
          let scenes: TreeNode[] = []
          scenes_list.forEach((scene) => {
            let cases: TreeNode[] = []
            scene.case_tree.forEach((cs) => {
              cases.push({
                label: cs.title,
                value: cs.case_id,
                leaf: true
              })
            })
            sceneCaseMap[scene.scene_id] = cases
            scenes.push({
              label: scene.scene_name,
              value: scene.scene_id,
              leaf: false
            })
          })
          planSceneMap[plan_id] = scenes
        }
        res = planSceneMap[plan_id]
        break
      case 2:
        res = sceneCaseMap[data.value]
        break
      default:
        break
    }
    resolve(res)
    return
  }
}

const chooseImportCase = async (case_id, scene_id) => {
  await caseService.importCase({
    case_id,
    scene_id
  })
  await refreshSceneList()
}

const plusIcon = useIcon({ icon: 'ep:plus' })

const sceneEditRef = ref()
const httpCaseRef = ref()

const sceneSort = () => {
  nextTick(async () => {
    scene_list.value.forEach((item, index) => {
      item.sort = item.scene_type === 0 ? index + 1 : 0
    })
    await sceneService.updateSceneSort(unref(scene_list))
    await refreshSceneList()
  })
}
getAllPlan()
</script>

<template>
  <vue-draggable
    target=".scene-container"
    handle=".handle"
    filter=".unHandle"
    v-model="scene_list"
    :animation="1000"
    @end="sceneSort"
  >
    <div class="scene-container">
      <el-card
        v-for="scene in scene_list"
        :key="scene.scene_id"
        :class="[scene.disabled ? 'not-valid' : '', 'scene-card']"
      >
        <template #header>
          <div class="scene-head-div">
            <div class="flex items-center">
              <Icon
                icon="ep:rank"
                :class="scene.scene_type === 0 ? ' cursor-move handle' : 'unHandle'"
                :style="scene.scene_type === 0 ? 'color: #409eff' : 'color: #b1b3b8'"
              />
              <div class="scene-name-box">
                <div class="scene-head-div-left">
                  <el-badge v-if="scene.scene_type === 1" type="warning" value="前置" class="item">
                    <span>{{ scene.scene_name }}</span>
                  </el-badge>
                  <span v-else>{{ scene.scene_name }}</span>
                </div>
                <div class="operate-row">
                  <el-tag
                    size="small"
                    class="operate-col"
                    @click="preCreateCase(scene.scene_id, CaseType.HttpCase)"
                    v-tc="{ name: '添加请求', scene_id: scene.scene_id }"
                  >
                    + 添加请求
                  </el-tag>
                  <el-tag
                    size="small"
                    class="operate-col"
                    type="warning"
                    @click="preCreateCase(scene.scene_id, CaseType.LogicControl)"
                    v-tc="{ name: '添加逻辑控制器', scene_id: scene.scene_id }"
                  >
                    + 逻辑控制器
                  </el-tag>
                  <el-popover placement="bottom-start" width="676" trigger="click">
                    <variable-table v-model="scene.variable_list" :scene-id="scene.scene_id" />
                    <template #reference>
                      <el-tag
                        size="small"
                        class="operate-col"
                        type="success"
                        v-tc="{ name: '添加场景变量', scene_id: scene.scene_id }"
                      >
                        + 场景变量
                      </el-tag>
                    </template>
                  </el-popover>
                  <el-dropdown trigger="click">
                    <el-tag
                      size="small"
                      class="operate-col"
                      v-tc="{ name: '导入请求', scene_id: scene.scene_id }"
                    >
                      + 导入请求
                    </el-tag>
                    <template #dropdown>
                      <el-cascader-panel
                        :key="importRef"
                        :props="selectCase"
                        @change="(val) => chooseImportCase(val, scene.scene_id)"
                      />
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </div>
            <el-dropdown @command="(cmd) => operateScene(cmd, scene)">
              <el-button :icon="moreFillIcon" size="small" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="1">编辑</el-dropdown-item>
                  <el-dropdown-item :command="2">复制</el-dropdown-item>
                  <el-dropdown-item v-if="!scene.disabled" :command="11">禁用</el-dropdown-item>
                  <el-dropdown-item v-else :command="12">启用</el-dropdown-item>
                  <el-dropdown-item :command="3">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
        <case-tree ref="apiTree" :scene-data="scene" @operate-case="operateCase" />
        <data-export
          v-if="scene.scene_type === 1"
          :class="scene.disabled ? 'not-valid' : ''"
          :scene-data="scene"
          :plan-id="planId"
        />
      </el-card>
    </div>
  </vue-draggable>
  <el-affix position="bottom" style="width: 100%; height: 30px" :offset="20">
    <el-button
      type="primary"
      :icon="plusIcon"
      plain
      @click="operateScene(0, emptyScene)"
      v-tc="{ name: '预添加场景', plan_id: props.planId }"
      class="create-scene-btn"
    >
      添加场景
    </el-button>
  </el-affix>
  <scene-edit ref="sceneEditRef" @refresh="refreshSceneList" />
  <http-case ref="httpCaseRef" @refresh="refreshSceneList" />
</template>

<style scoped lang="less">
.not-valid {
  background-color: #f5f5f5;
  border: #e9e9eb solid 1px;
}

@keyframes fadeInLeft {
  0% {
    opacity: 0;
    -webkit-transform: translate3d(-100%, 0, 0);
    transform: translate3d(-100%, 0, 0);
  }

  to {
    opacity: 1;
    -webkit-transform: translateZ(0);
    transform: translateZ(0);
  }
}

.scene-card {
  width: 100%;
  margin-top: 6px;
  margin-bottom: 6px;

  :deep(.el-card__header) {
    height: 40px;
    padding: 0;
  }
  :deep(.el-card__body) {
    padding: 4px;
  }
  .scene-head-div {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px 0 10px;

    .scene-name-box {
      display: flex;
      align-items: center;
    }

    .scene-head-div-left {
      display: flex;
      align-items: center;
      margin-left: 6px;
      gap: 10px;
    }

    .operate-row {
      width: 200px;

      display: none;
      align-items: center;
      margin-left: -3px;

      .operate-col {
        margin-left: 10px;
        cursor: pointer;
      }
    }

    .scene-name-box:hover .scene-head-div-left {
      display: none;
    }

    .scene-name-box:hover .operate-row {
      display: flex;
      animation: fadeInLeft 0.4s;
    }
  }
}

.create-scene-btn {
  width: 100%;
  //margin-bottom: 30px;
}

.animate__fadeInLeft {
  -webkit-animation-name: fadeInLeft;
  animation-name: fadeInLeft;
}
</style>

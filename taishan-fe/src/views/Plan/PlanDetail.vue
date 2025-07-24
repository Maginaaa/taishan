<script setup lang="ts">
import { useRouter } from 'vue-router'

import { nextTick, reactive, ref, unref } from 'vue'
import { useIcon } from '@/hooks/web/useIcon'
import { EChartsOption } from 'echarts'
import { ElMessage } from 'element-plus'
import { Icon } from '@/components/Icon'
import { Echart } from '@/components/Echart'
import { PressType, SourceName } from '@/constants'
import { cloneDeep } from 'lodash-es'
import { isNumber } from '@/utils/is'
import 'vue3-cron-plus-picker/style.css'

import PlanService from '@/api/scene/plan'
import MachineService from '@/api/scene/machine'
import SceneService from '@/api/scene/scene'
import K8sService from '@/api/scene/k8s'

import SceneMng from '@/views/Plan/SceneMng.vue'
import DebugResult from '@/views/Plan/plan/DebugResult.vue'
import ReportSearch from '@/views/Report/component/ReportSearch.vue'
import DataConfig from '@/views/Plan/plan/DataConfig.vue'
import TargetRpsTree from './component/TargetRpsTree.vue'
import DebugResultList from '@/views/Plan/plan/DebugResultList.vue'
import OperationLog from '@/views/Plan/plan/OperationLog.vue'

import { emptyPlanInfo, pressGraphOptions, TreeNode } from '@/views/Plan/plan-data'
import { baseReportSearch, ReportSearchReq } from '@/views/Report/report-data'
import type { SceneInfo } from '@/views/Plan/scene/case-data'
import TagService from '@/api/scene/tag'

const router = useRouter()
const planService = new PlanService()
const sceneService = new SceneService()
const machineService = new MachineService()
const k8sService = new K8sService()
const tagService = new TagService()

const plan_id = ref<number>(0)
let scene_list = ref<Array<SceneInfo>>([])
const targetRpsTree = ref()
const getPlanId = () => {
  plan_id.value = Number(router.currentRoute.value.params.id)
}
const getMachineList = async () => {
  const res = await machineService.getMachineList()
  machine_list_len.value = Object.keys(res).length
}

const plan_loading = ref(false)
const plan_debug_pass = ref(false)
const machine_list_len = ref(0)
let plan_info = ref<any>(emptyPlanInfo)

const getPlanDetail = async () => {
  plan_loading.value = true
  try {
    const plan_dt = await planService.getPlanDetail(unref(plan_id))
    if (plan_dt) {
      plan_info.value = plan_dt
      plan_debug_pass.value = plan_dt.debug_success
    }
    await getPlanScenes()
    await refreshPressScene()
  } finally {
    plan_loading.value = false
  }
}

const tag_list_options = ref<any>([])
const getPlanTagList = async () => {
  const res = await tagService.GetTagList(1) // 1 为plan的tag类型
  tag_list_options.value = [...res]
}

const getPlanScenes = async () => {
  const scenes_dt = await sceneService.getSceneList(unref(plan_id))
  scene_list.value.length = 0
  if (scenes_dt) {
    scene_list.value = scenes_dt
  }
  setTimeout(() => {
    targetRpsTree.value?.channelTargetRps()
  }, 500)
  await refreshPressScene()
}
const refreshPressScene = () => {
  let arr: Array<any> = []
  for (let s of scene_list.value) {
    if (s.disabled) {
      continue
    }
    const scene = plan_info.value.press_info.scene_list.find(
      (scene) => scene.scene_id === s.scene_id
    )
    if (scene) {
      scene.scene_name = s.scene_name
      scene.scene_type = s.scene_type
      arr.push(scene)
    } else {
      arr.push({
        scene_id: s.scene_id,
        scene_name: s.scene_name,
        scene_type: s.scene_type,
        concurrency: 0,
        rate: 0,
        case_tree: []
      })
    }
  }
  plan_info.value.press_info.scene_list = arr
}

const title_edit = ref(false)
const planNameRef = ref()
const handleEditTitle = async () => {
  title_edit.value = false
  await planService.updatePlanSimple(unref(plan_info))
  const dt = await planService.getPlanDetail(unref(plan_id))
  if (dt) {
    plan_info.value = dt
  }
}
const preEditTitle = () => {
  title_edit.value = true
  nextTick(() => planNameRef.value.focus())
}

const active_tab = ref('scene_config')
const chooseHeaderTabs = (tabName: string) => {
  if (tabName === 'press_config') {
    drawPressGraph()
  }
  if (tabName === 'press_history') {
    getReportList()
  }
  if (tabName === 'other_config') {
    getNamespaces()
  }
}

// 压测策略相关
const press_type_list = [
  {
    value: 0,
    label: '并发模式'
  },
  {
    value: 1,
    label: '阶梯式加压'
  },
  // {
  //   value: 2,
  //   label: 'RPS模式',
  //   disable: true
  // },
  {
    value: 3,
    label: 'RPS比例模式'
  },
  {
    value: 4,
    label: '固定次数'
  }
]
let optionsData = reactive<EChartsOption>(pressGraphOptions) as EChartsOption

const chartsRef = ref()

const drawPressGraph = () => {
  checkGraphParam()
  if (active_tab.value !== 'press_config') {
    return
  }
  let seriesData: Array<any> = []
  const _press_info = { ...plan_info.value.press_info }
  if (_press_info.press_type === 0) {
    // 对每个场景进行循环处理
    _press_info.scene_list
      .filter((sceneInfo) => sceneInfo.scene_type !== 1)
      .forEach((sceneInfo) => {
        const data: number[][] = []
        for (let i = 0; i <= _press_info.duration; i++) {
          const tempData: number[] = [i * 60, sceneInfo.concurrency]
          data.push(tempData)
        }
        seriesData.push({
          data: data,
          type: 'line',
          name: sceneInfo.scene_name,
          step: 'end', // 阶梯式加压展示直角
          showSymbol: false
        })
      })
    // 所有场景总数据
    const total: number[][] = []
    for (let i = 0; i < seriesData[0].data.length; i++) {
      const sum = seriesData.reduce((acc, arr) => acc + arr.data[i][1], 0)
      total.push([seriesData[0].data[i][0], sum])
    }
    seriesData.unshift({
      data: total,
      type: 'line',
      name: '总计',
      step: 'end', // 阶梯式加压展示直角
      showSymbol: false
    })
  } else if (_press_info.press_type === 1) {
    const data: number[][] = []
    const startConcurrent = _press_info.start_concurrency
    const stepSize = _press_info.step_size
    const stepDuration = _press_info.step_duration
    const maxConcurrent = _press_info.concurrency
    const duration = _press_info.duration * 60
    let time = 0
    let concurrent = startConcurrent

    // 计算每个步长的时间
    if (stepDuration > 0) {
      const steps = Math.ceil(duration / stepDuration)
      for (let i = 0; i < steps; i++) {
        data.push([time, concurrent])
        time += stepDuration
        concurrent += stepSize
        concurrent = Math.min(concurrent, maxConcurrent)
      }
      data.push([duration, concurrent])
    }
    seriesData = [
      {
        data: data,
        type: 'line',
        step: 'end',
        showSymbol: false
      }
    ]
  }

  // 配置图表
  optionsData!.series = cloneDeep(seriesData)
  chartsRef.value?.refreshHard()
}

const checkGraphParam = () => {
  if (plan_info.value.press_info.duration > 120) {
    plan_info.value.press_info.duration = 120
  } else if (plan_info.value.press_info.duration < 1) {
    plan_info.value.press_info.duration = 1
  }

  if (!plan_info.value.press_info.rps) {
    plan_info.value.press_info.rps = 1
  } else if (plan_info.value.press_info.rps > 200000) {
    plan_info.value.press_info.rps = 200000
  }
  if (plan_info.value.press_info.concurrency > 20000) {
    plan_info.value.press_info.concurrency = 20000
  } else if (plan_info.value.press_info.concurrency < 1) {
    plan_info.value.press_info.concurrency = 1
  }
  if (plan_info.value.press_info.start_concurrency < 0) {
    plan_info.value.press_info.start_concurrency = 0
  }
  if (plan_info.value.press_info.step_size < 1) {
    plan_info.value.press_info.step_size = 1
  }
  if (plan_info.value.press_info.step_duration < 1) {
    plan_info.value.press_info.step_duration = 1
  }
  plan_info.value.press_info.scene_list
    .filter((sc) => sc.scene_type !== 1)
    .forEach((item) => {
      if (item.concurrency > 20000) {
        item.concurrency = 20000
      } else if (item.concurrency < 1) {
        item.concurrency = 1
      }
      if (item.rate > 100) {
        item.rate = 100
      } else if (item.rate < 1) {
        item.rate = 1
      }
      if (item.iteration < 1) {
        item.iteration = 1
      }
    })

  if (plan_info.value.press_info.press_type === 3) {
    if (
      plan_info.value.press_info.rps <
      getPressSceneList().length * plan_info.value.engine_count * 20
    ) {
      ElMessage.warning('RPS总数不能小于场景数 * 施压机数 * 20')
      plan_info.value.press_info.rps =
        getPressSceneList().length * plan_info.value.engine_count * 20
    }
  }
}

// 高级配置
const cron_dialog_visible = ref(false)
let expression = ref('* * * * * *')
const timeIcon = useIcon({ icon: 'ep:clock' })

const chooseCron = () => {
  cron_dialog_visible.value = true
  expression.value = plan_info.value.task_info.cron
}
const fillValue = (val) => {
  plan_info.value.task_info.cron = val
}
const sampling_type_list = [
  {
    value: 0,
    label: '不采样'
  },
  {
    value: 1,
    label: '采样-仅错误'
  },
  {
    value: 2,
    label: '全量采样'
  }
]

const cas_props = {
  multiple: true,
  lazy: true,
  async lazyLoad(node, resolve) {
    const { level, data } = node
    let _res
    switch (level) {
      case 0:
        _res = namespace_list.value
        break
      case 1:
        _res = await k8sService.getDeployments({ namespace: data.value })
        break
    }
    let _list: any = []
    _res.forEach((item) => {
      _list.push({ label: item, value: item, leaf: level === 1 })
    })
    resolve(_list)
    return
  }
}

let namespace_list = ref<TreeNode[]>([])
const getNamespaces = async () => {
  const _res = await k8sService.getNamespaces()
  let _namespace_list: any = []
  _res.forEach((item) => {
    _namespace_list.push({ label: item, value: item })
  })
  namespace_list.value = _namespace_list
}

// 压测记录
let report_search = reactive<ReportSearchReq>(cloneDeep(baseReportSearch))
const getReportList = () => {
  report_search.plan_id = unref(plan_id)
}

// 底部操作栏
const plan_debug_loading = ref(false)
const execute_button_loading = ref(false)

const planSave = async () => {
  if (!planInfoCheck()) {
    return
  }
  const res = await planService.updatePlan(unref(plan_info))
  if (res.code === 0) {
    ElMessage.success('保存成功')
  }
}
const planDebug = async () => {
  plan_debug_loading.value = true
  plan_debug_pass.value = false
  let all_pass = true
  try {
    await planService.updatePlan(unref(plan_info))
    const dt = await planService.debugPlan(unref(plan_id))
    debugResultRef.value.drawerOpen(dt)
    if (dt.length === 0) {
      all_pass = false
    } else {
      dt.forEach((item) => {
        all_pass = all_pass && item.passed
      })
    }
  } catch (e) {
    all_pass = false
  } finally {
    plan_debug_loading.value = false
    plan_debug_pass.value = all_pass
  }
}

const planExecute = async () => {
  if (!planInfoCheck()) {
    return
  }
  execute_button_loading.value = true
  try {
    await planService.updatePlan(unref(plan_info))
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
  } catch (e) {
    console.log(e)
  } finally {
    execute_button_loading.value = false
  }
}

const getPressSceneList = () => {
  return plan_info.value.press_info.scene_list.filter((sc) => sc.scene_type !== 1)
}

const getSumRate = () => {
  let sum = 0
  plan_info.value.press_info.scene_list.forEach((item) => {
    if (item.scene_type === 0) {
      sum += item.rate || 0
    }
  })
  return sum
}

const planInfoCheck = () => {
  if (plan_info && plan_info.value.press_info.press_type === 3) {
    //RPS比例模式判断总流量比例
    const sList = plan_info.value.press_info.scene_list
    let sum = 0
    for (let index = 0; index < sList.length; index++) {
      const element = sList[index]
      if (element.scene_type === 1) {
        continue
      }
      sum += element.rate || 0
    }
    if (sum !== 100) {
      ElMessage.error('RPS比例模式场景流量比例为[' + sum + ']不为100，请到“压测策略中”检测')
      return false
    }
  }

  if (!(plan_info.value.press_info.duration > 0)) {
    ElMessage.error('压测时长需大于 0')
    return false
  }

  if (!(plan_info.value.engine_count > 0)) {
    ElMessage.error('请选择压测机数量')
    return false
  }

  switch (plan_info.value.press_info.press_type) {
    case PressType.ConcurrentModel:
    case PressType.RPSModel:
      if (plan_info.value.press_info.scene_list.some((item) => item.duration === 0)) {
        ElMessage.error('并发数需大于 0')
        return false
      }
      break
    case PressType.StepModel:
      if (!(plan_info.value.press_info.concurrency > 0)) {
        ElMessage.error('最大并发数需大于 0')
        return false
      }
      break
  }
  return true
}

const samplingRateCheck = () => {
  if (!isNumber(plan_info.value.sampling_info.sampling_rate)) {
    plan_info.value.sampling_info.sampling_rate = 0
  }
  if (plan_info.value.sampling_info.sampling_rate < 0) {
    plan_info.value.sampling_info.sampling_rate = 0
  }
  if (plan_info.value.sampling_info.sampling_rate > 10000) {
    plan_info.value.sampling_info.sampling_rate = 10000
  }
}

const downloadIcon = useIcon({ icon: 'ep:download' })
const caretRightIcon = useIcon({ icon: 'ep:caret-right' })
const videoPlayIcon = useIcon({ icon: 'ep:video-play' })

const debugResultRef = ref()

getPlanId()
getPlanTagList()
getPlanDetail()
getMachineList()
</script>

<template>
  <div class="plan-header">
    <el-input
      v-if="title_edit"
      ref="planNameRef"
      v-model="plan_info.plan_name"
      autofocus
      maxlength="40"
      show-word-limit
      class="title-edit-input"
      @blur="handleEditTitle"
      style="margin-left: 10px"
    />
    <div v-else class="title-div">
      <span class="title-span" @click="preEditTitle">
        {{ plan_info.plan_name }}
      </span>
      <Icon icon="ep:edit" class="icon" />
    </div>
    <div style="margin-right: 10px">
      <el-button
        type="primary"
        :icon="downloadIcon"
        v-tc="{ name: '计划保存' }"
        plain
        @click="planSave"
      >
        计划保存
      </el-button>

      <el-button
        type="primary"
        :icon="caretRightIcon"
        :loading="plan_debug_loading"
        v-tc="{ name: '计划调试' }"
        plain
        @click="planDebug"
      >
        计划调试
      </el-button>

      <el-button
        type="primary"
        :icon="videoPlayIcon"
        :loading="execute_button_loading"
        :disabled="!plan_debug_pass"
        v-tc="{ name: '开始执行' }"
        plain
        @click="planExecute"
      >
        开始执行
      </el-button>
    </div>
  </div>
  <el-skeleton :loading="plan_loading" :rows="12" animated>
    <el-tabs
      v-model="active_tab"
      class="plan-tabs"
      tabPosition="left"
      lazy
      @tab-change="chooseHeaderTabs"
    >
      <el-tab-pane label="计划信息" name="base_info">
        <el-descriptions title="计划信息" style="margin-top: 10px">
          <el-descriptions-item label="创建人">
            {{ plan_info.create_user_name }}
          </el-descriptions-item>
          <el-descriptions-item label="修改人">
            {{ plan_info.update_user_name }}
          </el-descriptions-item>
          <el-descriptions-item label="最近压测时间">
            {{ plan_info.last_press_time }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ plan_info.create_time }}</el-descriptions-item>
          <el-descriptions-item label="修改时间">{{ plan_info.update_time }}</el-descriptions-item>
          <el-descriptions-item label="累计压测次数">
            {{ plan_info.press_count }}
          </el-descriptions-item>
          <el-descriptions-item label="标签信息">
            <el-select multiple v-model="plan_info.tag_list" style="width: 300px">
              <el-option
                v-for="item in tag_list_options"
                :key="item"
                :label="item.label"
                :value="item.id"
              >
                <el-tag>{{ item.label }}</el-tag>
              </el-option>
              <template #tag>
                <el-tag v-for="tag in plan_info.tag_list" :key="tag.index">
                  {{ tag_list_options.find((opt) => opt.id === tag)?.label || tag }}
                </el-tag>
              </template>
            </el-select>
          </el-descriptions-item>
          <el-descriptions-item label="备注信息">
            <div class="remark">
              {{ plan_info.remark }}
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>
      <el-tab-pane label="场景管理" name="scene_config">
        <scene-mng :plan-id="plan_id" v-model="scene_list" @refresh="getPlanScenes" />
      </el-tab-pane>
      <el-tab-pane label="数据配置" name="data_config" class="data-cfg-container">
        <data-config v-model="plan_info" />
      </el-tab-pane>
      <el-tab-pane label="压测策略" name="press_config">
        <div class="press-cfg-div">
          <el-form
            label-position="right"
            label-width="100px"
            :model="plan_info"
            style="min-width: 520px"
          >
            <el-form-item label="压测模式">
              <el-radio-group v-model="plan_info.press_info.press_type" @change="drawPressGraph">
                <el-radio v-for="item in press_type_list" :key="item.value" :value="item.value">
                  {{ item.label }}
                </el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item
              v-show="[0, 2, 3].indexOf(plan_info.press_info.press_type) > -1"
              label="总时长"
            >
              <el-input
                v-model.number="plan_info.press_info.duration"
                class="input-width"
                placeholder="输入范围： 1~120"
                @change="drawPressGraph"
              />
              <span v-html="'&nbsp; &nbsp;分钟'"></span>
            </el-form-item>
            <div v-show="plan_info.press_info.press_type === 1">
              <el-form-item label="起始并发数">
                <el-input
                  v-model.number="plan_info.press_info.start_concurrency"
                  clearable
                  class="input-width"
                  size="small"
                  @change="drawPressGraph"
                />
              </el-form-item>
              <el-form-item label="步 长">
                <el-input
                  v-model.number="plan_info.press_info.step_size"
                  clearable
                  class="input-width"
                  size="small"
                  @change="drawPressGraph"
                />
              </el-form-item>
              <el-form-item label="每步时长">
                <el-input
                  v-model.number="plan_info.press_info.step_duration"
                  clearable
                  class="input-width"
                  size="small"
                  @change="drawPressGraph"
                />
                <span v-html="'&nbsp; &nbsp;秒'"></span>
              </el-form-item>
              <el-form-item label="最大并发数">
                <el-input
                  v-model.number="plan_info.press_info.concurrency"
                  clearable
                  class="input-width"
                  size="small"
                  @change="drawPressGraph"
                />
              </el-form-item>
              <el-form-item label="持续时间">
                <el-input
                  v-model.number="plan_info.press_info.duration"
                  clearable
                  class="input-width"
                  size="small"
                  placeholder="输入范围： 1~120"
                  @change="drawPressGraph"
                />
                <span v-html="'&nbsp; &nbsp;分钟'"></span>
              </el-form-item>
            </div>
            <el-form-item label="机器数量">
              <span v-if="machine_list_len === 0">当前无可用压测机</span>
              <el-input-number
                v-else
                v-model="plan_info.engine_count"
                class="input-width"
                controls-position="right"
                :min="1"
                :max="machine_list_len"
                @change="drawPressGraph"
              />
            </el-form-item>
            <el-form-item label="RPS总数" v-if="plan_info.press_info.press_type == 3">
              <el-input
                v-model.number="plan_info.press_info.rps"
                clearable
                class="input-width"
                @change="drawPressGraph"
              />
            </el-form-item>
            <el-form-item
              label="压测场景"
              v-show="[0, 1, 3, 4].indexOf(plan_info.press_info.press_type) > -1"
            >
              <div style="display: flex; flex-direction: column">
                <div v-show="[1, 3].indexOf(plan_info.press_info.press_type) > -1">
                  流量比例和: {{ getSumRate() }}%
                </div>
                <div
                  v-for="item in getPressSceneList()"
                  :key="item.scene_id"
                  style="display: inline-flex; margin-top: 5px"
                >
                  <div>
                    <el-input v-model="item.scene_name" :disabled="true" placeholder="场景名称" />
                  </div>
                  <div style="display: inline-flex">
                    <div
                      v-if="[1, 3].indexOf(plan_info.press_info.press_type) > -1"
                      style="display: flex"
                    >
                      <span class="concurrency-span" style="width: 60px"> 流量比例 </span>
                      <el-input
                        v-model.number="item.rate"
                        placeholder="比例"
                        style="width: 150px"
                        @change="drawPressGraph"
                      >
                        <template #suffix> % </template>
                      </el-input>
                    </div>
                    <div
                      v-if="[0, 4].indexOf(plan_info.press_info.press_type) > -1"
                      style="display: flex"
                    >
                      <span class="concurrency-span" style="width: 46px"> 并发数 </span>
                      <el-input
                        v-model.number="item.concurrency"
                        placeholder="并发数"
                        style="width: 150px"
                        @change="drawPressGraph"
                      />
                    </div>
                    <div
                      v-if="[4].indexOf(plan_info.press_info.press_type) > -1"
                      style="display: flex"
                    >
                      <span class="concurrency-span" style="width: 30px"> 次数 </span>
                      <el-input
                        v-model.number="item.iteration"
                        placeholder="压测次数"
                        style="width: 150px"
                        @change="drawPressGraph"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </el-form-item>
          </el-form>
          <echart
            v-if="[0, 1].indexOf(plan_info.press_info.press_type) > -1"
            ref="chartsRef"
            :options="optionsData"
            :height="400"
            :width="560"
            style="margin-left: 20px"
          />
          <div style="flex: 1; min-height: 0" v-if="plan_info.press_info.press_type === 2">
            <div>
              <TargetRpsTree
                ref="targetRpsTree"
                style="padding-bottom: 10px"
                :sceneList="scene_list"
                v-model:config="plan_info.press_info.scene_list"
              />
            </div>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="高级配置" name="other_config">
        <el-form
          label-position="right"
          label-width="100px"
          :model="plan_info"
          style="width: 1000px; margin-top: 30px"
        >
          <el-form-item label="定时任务">
            <el-row style="width: 500px">
              <el-col :span="3">
                <el-switch
                  v-model="plan_info.task_info.enable"
                  style="--el-switch-on-color: #13ce66"
                />
              </el-col>
              <el-col :span="21">
                <el-input
                  style="width: 266px"
                  v-show="plan_info.task_info.enable"
                  v-model="plan_info.task_info.cron"
                  placeholder="cron表达式"
                >
                  <template #append>
                    <el-button :icon="timeIcon" @click="chooseCron" />
                  </template>
                </el-input>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item label="容量采样">
            <el-switch v-model="plan_info.capacity_switch" style="--el-switch-on-color: #13ce66" />
          </el-form-item>
          <el-form-item label="采样策略">
            <el-row style="width: 500px">
              <el-col :span="8">
                <el-select v-model="plan_info.sampling_info.sampling_type" style="width: 140px">
                  <el-option
                    v-for="item in sampling_type_list"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-col>
              <el-col :span="3" v-show="plan_info.sampling_info.sampling_type > 0">
                <div>万分之</div>
              </el-col>
              <el-col :span="11">
                <el-input
                  v-model.number="plan_info.sampling_info.sampling_rate"
                  v-show="plan_info.sampling_info.sampling_type > 0"
                  style="width: 100px"
                  placeholder="输入范围： 1~100"
                  :onchange="samplingRateCheck"
                />
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item label="管理服务">
            <el-cascader
              v-model="plan_info.server_info"
              :show-all-levels="false"
              collapse-tags
              collapse-tags-tooltip
              :max-collapse-tags="3"
              :options="namespace_list"
              :props="cas_props"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <!--      <el-tab-pane label="熔断策略" name="press_break_config" />-->
      <el-tab-pane label="压测记录" name="press_history">
        <div style="margin: 10px 0 80px 0">
          <report-search v-if="active_tab === 'press_history'" v-model="report_search" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="调试记录" name="debug_record">
        <div style="margin: 10px 0 80px 0">
          <debug-result-list v-if="active_tab === 'debug_record'" :plan-id="plan_id" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="操作日志" name="operation_log">
        <operation-log
          v-if="active_tab === 'operation_log'"
          :source-id="plan_id"
          :source-name="SourceName.Plan"
        />
      </el-tab-pane>
    </el-tabs>
  </el-skeleton>
  <debug-result ref="debugResultRef" />
  <el-dialog v-model="cron_dialog_visible">
    <vue3-cron-plus-picker
      @hide="cron_dialog_visible = false"
      @fill="fillValue"
      :expression="expression"
    />
  </el-dialog>
</template>

<style scoped lang="less">
.plan-header {
  height: 50px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: -10px;
  background-color: #fff;
  box-shadow: 0 2px 15px 0 rgba(0, 0, 0, 0.1);
  border-top: 1px solid #ebeef5;

  .title-edit-input {
    width: 300px;
    cursor: text;
  }

  .title-div {
    display: flex;
    align-items: center;
    gap: 5px;

    .title-span {
      margin-left: 10px;
      font-size: 20px;
      font-weight: 500;
      color: #409eff;
    }
    .icon {
      display: none;
    }
  }

  .title-div:hover .icon {
    display: block;
    color: #409eff;
  }
}

.plan-tabs {
  height: calc(100vh - 160px);
  background-color: #fff;

  box-shadow: 0 2px 15px 0 rgba(0, 0, 0, 0.1);
  border-top: 1px solid #ebeef5;

  :deep(.el-tabs__content) {
    padding-right: 12px;
    height: calc(100vh - 160px);
    overflow-y: auto;
  }

  .remark {
    max-width: 90%;
    overflow-wrap: break-word;
    word-break: break-all;
  }

  .press-cfg-div {
    margin-top: 20px;
    display: flex;

    .input-width {
      width: 220px;
    }

    .concurrency-span {
      margin: 0 10px 0 10px;
      color: var(--el-text-color-regular);
    }

    .base-info-div {
      width: 100%;
      display: flex;

      .scene-edit {
        margin-left: 30px;
        width: 800px;

        .scene-row {
          height: 20px;
          width: 100%;
          margin-bottom: 16px;
          display: flex;

          .scene-col {
            width: 160px;
            margin-right: 12px;
          }
        }
      }
    }

    .machine-info-div {
      display: flex;

      .machine-info-card {
        margin-left: 30px;
      }
    }

    .press-div {
      display: flex;
    }

    .break-div {
      display: flex;
      //flex-direction: column;

      .break-info-card {
        width: 300px;
        height: 200px;
      }
    }
  }

  .data-cfg-container {
    height: calc(100% - 30px);
  }
}
</style>

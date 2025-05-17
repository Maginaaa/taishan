<template>
  <div>
    <span style="font-size: 15px; font-weight: 400"> RPS场景配置, 总目标RPS ( {{ sumRps }} )</span>
    <el-card
      shadow="hover"
      v-for="(scene, index) in scene_list"
      :key="scene.scene_id"
      style="margin-top: 10px"
      :body-style="{ padding: '10px' }"
    >
      <div style="font-size: 14px; font-weight: 500; margin-bottom: 2px">
        {{ scene.scene_name }}
      </div>
      <el-tree
        :ref="(el) => (treeRefs[index] = el)"
        v-if="scene.case_tree"
        :indent="15"
        node-key="case_id"
        default-expand-all
        :data="scene.case_tree || []"
      >
        <template #default="{ data }">
          <div style="flex: auto">
            <el-card
              shadow="hover"
              :body-style="{ padding: '5px 10px' }"
              style="height: 42px; line-height: 30px"
            >
              <span v-if="data.type === 1">
                <el-tag size="small">HTTP请求</el-tag> {{ data.title }}
                <span style="float: right">
                  <span style="margin-right: 5px"> 预期RPS: </span>
                  <el-input-number
                    style="width: 90px"
                    v-model="data.target_rps"
                    :min="1"
                    :max="100000"
                    :controls="false"
                    @change="channelTargetRps"
                  />
                </span>
              </span>
              <span v-if="data.type === 11">
                <el-tag size="small" type="warning">逻辑控制器</el-tag> {{ data.title }}
              </span>
            </el-card>
          </div>
        </template>
      </el-tree>
    </el-card>
  </div>
</template>

<script setup lang="ts">
// import { computed } from 'vue'
import { ref, watch } from 'vue'
const props = defineProps({
  sceneList: {
    type: Array,
    defalt: []
  },
  config: {
    type: Array,
    defalt: []
  }
})

const treeRefs = ref<any>([])
const sceneConfigList = ref<any>([])
const sumRps = ref(0)

const scene_list = ref<any>([])
const emits = defineEmits(['update:config'])
// const sumRps = computed(() => {
//   if (!sceneConfigList.value) return 0

//   let rpsCount = 0
//   for (let index = 0; index < sceneConfigList.value.length; index++) {
//     const sceneConfig = sceneConfigList.value[index]
//     if (sceneConfig && sceneConfig.case_tree) {
//       for (let index1 = 0; index1 < sceneConfig.case_tree.length; index1++) {
//         const treeCase = sceneConfig.case_tree[index1]
//         rpsCount = rpsCount + treeCase.target_rps || 1
//       }
//     }
//   }

//   return rpsCount
// })

watch(
  [() => props.config, () => props.sceneList],
  () => {
    const tree = JSON.parse(JSON.stringify(props.sceneList || [])).filter((item) => !item.disabled)
    const propsConfig = props.config || []
    sumRps.value = 0
    sceneConfigList.value = propsConfig
    for (let index = 0; index < tree.length; index++) {
      const scene = tree[index]
      const configScene: any = propsConfig.find((item: any) => item.scene_id === scene.scene_id)
      if (scene && scene.case_tree) {
        buildCaseTree(scene.case_tree, configScene?.case_tree || [])
      }
    }
    scene_list.value = tree
  },
  {
    immediate: true,
    deep: true
  }
)

function buildCaseTree(caseTree, caseConfigs) {
  let delIndexs: any = []
  for (let index1 = 0; index1 < caseTree.length; index1++) {
    const case1 = caseTree[index1]
    if (case1.disabled) {
      if (caseConfigs) {
        caseConfigs = caseConfigs.filter((item) => item.case_id !== case1.case_id)
      }
      delIndexs.push(index1)
      continue
    }

    if (case1.type === 1) {
      if (caseConfigs) {
        const caseConfig = caseConfigs.find((item) => item.case_id === case1.case_id)
        case1.target_rps = caseConfig?.target_rps || 1
      } else {
        case1.target_rps = 1
      }
      sumRps.value = sumRps.value + case1.target_rps
    }

    if (case1.children) {
      buildCaseTree(case1.children, caseConfigs)
    }
  }
  let delCount = 0
  for (let item of delIndexs) {
    caseTree.splice(item - delCount, 1)
    delCount = delCount + 1
  }
}

function channelTargetRps() {
  const sceneConfigs: any = []
  const tree = scene_list.value || []
  for (let index = 0; index < tree.length; index++) {
    const scene = tree[index]
    const sceneConfig = {
      scene_id: scene.scene_id,
      scene_name: scene.scene_name,
      scene_type: scene.scene_type,
      concurrency: scene.concurrency || 1,
      rate: scene.rate || 1,
      case_tree: []
    }
    sceneConfigs.push(sceneConfig)
    buildCaseConfig(sceneConfig.case_tree, scene.case_tree)
  }
  emits('update:config', sceneConfigs)
}

function buildCaseConfig(caseConfigList, caseTree) {
  if (!caseTree) {
    return
  }
  for (let index = 0; index < caseTree.length; index++) {
    const element = caseTree[index]
    if (element.type === 1) {
      //只支持http接口设置RPS
      caseConfigList.push({
        case_id: element.case_id,
        target_rps: element.target_rps || 1
      })
    }
    if (element.children) {
      buildCaseConfig(caseConfigList, element.children)
    }
  }
}

defineExpose({
  channelTargetRps
})
</script>

<style scoped>
:deep(.el-tree-node__content) {
  height: 48px;
  /* width: 100%; */
  flex: 1;
}
</style>

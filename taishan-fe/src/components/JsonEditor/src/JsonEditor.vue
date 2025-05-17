<script setup lang="ts">
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { propTypes } from '@/utils/propTypes'
import { computed, ref, unref } from 'vue'
import useClipboard from 'vue-clipboard3'
import { isEmptyVal } from '@/utils/is'
import { ClickOutside as vClickOutside } from 'element-plus'

const emits = defineEmits([
  'update:modelValue',
  'node-click',
  'brackets-click',
  'icon-click',
  'selected-value',
  'value-extract'
])

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  },
  readOnly: {
    type: Boolean,
    default: true
  },
  deep: propTypes.number.def(5),
  showLength: propTypes.bool.def(true),
  showLineNumbers: propTypes.bool.def(true),
  showLineNumber: propTypes.bool.def(true),
  showIcon: propTypes.bool.def(true),
  showDoubleQuotes: propTypes.bool.def(true),
  virtual: propTypes.bool.def(false),
  height: propTypes.number.def(400),
  itemHeight: propTypes.number.def(20),
  rootPath: propTypes.string.def('$'),
  nodeSelectable: propTypes.func.def(),
  selectableType: propTypes.oneOf<'multiple' | 'single'>(['multiple', 'single']).def(),
  showSelectController: propTypes.bool.def(false),
  selectOnClickNode: propTypes.bool.def(true),
  highlightSelectedNode: propTypes.bool.def(true),
  collapsedOnClickBrackets: propTypes.bool.def(true),
  renderNodeKey: propTypes.func.def(),
  renderNodeValue: propTypes.func.def(),
  editable: propTypes.bool.def(false),
  editableTrigger: propTypes.oneOf<'click' | 'dblclick'>(['click', 'dblclick']).def('click')
})

const data = computed(() => props.modelValue)

const localModelValue = computed({
  get: () => data.value,
  set: (val) => {
    emits('update:modelValue', val)
  }
})
const nodeClick = (node: any) => {
  emits('node-click', node)
}

const bracketsClick = (collapsed: boolean) => {
  emits('brackets-click', collapsed)
}

const iconClick = (collapsed: boolean) => {
  emits('icon-click', collapsed)
}

const selectedChange = (newVal: any, oldVal: any) => {
  emits('selected-value', newVal, oldVal)
}

const { toClipboard } = useClipboard()
const valueExtract = (type) => {
  popoverVisible.value = false
  if (type === 1 || type === 2) {
    emits('value-extract', type, current_data.value)
  } else {
    toClipboard(JSON.stringify(data.value))
  }
}
const getValueStyle = (val) => {
  const style = 'font-weight: bold;'
  switch (typeof val) {
    case 'number':
      return style + 'color: #25aae2;'
    case 'string':
      if (val === 'null') {
        return style + 'color: #f98280;'
      }
      if (val.startsWith('"http://') || val.startsWith('"https://')) {
        return style + 'color: #61D2D6;text-decoration: underline;'
      }
      return style + 'color: #3ab54a;'
    case 'boolean':
      return style + 'color: #f98280;'
  }
}

const click = (event) => {
  const target = event.target
  if (isEmptyVal(target.id)) {
    popoverVisible.value = false
    unref(popoverRef).popperRef?.delayHide?.()
  }
}

const popoverRef = ref()
const popoverVisible = ref(false)

let nodeRef = ref<any>('')
let current_data = ref()
const valClick = (node) => {
  popoverVisible.value = true
  current_data.value = node
  nodeRef.value = document.getElementById(node.path)
}

const onClickOutside = () => {
  popoverVisible.value = false
}
</script>

<template>
  <vue-json-pretty
    v-model:data="localModelValue"
    :deep="deep"
    :show-length="showLength"
    :show-line-numbers="showLineNumbers"
    :show-line-number="showLineNumber"
    :show-icon="showIcon"
    :show-double-quotes="showDoubleQuotes"
    :virtual="virtual"
    :height="height"
    :item-height="itemHeight"
    :root-path="rootPath"
    :node-selectable="nodeSelectable"
    :selectable-type="selectableType"
    :show-select-controller="showSelectController"
    :select-on-click-node="selectOnClickNode"
    :highlight-selected-node="highlightSelectedNode"
    :collapsed-on-click-brackets="collapsedOnClickBrackets"
    :render-node-key="renderNodeKey"
    :render-node-value="renderNodeValue"
    :editable="editable"
    :editable-trigger="editableTrigger"
    @click="click"
    @node-click="nodeClick"
    @brackets-click="bracketsClick"
    @icon-click="iconClick"
    @selected-change="selectedChange"
  >
    <template #renderNodeKey="{ defaultKey }">
      <span style="color: #92278f; font-weight: bold">
        {{ defaultKey }}
      </span>
    </template>
    <template #renderNodeValue="{ node, defaultValue }">
      <span
        :id="node.path"
        :style="getValueStyle(defaultValue)"
        @click="valClick(node)"
        v-click-outside="onClickOutside"
      >
        {{ defaultValue }}
      </span>
    </template>
  </vue-json-pretty>
  <el-popover
    :visible="popoverVisible"
    :virtual-ref="nodeRef"
    ref="popoverRef"
    trigger="click"
    virtual-triggering
    popper-class="json-popover"
  >
    <el-menu @select="valueExtract">
      <el-menu-item v-show="!props.readOnly" :index="1"> 添加断言 </el-menu-item>
      <el-menu-item v-show="!props.readOnly" :index="2"> 变量提取 </el-menu-item>
      <el-menu-item :index="3"> 全文复制 </el-menu-item>
    </el-menu>
  </el-popover>
</template>

<style lang="less">
.json-popover {
  --el-popover-padding: 0;
}
</style>

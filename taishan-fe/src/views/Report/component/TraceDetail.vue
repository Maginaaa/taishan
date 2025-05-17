<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
  traceDetail: {
    type: Array
  }
})

let basePxWidth = ref(0.0)
let baseEndTimestamp = ref(0)
let totalPxWidth = 0.0
const calculateBasePxWidth = () => {
  const _trace_dt = props.traceDetail[0]
  const pageWidth = window.innerWidth * 0.6 // 0.6是因为dialog宽度设置为60%
  totalPxWidth = pageWidth - (18 * getMaxDepth(_trace_dt) + 350 + 100) // 18是一个箭头的宽度，350是文本最大宽度，100是右侧文字宽度
  baseEndTimestamp.value = _trace_dt.timestamp + _trace_dt.duration
  basePxWidth.value = totalPxWidth / _trace_dt.duration
}
const getMaxDepth = (node) => {
  if (!node) {
    return 0
  }
  if (node.children && node.children.length === 0) {
    return 1
  }
  let maxDepth = 0
  for (let child of node.children) {
    let childDepth = getMaxDepth(child)
    if (childDepth > maxDepth) {
      maxDepth = childDepth
    }
  }
  return maxDepth + 1
}

calculateBasePxWidth()
</script>

<template>
  <el-tree :data="traceDetail" default-expand-all node-key="id">
    <template #default="{ data }">
      <div class="node-div">
        <div>
          <div class="content ope-name">{{ data.operation_name }}</div>
          <div class="content gray">{{ data.app }} {{ data.ip }}</div>
        </div>
        <div
          :style="{
            borderBottom: '1px dashed rgb(221.7, 222.6, 224.4)',
            height: '1px',
            flex: '1 1 auto'
          }"
        >
        </div>
        <div
          :style="{
            backgroundColor: '#67C23A',
            width: `${basePxWidth * data.duration}px`,
            height: '10px'
          }"
        >
        </div>
        <div
          :style="{
            borderBottom: '1px dashed rgb(221.7, 222.6, 224.4)',
            height: '1px',
            width: `${basePxWidth * (baseEndTimestamp - data.timestamp - data.duration)}px`
          }"
        >
        </div>
        <div class="time gray">{{ data.duration }}ms</div>
      </div>
    </template>
  </el-tree>
  <el-table
    :data="traceDetail"
    border
    stripe
    style="width: 100%"
    max-height="600px"
    default-expand-all
    row-key="id"
    :tree-props="{ children: 'children' }"
  >
    <el-table-column header-align="center" prop="operation_name" label="接口名" />
    <el-table-column align="center" header-align="center" prop="app" label="服务名" />
    <el-table-column align="center" header-align="center" prop="ip" label="所属IP" />
    <el-table-column align="center" header-align="center" prop="duration" label="时间(ms)" />
  </el-table>
</template>

<style scoped lang="less">
:deep(.el-tree-node__content) {
  height: auto;
}

.node-div {
  height: 50px;
  width: 100%;
  display: flex;
  align-items: center;

  .ope-name {
    color: black;
  }

  .gray {
    color: gray;
  }

  .content {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 320px;
  }

  .time {
    width: 60px;
    margin-left: 10px;
  }
}
</style>

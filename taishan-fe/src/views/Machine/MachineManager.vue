<script setup lang="ts">
import { ref } from 'vue'
import CpuChart from '@/views/Machine/chart/CpuChart.vue'
import MemChart from '@/views/Machine/chart/MemChart.vue'
import MachineService from '@/api/scene/machine'

const machineService = new MachineService()

let auto_refresh = ref(false)
let graph_data = ref({
  cpu: [],
  mem: [],
  network: []
})

const getMachineUseInfo = async () => {
  graph_data.value = await machineService.getAllMachineUseInfo()
}
let data_interval: any = null

const startGetMachineData = () => {
  if (!data_interval) {
    data_interval = setInterval(getMachineUseInfo, 5000)
  }
}

const autoRefreshChange = () => {
  if (auto_refresh.value) {
    startGetMachineData()
  } else {
    clearInterval(data_interval)
    data_interval = null
  }
}

getMachineUseInfo()
</script>

<template>
  <el-card class="operation-row">
    <div style="float: right">
      5s 自动刷新
      <el-switch
        v-model="auto_refresh"
        @change="autoRefreshChange"
        style="--el-switch-on-color: #13ce66"
      />
    </div>
  </el-card>
  <el-card class="graph-row">
    <cpu-chart :data="graph_data.cpu" />
    <mem-chart :data="graph_data.mem" />
  </el-card>
</template>

<style scoped lang="less">
.operation-row {
  width: 100%;
  display: flex;
  flex-direction: row;
}

.graph-row {
  width: 100%;
  height: calc(100% - 80px);
  display: flex;
  flex-direction: row;
  margin-top: 10px;

  :deep(.el-card__body) {
    width: 100%;
  }

  .graph-item {
    width: calc((100% - 20px) / 3);
    height: 280px;
  }
}
</style>

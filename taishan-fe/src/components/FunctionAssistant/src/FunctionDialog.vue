<script setup lang="ts">
import { reactive, ref, unref } from 'vue'
import useClipboard from 'vue-clipboard3'
import { ElMessage } from 'element-plus'
import { useIcon } from '@/hooks/web/useIcon'
import UtilsService from '@/api/scene/utils'

const visible = ref(false)
let current_function = ref<string>('')
let current_function_template = ref<string>('')
let current_function_params = ref([])

const utilsService = new UtilsService()

const copyIcon = useIcon({ icon: 'ep:copy-document' })

interface formType {
  expression: string
  result: string
}

let expressionForm = reactive<formType>({
  expression: '',
  result: ''
})

const openDialog = () => {
  visible.value = true
}

const chooseParam = (val) => {
  if (val === '') {
    current_function_template.value = ''
    current_function_params.value = []
    return
  }
  const obj = function_params_mp[val]
  current_function_template.value = obj.function_temp
  current_function_params.value = [...obj.params]
  updateExpression()
}

const { toClipboard } = useClipboard()
const copyExpr = () => {
  if (expressionForm.expression === '') {
    ElMessage.warning('表达式为空')
    return
  }
  toClipboard(`\${${expressionForm.expression}\}`)
  ElMessage.success('复制成功')
}

const debug = async () => {
  const res = await utilsService.Debug(expressionForm.expression)
  expressionForm.result = res
}

const updateExpression = () => {
  let _expression = current_function_template
  current_function_params.value.forEach((item) => {
    _expression.value = _expression.value.replace(item.param_name, item.param_val)
  })
  expressionForm.expression = unref(_expression)
}

const function_list = [
  {
    label: '随机函数',
    options: [
      {
        function_name: '随机数字',
        function_key: '__RandomInt',
        function_temp: '__RandomInt(param1, param2, param3)'
      },
      {
        function_name: '随机字符串',
        function_key: '__RandomStr',
        function_temp: '__RandomStr(param1, param2, param3)'
      },
      {
        function_name: '随机选择字符串',
        function_key: '__RandomChooseStr',
        function_temp: '__RandomChooseStr(param1, param2)'
      },
      {
        function_name: '身份证号生成',
        function_key: '__IDCard',
        function_temp: '__IDCard(param1)'
      },
      {
        function_name: '中文名生成',
        function_key: '__ChineseName',
        function_temp: '__ChineseName(param1)'
      },
      {
        function_name: '手机号生成',
        function_key: '__Mobile',
        function_temp: '__Mobile(param1)'
      },
      {
        function_name: 'UUID生成',
        function_key: '__UUID',
        function_temp: '__UUID(param1)'
      }
    ]
  },
  {
    label: '字符串处理',
    options: [
      {
        function_name: '大小写转换',
        function_key: '__ChangeCase',
        function_temp: '__ChangeCase(param1, param2, param3)'
      },
      {
        function_name: '时间格式化',
        function_key: '__TimeFormat',
        function_temp: '__TimeFormat(param1, param2)'
      },
      {
        function_name: '字符串截取',
        function_key: '__SubStr',
        function_temp: '__SubStr(param1, param2, param3, param4)'
      },
      {
        function_name: '四则运算',
        function_key: '__Calculate',
        function_temp: '__Calculate(param1, param2, param3, param4)'
      }
    ]
  },
  {
    label: '加解密',
    options: [
      {
        function_name: 'MD5加密',
        function_key: '__MD5',
        function_temp: '__MD5(param1, param2)'
      },
      {
        function_name: 'Base64加密',
        function_key: '__Base64Encode',
        function_temp: '__Base64Encode(param1, param2)'
      },
      {
        function_name: 'Base64解密',
        function_key: '__Base64Decode',
        function_temp: '__Base64Decode(param1, param2)'
      },
      {
        function_name: 'RSA加密',
        function_key: '__RSA',
        function_temp: '__RSA(param1, param2, param3)'
      }
    ]
  }
]
const function_params_mp = {
  __RandomInt: {
    function_temp: '__RandomInt(param1, param2, param3)',
    params: [
      {
        param_name: 'param1',
        param_val: '1',
        desc: '随机范围的左(开)区间',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '10',
        required: false,
        desc: '随机范围的右(闭)区间',
        val_type: 'input'
      },
      {
        param_name: 'param3',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __RandomStr: {
    function_temp: '__RandomStr(param1, param2, param3)',
    params: [
      {
        param_name: 'param1',
        param_val: '1',
        required: false,
        desc: '随机字符串长度，默认为1',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '随机字符的取值范围，默认为所有大小写字母和数字',
        val_type: 'input'
      },
      {
        param_name: 'param3',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __RandomChooseStr: {
    function_temp: '__RandomChooseStr(param1, param2)',
    params: [
      {
        param_name: 'param1',
        param_val: 'a|b|c',
        required: true,
        desc: '随机字符串选项列表，用|做分割',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __IDCard: {
    function_temp: '__IDCard(param1)',
    params: [
      {
        param_name: 'param1',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __ChineseName: {
    function_temp: '__ChineseName(param1)',
    params: [
      {
        param_name: 'param1',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __Mobile: {
    function_temp: '__Mobile(param1)',
    params: [
      {
        param_name: 'param1',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __UUID: {
    function_temp: '__UUID(param1)',
    params: [
      {
        param_name: 'param1',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __ChangeCase: {
    function_temp: '__ChangeCase(param1, param2, param3)',
    params: [
      {
        param_name: 'param1',
        param_val: 'abc',
        required: true,
        desc: '需要改变大小的字符串',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: 'UPPER',
        required: true,
        desc: 'UPPER:转大写、 LOWER:转小写',
        val_type: 'select',
        select_options: []
      },
      {
        param_name: 'param3',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __TimeFormat: {
    function_temp: '__TimeFormat(param1, param2)',
    params: [
      {
        param_name: 'param1',
        param_val: 'YYYY-MM-DD HH:MM:SS',
        required: false,
        desc: '时间格式，不传即为获取时间戳',
        val_type: 'select'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __SubStr: {
    function_temp: '__SubStr(param1, param2, param3, param4)',
    params: [
      {
        param_name: 'param1',
        param_val: 'abcde',
        required: true,
        desc: '需要截取的字符串',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '1',
        required: false,
        desc: '需要截取内容的左(开)索引',
        val_type: 'select'
      },
      {
        param_name: 'param3',
        param_val: '3',
        required: false,
        desc: '需要截取内容的右(闭)索引'
      },
      {
        param_name: 'param4',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __Calculate: {
    function_temp: '__Calculate(param1, param2, param3, param4)',
    params: [
      {
        param_name: 'param1',
        param_val: '+',
        required: true,
        desc: '操作符，仅支持+、-、*、/',
        val_type: 'select'
      },
      {
        param_name: 'param2',
        param_val: '1',
        required: true,
        desc: '需计算的左内容(仅支持整型)',
        val_type: 'select'
      },
      {
        param_name: 'param3',
        param_val: '2',
        required: true,
        desc: '需计算的右内容(仅支持整型)'
      },
      {
        param_name: 'param4',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __MD5: {
    function_temp: '__MD5(param1, param2)',
    params: [
      {
        param_name: 'param1',
        param_val: 'abc',
        required: true,
        sup_param: true,
        desc: '需加密的参数',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __Base64Encode: {
    function_temp: '__Base64Encode(param1, param2)',
    params: [
      {
        param_name: 'param1',
        param_val: 'abc',
        required: true,
        desc: '需加密的参数',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __Base64Decode: {
    function_temp: '__Base64Decode(param1, param2)',
    params: [
      {
        param_name: 'param1',
        param_val: 'YWJj',
        required: true,
        desc: '需解密的参数',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '',
        required: false,
        desc: '另存为新变量'
      }
    ]
  },
  __RSA: {
    function_temp: '__RSA(param1, param2, param3)',
    params: [
      {
        param_name: 'param1',
        param_val: 'pubKey',
        required: true,
        sup_param: true,
        desc: '公钥',
        val_type: 'input'
      },
      {
        param_name: 'param2',
        param_val: '123456',
        required: true,
        sup_param: true,
        desc: '加密变量'
      },
      {
        param_name: 'param3',
        param_val: '',
        desc: '另存为新变量'
      }
    ]
  }
}

defineExpose({
  openDialog
})
</script>

<template>
  <div class="container-func">
    <el-dialog title="函数助手" v-model="visible" width="700px">
      <div class="func-container">
        <div class="row">
          <el-select
            v-model="current_function"
            class="func-select"
            placeholder="请选择函数"
            filterable
            clearable
            @change="chooseParam"
          >
            <el-option-group v-for="group in function_list" :key="group.label" :label="group.label">
              <el-option
                v-for="item in group.options"
                :key="item.function_key"
                :label="item.function_name"
                :value="item.function_key"
              >
                <span style="float: left">{{ item.function_name }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">
                  {{ item.function_temp }}
                </span>
              </el-option>
            </el-option-group>
          </el-select>
          <div class="func-temp">
            {{ current_function_template }}
          </div>
        </div>
        <el-table :data="current_function_params" border class="param-table">
          <el-table-column
            align="center"
            header-align="center"
            prop="param_name"
            label="参数名"
            width="150"
          />
          <el-table-column
            align="center"
            header-align="center"
            prop="param_val"
            label="参数值"
            width="150"
          />
          <el-table-column
            align="center"
            header-align="center"
            prop="required"
            label="必填"
            width="60"
          >
            <template #default="scope">
              <el-tag v-if="scope.row.required" type="success" size="small">是</el-tag>
              <el-tag v-else type="danger" size="small">否</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="center" header-align="center" label="支持参数" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.sup_param" type="success" size="small">是</el-tag>
              <el-tag v-else type="danger" size="small">否</el-tag>
            </template>
          </el-table-column>
          <el-table-column header-align="center" prop="desc" label="描述" />
        </el-table>
        <el-form
          ref="form"
          class="res-form"
          :model="expressionForm"
          label-position="left"
          label-width="70px"
        >
          <el-form-item label="表达式">
            <el-input v-model="expressionForm.expression" class="exp-input">
              <template #append>
                <el-button type="primary" :icon="copyIcon" @click="copyExpr" />
              </template>
            </el-input>
            <el-button
              type="primary"
              style="margin-left: 15px"
              v-tc="{ name: '函数助手', param: expressionForm.expression }"
              @click="debug"
            >
              生成
            </el-button>
          </el-form-item>
          <el-form-item label="函数结果" style="margin-top: 10px">
            <el-input v-model="expressionForm.result" />
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped lang="less">
.container-func {
  :deep(.el-dialog) {
    padding: 16px;
  }

  :deep(.el-dialog__title) {
    font-size: 20px;
    font-weight: bold;
  }

  :deep(.el-form-item) {
    margin-bottom: 0;
  }

  .func-container {
    overflow-y: auto;

    .row {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .func-select {
        width: 360px;
      }

      .func-temp {
        float: right;
        font-weight: bold;
        color: #409eff;
      }
    }

    .param-table {
      margin-top: 15px;
      width: 100%;
    }

    .res-form {
      margin-top: 10px;

      .exp-input {
        margin-top: 5px;
        width: 518px;
      }
    }
  }
}
</style>

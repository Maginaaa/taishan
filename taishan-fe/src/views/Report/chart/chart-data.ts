import { EChartsOption } from 'echarts/types/dist/echarts'

export const baseChartOption: EChartsOption = {
  title: {
    left: 'center',
    textStyle: {
      // 文字颜色
      // color: '#6e6e6e',
      // // 字体风格,'normal','italic','oblique'
      // fontStyle: 'normal',
      // // 字体粗细 'normal','bold','bolder','lighter',100 | 200 | 300 | 400...
      // fontWeight: 'bold',
      // // 字体系列
      // fontFamily: 'sans-serif',
      // // 字体大小
      fontSize: 14
    }
  },
  grid: {
    left: '45px',
    right: '30px',
    bottom: '40px'
  },
  dataset: {
    source: []
  },
  tooltip: {
    trigger: 'axis'
  },
  xAxis: {
    type: 'category'
  },
  yAxis: {
    type: 'value',
    splitNumber: 4,
    splitLine: {
      show: true, // 显示分割线
      lineStyle: {
        color: ['rgba(0,0,0,0.1)'], // 分割线的颜色，可以设置为数组来显示不同颜色的虚线段
        width: 1, // 分割线的宽度
        type: 'dashed' // 分割线类型，'dashed' 表示虚线
      }
    }
  },
  series: [
    {
      type: 'line',
      showSymbol: false
    }
  ]
}

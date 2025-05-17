/**
 * 指令 v-tc
 * @Example v-tc="{name:'保存接口场景用例',params2:'XX'}"
 */
import type { App, Directive, DirectiveBinding } from 'vue'
import { trackClick } from '@/point/utils.js'

const mounted = (el: Element, binding: DirectiveBinding<any>) => {
  //实现
  el.addEventListener('click', () => {
    const data = binding.value || {} //接收传参
    trackClick(data)
  })
}

const point: Directive = {
  mounted
}

export function setupTrackClickDirective(app: App) {
  app.directive('tc', point)
}

export default point

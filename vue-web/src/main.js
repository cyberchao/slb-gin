import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
ElementUI.Dialog.props.closeOnClickModal.default = false
Vue.use(ElementUI);

import { store } from '@/store/index'
// 引入封装的router
import router from '@/router/index'

Vue.config.productionTip = false
export default new Vue({
  render: h => h(App),
  router,
  store
}).$mount('#app')

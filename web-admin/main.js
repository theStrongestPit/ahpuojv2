import Vue from 'vue';
import App from '@/web-admin/js/App.vue';
import router from '@/web-admin/js/router';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import SvgIcon from '@/web-common/components/svgicon.vue';
import store from '@/web-admin/js/store';
import Base from '@/web-common/base';

Vue.use(ElementUI);
Vue.use(Base); // 注册的全局函数

// Vue.config.productionTip = false
Vue.config.devtools = true;

// svg图标
Vue.component('svg-icon', SvgIcon);
const requireAll = requireContext => requireContext.keys().map(requireContext);
const req = require.context('@/web-common/assets/icons', false, /\.svg$/);
requireAll(req);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');

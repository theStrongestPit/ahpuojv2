import Vue from 'vue';
import Vuex from 'vuex';
import user from '@/web-user/js/store/modules/user';
import bus from '@/web-user/js/store/modules/bus';
import getters from '@/web-user/js/store/getters';
Vue.use(Vuex);
export default new Vuex.Store({
  modules: {
    user,
    bus
  },
  getters
});

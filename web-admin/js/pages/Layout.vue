<template lang="pug">
  div
    top-bar(:screenWidth="screenWidth")
    transition(name="slide-fade")
      sider-bar(v-if="screenWidth>=960")
    transition(name="el-fade-in-linear")
      router-view(v-if="$route.meta.keepAlive === false")
    transition(name="el-fade-in-linear")
      keep-alive
        router-view(v-if="$route.meta.keepAlive !== false")
</template>

<script>
import SiderBar from '@/web-admin/js/common/siderbar.vue';
import TopBar from '@/web-admin/js/common/topbar.vue';
import '@/web-admin/sass/main.scss';

export default {
  name: 'Layout',
  components: {
    SiderBar,
    TopBar
  },
  data() {
    return {
      screenWidth: document.body.clientWidth
    };
  },

  watch: {
    screenWidth(val) {
      // 为了避免频繁触发resize函数导致页面卡顿，使用定时器
      if (!this.timer) {
        // 一旦监听到的screenWidth值改变，就将其重新赋给data里的screenWidth
        this.screenWidth = val;
        this.timer = true;
        let self = this;
        setTimeout(function() {
          // 打印screenWidth变化的值
          console.log(self.screenWidth);
          self.timer = false;
        }, 400);
      }
    }
  },
  mounted() {
    const self = this;
    window.onresize = () => {
      return (() => {
        window.screenWidth = document.body.clientWidth;
        self.screenWidth = window.screenWidth;
      })();
    };
  }
};
</script>

<style lang="scss" scoped type="text/css">
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.5s ease;
  max-width: 1000px;
}
.slide-fade-enter, .slide-fade-leave-to
/* .slide-fade-leave-active for below version 2.1.8 */ {
  max-width: 0px;
  opacity: 1;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>

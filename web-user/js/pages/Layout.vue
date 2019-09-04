<template lang="pug">
  .front__page
    title(v-if="$route.meta.title") {{$route.meta.title}}
    top-bar(:screenWidth="screenWidth")
    transition(name='fade', mode="out-in", :duration="{ enter: 500, leave: 0 }",appear, appear-active-class='animated fadeInUp fade-enter-active' enter-active-class='animated fadeInUp fade-enter-active' leave-active-class='animated fade-leave-active')
      router-view(v-if="$route.meta.keepAlive === false",:screenWidth="screenWidth")
    transition(name='fade', mode="out-in", :duration="{ enter: 500, leave: 0 }",appear, appear-active-class='animated fadeInUp fade-enter-active' enter-active-class='animated fadeInUp fade-enter-active' leave-active-class='animated fade-leave-active')
      keep-alive
        router-view(v-if="$route.meta.keepAlive !== false",:screenWidth="screenWidth")
</template>

<script>
import TopBar from "@/web-user/js/common/topbar.vue";

export default {
  name: "app",
  components: {
    TopBar
  },
  data() {
    return {
      screenWidth: document.body.clientWidth
    };
  },
  mounted() {
    const self = this;
    window.onresize = () => {
      return (() => {
        window.screenWidth = document.body.clientWidth;
        self.screenWidth = window.screenWidth;
      })();
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
  }
};
</script>

<style lang="scss" scoped type="text/css">
.animated {
  animation-duration: 0.5s;
  animation-fill-mode: both;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
.userlayout__transtion {
  animation-duration: 3s;
  animation-delay: 0s;
}
</style>

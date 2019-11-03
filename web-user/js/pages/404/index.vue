<template lang="pug">
  .content
    .content__main
      .home__wrapper
        h1 对不起，你要访问的页面不存在，即将跳转到主页
        img(src="@/web-common/assets/images/404.jpg")
</template>

<script>
export default {
  data() {
    return {
      timer: 0
    };
  },
  mounted() {
    const self = this;
    self.timer = setTimeout(() => {
      console.log('index');
      self.$router.replace({name: 'index'});
    }, 3000);
  },
  beforeDestroy() {
    // 关闭定时器
    if (this.timer) {
      clearTimeout(this.timer);
    }
  }
};
</script>

<style lang="scss" scoped>
.home__wrapper {
  height: 100%;
  background: $--color-level15;
}
</style>
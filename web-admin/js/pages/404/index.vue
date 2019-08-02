<template lang="pug">
.admin-content
  .content__main
    h1 对不起，你要访问的页面不存在，即将跳转到主页
    img(:src="imgUrl('static/images/404.jpg')")
</template>

<script>
export default {
  name: "",
  data() {
    return {};
  },
  mounted() {
    const self = this;
    setTimeout(() => {
      console.log("home");
      self.$router.replace({ name: "home" });
    }, 3000);
  }
};
</script>

<style lang="scss" scoped>
</style>
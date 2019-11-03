<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .form__wrapper
      el-form(ref="form",:model="form")
        el-form-item(label="讨论版")
          el-switch(v-model="form.enable_issue",active-text="开",inactive-text="关",active-value=true,inactive-value=false)
        el-form-item
          el-button(type="primary",@click="submitSetSettings") 保存
</template>
<script>
import { getSettings, setSettings } from "@/web-admin/js/api/settings.js";

export default {
  data() {
    return {
      form: {
        enable_issue: true
      }
    };
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      try {
        let res = await getSettings();
        let data = res.data;
        self.form = data.config;
      } catch (err) {
        console.log(err);
      }
    },
    async submitSetSettings() {
      const self = this;
      try {
        let res = await setSettings(self.form);
        self.$message({
          type: "success",
          message: res.data.message
        });
      } catch (err) {
        console.log(err);
        self.$message({
          type: "info",
          message: err.response.data.message
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.content__generator__title {
  margin-bottom: 1em;
  color: $--color-level3;
}
.content_generator_card {
  height: 400px;
}
.content__main__item {
  height: 40em;
  padding: 2em;
  .el-card {
    height: 100%;
  }
}
.dialog__info {
  font-size: 16px;
  text-align: left;
  padding: 0.5em;
}
</style>
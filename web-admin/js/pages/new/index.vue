<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-form(:model="form",:rules="rules",ref="form",label-width="75px")
      el-form-item(label="标题", prop="title")
       el-input(v-model="form.title",placeholder="请输入标题",:autofocus="true")
      el-form-item(label="内容")
        tinymce-editor(v-model="form.content",:height="600")
      el-form-item(class="fl")
        el-button(type="primary",@click="submit") 提交
</template>

<script>
import TinymceEditor from "@/web-common/components/tinymce_editor.vue";
import { createNew, editNew, getNew } from "@/web-admin/js/api/new.js";
export default {
  components: {
    TinymceEditor
  },
  data() {
    return {
      form: {
        title: "",
        content: ""
      },
      rules: {
        title: [
          {
            required: true,
            message: "请输入标题",
            trigger: "blur"
          },
          {
            max: 20,
            message: "超出长度限制",
            trigger: "blur"
          }
        ]
      }
    };
  },
  async mounted() {
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      try {
        if (self.$route.name == "adminEditNew") {
          let id = self.$route.params.id;
          let res = await getNew(id);
          self.form.title = res.data.new.title;
          self.form.content = res.data.new.content;
        } else {
          Object.assign(self.form, {
            title: "",
            content: ""
          });
        }
      } catch (err) {
        console.log(err);
        self.$router.replace({ name: "admin404Page" });
      }
    },
    async submit() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res;
            if (self.$route.name == "adminAddNew") {
              res = await createNew(self.form);
            } else {
              let id = self.$route.params.id;
              res = await editNew(id, self.form);
            }
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.$router.push({ name: "adminNewlist" });
          } catch (err) {
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
        } else {
          self.$message({
            message: "表单必填项不能为空",
            type: "error"
          });
          return false;
        }
      });
    }
  },
  watch: {
    $route(to, from) {
      if (to.name == "adminAddNew" || to.name == "adminEditNew") {
        this.init();
      }
    }
  }
};
</script>

 <style lang="scss" scoped>
</style>
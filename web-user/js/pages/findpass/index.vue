<template lang="pug">
  .content
    .content__main
      .one-main
        h1.content__panel__title 找回密码
        .steps-wrapper
          el-steps(:active="active",finish-status="success")
            el-step(title="确认账号")
            el-step(title="重设密码")
        .form-wrapper
          el-form(ref="form", :model="form", :rules="formRules")
            el-form-item(prop="email")
              el-input(v-model="form.email", placeholder="请输入绑定邮箱地址")
                svg-icon(slot="prefix", name="email", class="auth__input__prefix")
            el-form-item
              el-button(type="primary",style="width:100%;",@click="sendEmail") 发送确认邮件
</template>

<script>
import { getNewList } from "@/web-user/js/api/nologin.js";
import { sendFindPassEmail } from "@/web-user/js/api/auth.js";
export default {
  name: "",
  data() {
    return {
      active: 0,
      form: {
        email: ""
      },
      formRules: {
        email: [
          {
            required: true,
            message: "请输入邮箱地址",
            trigger: "blur"
          },
          {
            type: "email",
            message: "请输入正确的邮箱地址",
            trigger: "blur"
          }
        ]
      }
    };
  },
  mounted() {},
  methods: {
    sendEmail() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res = await sendFindPassEmail(self.form);
            self.$message({
              message: res.data.message,
              type: "success"
            });
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
        } else {
          return false;
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.steps-wrapper {
  padding-top: 30px;
  width: 600px;
  margin: 0 auto;
}
.form-wrapper {
  width: 600px;
  margin: 30px auto 0;
  svg {
    margin-left: -5px;
    height: 20px;
  }
}
</style>
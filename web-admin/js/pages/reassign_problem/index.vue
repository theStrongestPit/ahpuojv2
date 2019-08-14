<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-card.box-card.rejudge-problem__card
      el-form(:model="form",ref="form",:rules="rules",@submit.native.prevent)
       el-form-item(label="原ID",prop="oldId")
        el-input(v-model.number="form.oldId",@keyup.enter.native="submitRejudgeSolution",placeholder="请输入要移动的问题ID") 
       el-form-item(label="新ID",prop="newId")
        el-input(v-model.number="form.newId",@keyup.enter.native="submitRejudgeSolution",placeholder="请输入要移动到的ID") 
       el-form-item
        el-button(type="primary",@click="submit") 提交
</template>

<script>
import { reassignProblem } from "@/web-admin/js/api/problem.js";
export default {
  name: "importproblem",
  data() {
    return {
      form: {
        oldId: null,
        newId: null
      },
      // problemId: null,
      rules: {
        oldId: [
          {
            required: true,
            message: "原ID不能为空",
            trigger: "blur"
          }
        ],
        newId: [
          {
            required: true,
            message: "新ID不能为空",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    submit() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res = await reassignProblem(self.form.oldId, self.form.newId);
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
          self.$message({
            message: "表单必填项不能为空",
            type: "error"
          });
          return false;
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.rejudge-problem__card {
  width: 500px;
  height: 300px;
}
</style>
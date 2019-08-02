<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`adminIndex`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-card.box-card.rejudge-problem__card
      el-form(:inline="true",:model="rejudgeSolutionForm",ref="rejudgeSolutionForm",:rules="rejudgeSolutionFormRules",@submit.native.prevent)
       el-form-item(label="重判提交",prop="id")
        el-input(v-model.number="rejudgeSolutionForm.id",@keyup.enter.native="submitRejudgeSolution",placeholder="请输入提交ID") 
       el-form-item
        el-button(type="primary",@click="submitRejudgeSolution") 提交
      el-form(:inline="true",:model="rejudgeProblemForm",ref="rejudgeProblemForm",:rules="rejudgeProblemFormRules",@submit.native.prevent)
       el-form-item(label="重判问题",prop="id")
        el-input(v-model.number="rejudgeProblemForm.id",@keyup.enter.native="submitRejudgeProblem",placeholder="请输入问题ID") 
       el-form-item
        el-button(type="primary",@click="submitRejudgeProblem") 提交
      el-alert(title="由于前台和后台管理系统是的两个独立系统，请手动搜索重判结果。" ,type="info", :closable="false")
</template>

<script>
import { rejudgeSolution, rejudgeProblem } from "@/web-admin/js/api/problem.js";
export default {
  name: "importproblem",
  data() {
    return {
      rejudgeSolutionForm: {
        id: null
      },
      rejudgeProblemForm: {
        id: null
      },
      // problemId: null,
      rejudgeSolutionFormRules: {
        id: [
          {
            required: true,
            message: "请输入提交ID名称",
            trigger: "blur"
          }
        ]
      },
      rejudgeProblemFormRules: {
        id: [
          {
            required: true,
            message: "请输入提交ID名称",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    submitRejudgeSolution() {
      const self = this;
      self.$refs["rejudgeSolutionForm"].validate(async valid => {
        if (valid) {
          try {
            console.log(self.rejudgeSolutionForm);
            let res = await rejudgeSolution(self.rejudgeSolutionForm.id);
            console.log(res);
            self.$message({
              message: res.data.message,
              type: "success"
            });

            let routerResolve = self.$router.resolve({
              name: "solution",
              params: {
                id: self.rejudgeSolutionForm.id
              }
            });
            window.open(routerResolve.href, "_blank");
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
    },
    submitRejudgeProblem() {
      const self = this;
      self.$refs["rejudgeProblemForm"].validate(async valid => {
        if (valid) {
          try {
            console.log(self.rejudgeProblemForm);
            let res = await rejudgeProblem(self.rejudgeProblemForm.id);
            console.log(res);
            self.$message({
              message: res.data.message,
              type: "success"
            });
            window.open("/status", "_blank");
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
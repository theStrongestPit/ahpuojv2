<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__card__wrapper
      el-card
        span 问题名称:{{problem.title}}
    .content__button__wrapper
      el-button(type="success" @click="handleAddData") 添加数据
    el-table(:data="tableData", style="width: 100%;margin-bottom:30px;")
      el-table-column(label="文件名",prop="filename", width="180")
      el-table-column(label="文件大小(字节)",prop="size",width="180")
      el-table-column(label="修改时间",prop="mod_time",min-width="180")
      el-table-column(label="操作", width="180")
        template(slot-scope="scope")
          el-button(size="mini",@click="handleEditData(scope.row)") 编辑
          el-button(size="mini" type="danger" @click="deleteProblemData(scope.row)") 删除
  el-dialog(:title="dialogFormTitle",:visible.sync="dialogFormVisible",@closed="closeDialog",@opened="openDialog",width="600px",:close-on-click-modal="false")
    el-form(v-show="dialogType == 'addData'",:model="addDataForm",ref="addDataForm",:rules="addDataRules",@submit.native.prevent)
      el-form-item(abel="数据名称", prop="filename")
        el-input(v-model="addDataForm.filename",autocomplete="off",@keyup.enter.native="submit")
    el-form(v-show="dialogType == 'editData'",:model="editDataForm",ref="editDataForm",:rules="editDataRules",@submit.native.prevent)
      el-form-item(label="数据内容", prop="content")
        el-input(v-model="editDataForm.content", type="textarea", :rows="20", autocomplete="off",resize="none")
    .dialog-footer(slot="footer")
      el-button(@click="cancel") 取消
      el-button(type="primary",native-type="submit",@click="submit") 确定
</template>

<script>
import {
  getProblem,
  getProblemDataList,
  editProblemData,
  addProblemData,
  deleteProblemData,
  getProblemDataFile
} from "@/web-admin/js/api/problem.js";

export default {
  data() {
    return {
      problem: null,
      dialogFormTitle: "",
      dialogFormVisible: false,
      dialogType: "",
      currentRow: null,
      addDataForm: {
        filename: ""
      },
      editDataForm: {
        content: ""
      },
      addDataRules: {
        filename: [
          {
            required: true,
            message: "请输入数据名称",
            trigger: "blur"
          },
          {
            max: 20,
            message: "超出长度限制",
            trigger: "blur"
          }
        ]
      },
      editDataRules: {
        content: [
          {
            required: true,
            message: "请输入数据内容",
            trigger: "blur"
          }
        ]
      },
      tableData: []
    };
  },
  async mounted() {
    let id = this.$route.params.id;
    try {
      let res = await getProblem(id);
      this.problem = res.data.problem;
      this.fetchProblemDataList();
    } catch (err) {
      this.$router.replace({ name: "admin404Page" });
      console.log(err);
    }
  },
  methods: {
    async fetchProblemDataList(id) {
      const self = this;
      try {
        let res = await getProblemDataList(self.problem.id);
        console.log(res);
        self.tableData = res.data.files;
      } catch (err) {
        console.log(err);
      }
    },
    openDialog() {
      if (this.dialogType == "addData") {
        this.$notify({
          title: "提示",
          message:
            "添加方法：输入数据名称点击确定，系统将生成对应的.in和.out文件，之后再编辑对应的数据文件。",
          duration: 6000
        });
      }
      console.log("123", this.$refs);
      this.$refs.addDataForm.clearValidate();
      this.$refs.editDataForm.clearValidate();
    },
    closeDialog() {
      this.$refs.addDataForm.resetFields();
      this.$refs.editDataForm.resetFields();
      this.addDataForm.filename = "";
      this.editDataForm.content = "";
    },
    submit() {
      const self = this;
      let refform = "";
      if (self.dialogType == "addData") {
        refform = "addDataForm";
      } else {
        refform = "editDataForm";
      }
      self.$refs[refform].validate(async valid => {
        if (valid) {
          try {
            let res;
            if (self.dialogType == "addData") {
              res = await addProblemData(self.problem.id, self.addDataForm);
            } else {
              res = await editProblemData(
                self.problem.id,
                self.currentRow.filename,
                self.editDataForm
              );
            }
            self.$message({
              message: res.data.message,
              type: "success"
            });
            for (let i in res.data.info) {
              setTimeout(() => {
                self.$notify({
                  title: "提示",
                  message: res.data.info[i],
                  duration: 6000
                });
              }, 500 * i);
            }
            self.fetchProblemDataList();
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
          self.dialogFormVisible = false;
        } else {
          return false;
        }
      });
    },
    cancel() {
      this.dialogFormVisible = false;
    },
    handleAddData() {
      this.dialogFormTitle = "添加数据";
      this.dialogType = "addData";
      this.dialogFormVisible = true;
    },
    async handleEditData(row) {
      this.dialogFormTitle = "编辑数据文件" + row.filename;
      this.dialogType = "editData";
      this.currentRow = row;
      let res = await getProblemDataFile(this.problem.id, row.filename);
      this.editDataForm.content = res.data.content;
      this.dialogFormVisible = true;
    },
    async deleteProblemData(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除数据文件${row.filename}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteProblemData(self.problem.id, row.filename);
          self.$message({
            type: "success",
            message: res.data.message
          });
          self.fetchProblemDataList();
        } catch (err) {
          self.$message({
            type: "error",
            message: err.response.data.message
          });
        }
      } catch (err) {
        self.$message({
          type: "info",
          message: "已取消删除"
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>
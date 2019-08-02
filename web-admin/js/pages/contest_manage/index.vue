<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`adminIndex`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-card.content__card__wrapper
      p 竞赛&作业名称: {{contest.name}}
      p 人员总数: {{total}}
    .content__button__wrapper
      el-button(type="success", @click="handleAddUser") 添加人员
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入用户名或昵称", v-model="queryParam", maxlength="20", clearable)
      el-button(icon="el-icon-search", type="primary", plain, @click="fetchContestUserList")
    el-table(:data="tableData", style="width:100%;", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="用户名", prop="username", width="180")
      el-table-column(label="昵称", prop="nick")
      el-table-column(label="操作", width="180")
        template(slot-scope="scope")
          el-button(size="mini", type="danger", @click="handleDeleteContestUser(scope.row)") 删除
    el-pagination(@size-change="handleSizeChange", @current-change="fetchContestUserList", :current-page.sync="currentPage", :page-sizes="[10, 20, 30, 40,50]", :page-size="10", layout="total, sizes, prev, pager, next, jumper", :total="total")
    el-dialog(title="添加竞赛&作业成员", :visible.sync="dialogFormVisible", @closed="closeDialog", @opened="openDialog", width="400px",:close-on-click-modal="false")
      el-form(:model="form", ref="form", :rules="rules", @submit.native.prevent)
        el-form-item(label="用户名列表", prop="userList")
          el-input(type="textarea", rows="20", v-model="form.userList", ref="input", autocomplete="off", resize="none")
      .dialog-footer(slot="footer")
        el-button(@click="cancel") 取消
        el-button(type="primary", native-type="submit", @click="submit") 确定

    el-dialog(title="操作信息",:visible.sync="dialogOperatorInfoVisible",:close-on-click-modal="false",width="600px")
      template(v-for="(item,index) in infoList")
        p(:key="index") {{item}}
      .dialog-footer(slot="footer")
        el-button(@click="dialogOperatorInfoVisible = false") 取消
        el-button(type="primary",@click="dialogOperatorInfoVisible = false") 确定
</template>

<script>
import {
  getContestUserList,
  getContest,
  addContestUser,
  deleteContestUser
} from "@/web-admin/js/api/contest.js";

export default {
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      queryParam: "",
      total: 0,
      dialogFormVisible: false,
      dialogOperatorInfoVisible: false,
      infoList: [],
      contest: null,
      form: {
        userList: ""
      },
      rules: {
        userList: [
          {
            required: true,
            message: "请输入要添加的用户名列表",
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
      let res = await getContest(id);
      this.contest = res.data.contest;
      this.fetchContestUserList();
    } catch (err) {
      this.$router.replace({ name: "admin404Page" });
      console.log(err);
    }
  },
  methods: {
    async fetchContestUserList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getContestUserList(
          self.contest.id,
          self.currentPage,
          self.perpage,
          self.queryParam
        );
        console.log(res);
        let data = res.data;
        setTimeout(() => {
          self.tableData = data.data;
          self.total = data.total;
          self.loading = false;
        }, 200);
      } catch (err) {
        console.log(err);
      }
    },
    handleSizeChange(val) {
      this.perpage = val;
      console.log(this.perpage);
      this.fetchContestUserList();
    },
    openDialog() {
      this.$notify({
        title: "提示",
        message: "每一行对应一个用户名，若对应账号存在则加入团队，否则将忽略。",
        duration: 6000
      });
      this.$refs.form.clearValidate();
      this.$refs.input.focus();
    },
    closeDialog() {
      this.$refs.form.resetFields();
      this.$refs.input.blur();
    },
    submit() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res;
            let id = self.$route.params.id;
            res = await addContestUser(id, self.form);
            self.infoList = res.data.info;
            self.dialogOperatorInfoVisible = true;
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.fetchContestUserList();
          } catch (err) {
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
    handleAddUser() {
      this.dialogFormVisible = true;
    },
    async handleDeleteContestUser(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除团队成员${row.username}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteContestUser(self.contest.id, row.id);
          self.$message({
            type: "success",
            message: res.data.message
          });
          // 删除最后一页最后一条记录，如果不是第一页，则当前页码-1
          if (self.tableData.length == 1) {
            if (self.currentPage > 1) {
              self.currentPage--;
            }
          }
          self.fetchContestUserList();
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
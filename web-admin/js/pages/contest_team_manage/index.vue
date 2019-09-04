<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    el-card.content__card__wrapper
      p 竞赛名称: {{contest.name}}
      p 包含团队总数: {{contestTeamList.length}}
    .content__button__wrapper
      el-button(type="success", @click="handleAddUser") 添加团队
    el-card
      el-collapse(v-model="activeNames")
        template(v-for="item,index in contestTeamList")
          el-collapse-item(:title="item.name",:name="item.name")
            .content__button__wrapper
              el-button(type="success", size="small",@click="handleAddContestTeamUsers(item)") 添加成员
              el-button(type="primary", size="small",@click="handleAddContestTeamAllUsers(item)") 添加全部成员
              el-button(type="danger", size="small",@click="handleDeleteContestTeam(item)") 删除团队
            el-table(:data="item.userinfos", style="width:100%;")
              el-table-column(label="ID", prop="id", width="180")
              el-table-column(label="用户名", prop="username", width="180")
              el-table-column(label="昵称", prop="nick")
              el-table-column(label="操作", width="180")
                template(slot-scope="scope")
                  el-button(size="mini", type="danger", @click="handleDeleteContestTeamUser(item,scope.row)") 删除
    el-dialog(title="添加团队", :visible.sync="dialogAddTeamFormVisible", @closed="closeAddTeamDialog", @opened="openAddTeamDialog", width="400px",:close-on-click-modal="false")
      el-form(:model="addTeamForm", ref="addTeamForm", :rules="addTeamFormRules", @submit.native.prevent)
        el-form-item(label="选择团队", prop="team_id")
          el-select(v-model="addTeamForm.team_id",filterable,placeholder="请选择")
            el-option(v-for="item in teams",:key="item.id",:label="item.name",:value="item.id")
      .dialog-footer(slot="footer")
        el-button(@click="cancelDialogAddTeam") 取消
        el-button(type="primary", native-type="submit", @click="submitAddTeam") 确定

    el-dialog(title="添加团队成员", :visible.sync="dialogAddTeamUserFormVisible", @closed="closeAddTeamUserDialog", @opened="openAddTeamUserDialog", width="400px",:close-on-click-modal="false")
      el-form(:model="addTeamUserForm", ref="addTeamUserForm", :rules="addTeamUserFormRules", @submit.native.prevent)
        el-form-item(label="用户名列表", prop="userList")
          el-input(type="textarea", rows="20", v-model="addTeamUserForm.userList", ref="input", autocomplete="off",resize="none")
      .dialog-footer(slot="footer")
        el-button(@click="cancelDialogAddTeamUser") 取消
        el-button(type="primary", native-type="submit", @click="submitAddTeamUsers") 确定
    
    el-dialog(title="操作信息",:visible.sync="dialogOperatorInfoVisible",:close-on-click-modal="false",width="600px")
      template(v-for="(item,index) in infoList")
        p(:key="index") {{item}}
      .dialog-footer(slot="footer")
        el-button(@click="dialogOperatorInfoVisible = false") 取消
        el-button(type="primary",@click="dialogOperatorInfoVisible = false") 确定
</template>

<script>
import {
  getContestTeams,
  getContest,
  addContestTeam,
  addContestTeamUsers,
  addContestTeamAllUsers,
  deleteContestTeam,
  deleteContestTeamUser
} from "@/web-admin/js/api/contest.js";

import { getAllTeams } from "@/web-admin/js/api/team.js";
import FileSaver from "file-saver";
import XLSX from "xlsx";

export default {
  data() {
    return {
      dialogAddTeamFormVisible: false,
      dialogAddTeamUserFormVisible: false,
      dialogOperatorInfoVisible: false,
      contest: null,
      teams: [],
      currentTeam: null,
      activeNames: [],
      contestTeamList: [],
      infoList: [],
      addTeamForm: {
        team_id: null
      },
      addTeamUserForm: {
        userList: ""
      },
      addTeamFormRules: {
        team_id: [
          {
            required: true,
            message: "请选择要添加的团队",
            trigger: "change"
          }
        ]
      },
      addTeamUserFormRules: {
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
      res = await getAllTeams();
      this.teams = res.data.teams;
      this.fetchDataList();
    } catch (err) {
      this.$router.replace({ name: "admin404Page" });
      console.log(err);
    }
  },
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getContestTeams(self.contest.id);
        self.contestTeamList = res.data.data;
        // 默认全部展开面板
        self.contestTeamList.forEach((val, index) => {
          self.activeNames.push(val.name);
        });
      } catch (err) {
        console.log(err);
      }
    },
    openAddTeamDialog() {
      this.$refs.addTeamForm.clearValidate();
    },
    closeAddTeamDialog() {
      this.$refs.addTeamForm.resetFields();
    },
    openAddTeamUserDialog() {
      this.$notify({
        title: "提示",
        message:
          "每一行对应一个用户名，若对应账号属于该团队并且没有以其他团队成员参与该竞赛则加入，否则将忽略。",
        duration: 6000
      });
      this.$refs.addTeamUserForm.clearValidate();
    },
    closeAddTeamUserDialog() {
      this.$refs.addTeamUserForm.resetFields();
    },
    submitAddTeam() {
      const self = this;
      self.$refs.addTeamForm.validate(async valid => {
        if (valid) {
          try {
            let res;
            let id = self.$route.params.id;
            res = await addContestTeam(id, self.addTeamForm.team_id);
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.fetchDataList();
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
          self.dialogAddTeamFormVisible = false;
        } else {
          return false;
        }
      });
    },
    submitAddTeamUsers() {
      const self = this;
      self.$refs.addTeamUserForm.validate(async valid => {
        if (valid) {
          try {
            let res;
            let id = self.$route.params.id;
            res = await addContestTeamUsers(
              id,
              self.currentTeam.id,
              self.addTeamUserForm
            );
            console.log(res);
            self.infoList = res.data.info;
            self.dialogOperatorInfoVisible = true;
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.fetchDataList();
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
          self.dialogAddTeamUserFormVisible = false;
        } else {
          return false;
        }
      });
    },
    cancelDialogAddTeam() {
      this.dialogAddTeamFormVisible = false;
    },
    cancelDialogAddTeamUser() {
      this.dialogAddTeamUserFormVisible = false;
    },
    handleAddUser() {
      this.dialogAddTeamFormVisible = true;
    },
    handleAddContestTeamUsers(team) {
      this.currentTeam = team;
      this.dialogAddTeamUserFormVisible = true;
    },
    async handleAddContestTeamAllUsers(team) {
      const self = this;
      try {
        await self.$confirm(`确认要添加团队${team.name}的全部成员吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await addContestTeamAllUsers(self.contest.id, team.id);
          self.infoList = res.data.info;
          self.dialogOperatorInfoVisible = true;
          self.$message({
            type: "success",
            message: res.data.message
          });
          self.fetchDataList();
        } catch (err) {
          console.log(err);
          self.$message({
            type: "error",
            message: err.response.data.message
          });
        }
      } catch (err) {
        console.log(err);
        self.$message({
          type: "info",
          message: "已取消删除"
        });
      }
    },
    async handleDeleteContestTeam(team) {
      const self = this;
      try {
        await self.$confirm(`确认要删除团队${team.name}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteContestTeam(self.contest.id, team.id);
          self.$message({
            type: "success",
            message: res.data.message
          });
          self.fetchDataList();
        } catch (err) {
          console.log(err);
          self.$message({
            type: "error",
            message: err.response.data.message
          });
        }
      } catch (err) {
        console.log(err);
        self.$message({
          type: "info",
          message: "已取消删除"
        });
      }
    },
    async handleDeleteContestTeamUser(team, row) {
      const self = this;
      try {
        await self.$confirm(
          `确认要删除团队${team.name}中成员${row.username}吗?`,
          "提示",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning"
          }
        );
        try {
          let res = await deleteContestTeamUser(
            self.contest.id,
            team.id,
            row.id
          );
          self.$message({
            type: "success",
            message: res.data.message
          });
          self.fetchDataList();
        } catch (err) {
          console.log(err);
          self.$message({
            type: "error",
            message: err.response.data.message
          });
        }
      } catch (err) {
        console.log(err);
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
<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__button__wrapper
      el-button(type="success", @click="handleCreateTeam") 新建团队
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入团队名称", v-model="queryParam", @keyup.enter.native="handleSearchByParam",maxlength="20", clearable)
      el-button(icon="el-icon-search",type="primary",plain,@click="handleSearchByParam")
    el-table(:data="tableData", style="width:100%;", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="团队名称", prop="name")
      el-table-column(label="操作", width="240")
        template(slot-scope="scope")
          el-button(size="mini", type="primary" @click="$router.push({name:'adminTeamManage',params:{id:scope.row.id}})") 管理
          el-button(size="mini", @click="handleEditTeam(scope.row)") 编辑
          el-button(size="mini", type="danger", @click="handleDeleteTeam(scope.row)") 删除
    el-pagination(@size-change="handleSizeChange",@current-change="fetchDataList",:current-page.sync="currentPage",
        :page-sizes="[10, 20, 30, 40,50]",:page-size="10",layout="total, sizes, prev, pager, next, jumper",:total="total")
  el-dialog(:title="dialogFormTitle",:visible.sync="dialogFormVisible",@closed="closeDialog",@opened="openDialog",width="400px",:close-on-click-modal="false")
    el-form(:model="form" ref="form" :rules="rules" @submit.native.prevent)
      el-form-item(label="团队名称", prop="name")
        el-input(v-model="form.name", ref="input", autocomplete="off", @keyup.enter.native="submit")
    .dialog-footer(slot="footer")
      el-button(@click="cancel") 取消
      el-button(type="primary", native-type="submit", @click="submit") 确定
</template>

<script>
import {
  getTeamList,
  createTeam,
  editTeam,
  deleteTeam
} from '@/web-admin/js/api/team.js';

export default {
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      queryParam: '',
      total: 0,
      dialogFormTitle: '',
      dialogFormVisible: false,
      submitMode: '',
      form: {
        name: ''
      },
      rules: {
        name: [
          {
            required: true,
            message: '请输入团队名称',
            trigger: 'blur'
          },
          {
            max: 20,
            message: '超出长度限制',
            trigger: 'blur'
          }
        ]
      },
      tableData: []
    };
  },
  activated() {
    this.fetchDataList();
  },
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getTeamList(
          self.currentPage,
          self.perpage,
          self.queryParam
        );
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
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchDataList();
    },
    handleSizeChange(val) {
      this.perpage = val;
      console.log(this.perpage);
      this.fetchDataList();
    },
    openDialog() {
      this.$refs.form.clearValidate();
      this.$refs.input.focus();
    },
    closeDialog() {
      this.$refs.form.resetFields();
      this.$refs.input.blur();
      this.form.name = '';
    },
    submit() {
      const self = this;
      self.$refs['form'].validate(async valid => {
        if (valid) {
          try {
            let res;
            if (self.submitMode == 'create') {
              res = await createTeam(self.form);
            } else {
              res = await editTeam(self.currentRowId, self.form);
            }
            self.$message({
              message: res.data.message,
              type: 'success'
            });
            self.fetchDataList();
          } catch (err) {
            self.$message({
              message: err.response.data.message,
              type: 'error'
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
    handleCreateTeam() {
      this.dialogFormTitle = '新建团队';
      this.submitMode = 'create';
      this.dialogFormVisible = true;
    },
    handleEditTeam(row) {
      this.dialogFormTitle = '编辑团队';
      this.submitMode = 'edit';
      this.currentRowId = row.id;
      this.form.name = row.name;
      this.dialogFormVisible = true;
    },
    async handleDeleteTeam(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除团队${row.name}吗?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        });
        try {
          let res = await deleteTeam(row.id);
          self.$message({
            type: 'success',
            message: res.data.message
          });
          self.fetchDataList();
        } catch (err) {
          self.$message({
            type: 'error',
            message: err.response.data.message
          });
        }
      } catch (err) {
        self.$message({
          type: 'info',
          message: '已取消删除'
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>
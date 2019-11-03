<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入用户名或昵称", v-model="queryParam", @keyup.enter.native="handleSearchByParam",maxlength="20", clearable)
      el-button(icon="el-icon-search", type="primary", plain, @click="handleSearchByParam")
      el-select(v-model="userType" @change="fetchDataList")
        el-option(label="普通用户",:value="0")
        el-option(label="比赛用户",:value="1")
    el-table(:data="tableData", style="width:100%;", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="用户名", width="180")
        template(slot-scope="scope")
          a(:href="`/userinfo/${scope.row.id}`",target="_blank") {{scope.row.username}}
      el-table-column(label="昵称")
        template(slot-scope="scope")
          a(:href="`/userinfo/${scope.row.id}`",target="_blank") {{scope.row.nick}}
      el-table-column(label="状态", width="180")
        template(slot-scope="scope")
          oj-tag(:type="[scope.row.defunct == 0 ? 'success':'danger']") {{scope.row.defunct == 0?'启用':'禁用'}}
      el-table-column(label="操作", width="180")
        template(slot-scope="scope")
          el-button(size="mini", type="warning", @click="handleChangePass(scope.row)") 修改密码
          el-button(size="mini", :type="scope.row.defunct == 0?'danger':'success'", @click="handleToggleUserStatus(scope.row)") {{scope.row.defunct == 0?'禁用':'启用'}}
    el-pagination(@size-change="handleSizeChange", @current-change="fetchDataList", :current-page.sync="currentPage", :page-sizes="[10, 20, 30, 40,50]", :page-size="10", layout="total, sizes, prev, pager, next, jumper", :total="total")
  el-dialog(title="修改密码", :visible.sync="dialogFormVisible", @closed="closeDialog", @opened="openDialog", width="400px",:close-on-click-modal="false")
    el-form(:model="form", ref="form", :rules="rules", @submit.native.prevent)
      el-form-item(label="新密码", prop="password")
        el-input(v-model="form.password", ref="input", autocomplete="off", @keyup.enter.native="submit")
    .dialog-footer(slot="footer")
      el-button(@click="cancel") 取消
      el-button(type="primary", native-type="submit", @click="submit") 确定
</template>

<script>
import OjTag from '@/web-common/components/ojtag';
import {
  getUserList,
  toggleUserStatus,
  changeUserPass
} from '@/web-admin/js/api/user.js';

export default {
  components: {
    OjTag
  },
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      total: 0,
      queryParam: '',
      userType: 0,
      dialogFormVisible: false,
      form: {
        password: ''
      },
      rules: {
        password: [
          {
            required: true,
            message: '请输入新的用户密码',
            trigger: 'blur'
          },
          {
            // 匹配ascii字符
            pattern: /^[\x00-\xff]+$/,
            message: '密码只能包含ascii字符',
            trigger: 'blur'
          },
          {
            min: 6,
            message: '密码最少为6位',
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
  mounted() {},
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getUserList(
          self.currentPage,
          self.perpage,
          self.userType,
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
    handleChangePass(row) {
      console.log(row);
      this.currentRowId = row.id;
      this.dialogFormVisible = true;
    },
    submit() {
      const self = this;
      self.$refs['form'].validate(async valid => {
        if (valid) {
          try {
            let res = await changeUserPass(self.currentRowId, self.form);
            self.$message({
              message: res.data.message,
              type: 'success'
            });
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
    async handleToggleUserStatus(row) {
      const self = this;
      let msg = `确认要${row.defunct == 0 ? '禁用' : '启用'}用户${
        row.username
      }吗?`;
      try {
        await self.$confirm(msg, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        });
        try {
          let res = await toggleUserStatus(row.id);
          self.$message({
            type: 'success',
            message: res.data.message
          });
          row.defunct = !row.defunct;
        } catch (err) {
          self.$message({
            type: 'error',
            message: err.response.data.message
          });
        }
      } catch (err) {
        self.$message({
          type: 'info',
          message: '已取消操作'
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>
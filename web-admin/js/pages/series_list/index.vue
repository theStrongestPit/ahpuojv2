<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__button__wrapper
      el-button(type="success", @click="handleCreateSeries") 新建系列赛
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入系列赛名称", v-model="queryParam", @keyup.enter.native="handleSearchByParam",maxlength="20", clearable)
      el-button(icon="el-icon-search",type="primary",plain,@click="handleSearchByParam")
    el-table(:data="tableData" style="width: 100%", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="名称", prop="name", min-width="300")
      el-table-column(label="状态", width="240")
        template(slot-scope="scope")
          oj-tag(:type="scope.row.defunct == 0 ? 'success':'danger'") {{scope.row.defunct == 0?'启用':'保留'}}
          oj-tag(:type="scope.row.team_mode == 0 ? 'success':'primary'") {{scope.row.team_mode == 0?'个人':'团队'}}
      el-table-column(label="操作", width="300")
        template(slot-scope="scope")
          el-button(size="mini", type="primary", @click="$router.push({name:'adminSeriesManage',params:{id:scope.row.id}})") 管理
          el-button(size="mini", type="primary", @click="handleEditSeries(scope.row)") 编辑
          el-button(size="mini", :type="scope.row.defunct == 0?'danger':'success'", @click="handleToggleSeriesStatus(scope.row)") {{scope.row.defunct == 0?'保留':'启用'}}
          el-button(size="mini", type="danger", @click="handleDeleteContest(scope.row)") 删除
    el-pagination(@size-change="handleSizeChange",@current-change="fetchDataList",:current-page.sync="currentPage",
        :page-sizes="[10, 20, 30, 40,50]",:page-size="10",layout="total, sizes, prev, pager, next, jumper",:total="total")
  el-dialog(:title="dialogFormTitle",:visible.sync="dialogFormVisible",@closed="closeDialog",@opened="openDialog",width="600px",:close-on-click-modal="false")
    el-form(:model="form", ref="form" :rules="rules", @submit.native.prevent)
      el-form-item(label="系列赛名称", prop="name")
        el-input(v-model="form.name", ref="input", autocomplete="off")
      el-form-item(label="系列赛模式", prop="team_mode")
        el-switch(v-model="form.team_mode", active-text="团队", inactive-text="个人", inactive-color="#99cc33", :active-value="1", :inactive-value="0")
      el-form-item(label="系列赛描述", prop="description")
        el-input(v-model="form.description", type="textarea", :rows="5", autocomplete="off")
    .dialog-footer(slot="footer")
      el-button(@click="cancel") 取消
      el-button(type="primary", native-type="submit", @click="submit") 确定
</template>

<script>
import OjTag from '@/web-common/components/ojtag';
import {
  getSeriesList,
  createSeries,
  editSeries,
  deleteSeries,
  toggleSeriesStatus
} from '@/web-admin/js/api/series.js';

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
      queryParam: '',
      total: 0,
      dialogFormTitle: '',
      dialogFormVisible: false,
      submitMode: '',
      form: {
        name: '',
        team_mode: 0,
        description: ''
      },
      rules: {
        name: [
          {
            required: true,
            message: '请输入系列赛名称',
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
  mounted() {
    this.fetchDataList();
  },
  activated() {
    this.fetchDataList();
  },
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getSeriesList(
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
      this.$notify({
        title: '提示',
        message: '系列赛中只会显示与系列赛模式相同的比赛数据！',
        duration: 6000
      });
      this.$refs.form.clearValidate();
      this.$refs.input.focus();
    },
    closeDialog() {
      this.$refs.form.resetFields();
      this.$refs.input.blur();
      this.form.name = '';
      this.form.description = '';
      this.form.team_mode = 0;
    },
    submit() {
      const self = this;
      self.$refs['form'].validate(async valid => {
        if (valid) {
          try {
            let res;
            if (self.submitMode == 'create') {
              res = await createSeries(self.form);
            } else {
              res = await editSeries(self.currentRowId, self.form);
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
    handleCreateSeries() {
      this.dialogFormTitle = '新建系列赛';
      this.submitMode = 'create';
      this.dialogFormVisible = true;
    },
    handleEditSeries(row) {
      this.dialogFormTitle = '编辑系列赛';
      this.submitMode = 'edit';
      this.currentRowId = row.id;
      this.form.name = row.name;
      this.form.description = row.description;
      this.form.team_mode = row.team_mode;
      this.dialogFormVisible = true;
    },
    async handleToggleSeriesStatus(row) {
      const self = this;
      let msg = `确认要${row.defunct == 0 ? '保留' : '启用'}系列赛${
        row.name
      }吗?`;
      try {
        await self.$confirm(msg, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        });
        try {
          let res = await toggleSeriesStatus(row.id);
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
    },
    async handleDeleteContest(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除问题${row.name}吗?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        });
        try {
          let res = await deleteContest(row.id);
          self.$message({
            type: 'success',
            message: res.data.message
          });
          self.fetchContestList();
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
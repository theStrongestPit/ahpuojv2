<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入标签名称", @keyup.enter.native="handleSearchByParam", maxlength="20", clearable)
      el-button(icon="el-icon-search",type="primary",plain,@click="handleSearchByParam")
    el-table(:data="tableData" style="width: 100%", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="名称",min-width="300")
        template(slot-scope="scope")
          a(:href="`/contest/${scope.row.id}`" target="_blank") {{scope.row.name}}
      el-table-column(label="状态", width="240")
        template(slot-scope="scope")
          span(class="text-button",:class="[scope.row.defunct == 0 ? 'text-button--success':'text-button--danger']") {{scope.row.defunct == 0?'启用':'保留'}}
          span(class="text-button",:class="[scope.row.private == 0 ? 'text-button--success':'text-button--danger']") {{scope.row.private == 0?'公开':'私有'}}
          span(class="text-button",:class="[scope.row.team_mode == 0 ? 'text-button--success':'text-button--primary']") {{scope.row.team_mode == 0?'个人':'团队'}}
      el-table-column(label="操作", width="300")
        template(slot-scope="scope")
          el-button(size="mini", type="primary", @click="$router.push({name:'adminEditContest',params:{id:scope.row.id}})") 编辑
          el-button(size="mini", @click="$router.push({name: (scope.row.team_mode == 0)? 'adminContestManage':'adminContestTeamManage' ,params:{id:scope.row.id}})", :disabled="scope.row.private == 0") 人员
          el-button(size="mini", :type="scope.row.defunct == 0?'danger':'success'", @click="handleToggleContestStatus(scope.row)") {{scope.row.defunct == 0?'保留':'启用'}}
          el-button(size="mini", type="danger", @click="handleDeleteContest(scope.row)") 删除
    el-pagination(@size-change="handleSizeChange",@current-change="fetchDataList",:current-page.sync="currentPage",
        :page-sizes="[10, 20, 30, 40,50]",:page-size="10",layout="total, sizes, prev, pager, next, jumper",:total="total")
</template>

<script>
import {
  getContestList,
  editContest,
  deleteContest,
  toggleContestStatus
} from "@/web-admin/js/api/contest.js";

export default {
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      total: 0,
      queryParam: "",
      tableData: []
    };
  },
  mounted() {
    this.fetchDataList();
  },
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getContestList(
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
      this.fetchDataList();
    },
    async handleToggleContestStatus(row) {
      const self = this;
      let msg = `确认要${row.defunct == 0 ? "保留" : "启用"}竞赛${row.name}吗?`;
      try {
        await self.$confirm(msg, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await toggleContestStatus(row.id);
          self.$message({
            type: "success",
            message: res.data.message
          });
          row.defunct = !row.defunct;
        } catch (err) {
          self.$message({
            type: "error",
            message: err.response.data.message
          });
        }
      } catch (err) {
        self.$message({
          type: "info",
          message: "已取消操作"
        });
      }
    },
    async handleDeleteContest(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除竞赛${row.name}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteContest(row.id);
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
          self.fetchDataList();
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
  },
  activated() {
    this.fetchDataList();
  }
};
</script>

<style lang="scss" scoped>
</style>
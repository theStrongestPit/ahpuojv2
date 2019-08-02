<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`adminIndex`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__searchbar__wrapper
      el-input(placeholder="请输入新闻名称", style="max-width:20em", v-model="queryParam", maxlength="20", clearable)
      el-button(icon="el-icon-search" type="primary" plain @click="fetchNewList")
    el-table(:data="tableData", style="width:100%;", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="标题", prop="title")
      el-table-column(label="状态", width="180")
        template(slot-scope="scope")
          span(class="text-button",:class="[scope.row.defunct == 0 ? 'text-button--success':'text-button--danger']") {{scope.row.defunct == 0?'显示':'隐藏'}}
          span(v-if="scope.row.top > 0", class="text-button text-button--success") 置顶
      el-table-column(label="操作", width="300")
        template(slot-scope="scope")
          el-button(size="mini", type="primary", @click="$router.push({name:'adminEditNew',params:{id:scope.row.id}})") 编辑
          el-button(size="mini", :type="scope.row.defunct == 0?'danger':'success'", @click="handleToggleNewStatus(scope.row)") {{scope.row.defunct == 0?'隐藏':'显示'}}
          el-button(size="mini", :type="scope.row.top == 0?'success':'danger'", @click="handleToggleNewTopStatus(scope.row)") {{scope.row.top == 0?'置顶':'取置'}}
          el-button(size="mini" type="danger" @click="handleDeleteNew(scope.row)") 删除
      .content__pagination__wrapper
    el-pagination(@size-change="handleSizeChange",@current-change="fetchNewList",:current-page.sync="currentPage",
        :page-sizes="[10, 20, 30, 40,50]",:page-size="10",layout="total, sizes, prev, pager, next, jumper",:total="total")

</template>

<script>
import {
  getNewList,
  editNew,
  deleteNew,
  toggleNewStatus,
  toggleNewTopStatus
} from "@/web-admin/js/api/new.js";

export default {
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      total: 0,
      queryParam: "",
      submitMode: "",
      form: {
        title: ""
      },
      rules: {
        title: [
          {
            required: true,
            message: "请输入新闻名称",
            trigger: "blur"
          },
          {
            max: 20,
            message: "超出长度限制",
            trigger: "blur"
          }
        ]
      },
      tableData: []
    };
  },
  mounted() {},
  methods: {
    async fetchNewList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getNewList(
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
    handleSizeChange(val) {
      this.perpage = val;
      this.fetchNewList();
    },
    async handleDeleteNew(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除新闻${row.title}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteNew(row.id);
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
          self.fetchNewList();
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
    },
    async handleToggleNewStatus(row) {
      const self = this;
      let msg = `确认要${row.defunct == 0 ? "隐藏" : "显示"}新闻${
        row.title
      }吗?`;
      try {
        await self.$confirm(msg, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await toggleNewStatus(row.id);
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
    async handleToggleNewTopStatus(row) {
      const self = this;
      let msg = `确认要${row.top == 0 ? "置顶" : "取置"}新闻${row.title}吗?`;
      try {
        await self.$confirm(msg, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await toggleNewTopStatus(row.id);
          self.$message({
            type: "success",
            message: res.data.message
          });
          self.fetchNewList();
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
    }
  },
  activated() {
    this.fetchNewList();
  }
};
</script>

<style lang="scss" scoped>
</style>
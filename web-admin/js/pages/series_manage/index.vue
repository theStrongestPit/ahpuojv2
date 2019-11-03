<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
      el-breadcrumb-item {{$route.meta.title}}
  .content__main
    .content__searchbar__wrapper
      el-card.content__card__wrapper
        p 系列赛名称: {{series&&series.name}}
        p 包含竞赛总数: {{tableData.length}}
    .content__button__wrapper
      el-button(type="success", @click="handleAddUser") 添加竞赛
    .content__searchbar__wrapper
      el-input(style="max-width:20em", placeholder="请输入竞赛名称", v-model="queryParam", @keyup.enter.native="handleSearchByParam",maxlength="20", clearable)
      el-button(icon="el-icon-search", type="primary", plain, @click="handleSearchByParam")
    el-table(:data="tableData", style="width:100%;", v-loading="loading")
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="名称", prop="name", width="180")
      el-table-column(label="模式", min-width="150")
        template(slot-scope="scope")
          oj-tag(:type="scope.row.private == 1 ? 'danger':'success'") {{ scope.row.private == 1?"私有赛":"公开赛" }}
          oj-tag(:type="scope.row.team_mode == 0 ? 'success':'primary'") {{ scope.row.team_mode == 0?"个人赛":"团体赛" }}
      el-table-column(label="操作", width="180")
        template(slot-scope="scope")
          el-button(size="mini", type="danger", @click="handleDeleteSeriesContest(scope.row)") 删除
    el-pagination(@size-change="handleSizeChange", @current-change="fetchDataList", :current-page.sync="currentPage", :page-sizes="[10, 20, 30, 40,50]", :page-size="10", layout="total, sizes, prev, pager, next, jumper", :total="total")
    el-dialog(title="添加竞赛", :visible.sync="dialogFormVisible", @closed="closeDialog", @opened="openDialog", width="400px",:close-on-click-modal="false")
      el-form(:model="form", ref="form", :rules="rules", @submit.native.prevent)
        el-form-item(label="选择竞赛", prop="contest_id")
          el-select(v-model="form.contest_id",filterable,placeholder="请选择")
            el-option(v-for="item in contests",:key="item.id",:label="item.name",:value="item.id")
      .dialog-footer(slot="footer")
        el-button(@click="cancel") 取消
        el-button(type="primary", native-type="submit", @click="submit") 确定
</template>

<script>
import OjTag from "@/web-common/components/ojtag";
import {
  getSeriesContestList,
  getSeries,
  addSeriesContest,
  deleteSeriesContest
} from "@/web-admin/js/api/series.js";

import { getAllContests } from "@/web-admin/js/api/contest.js";
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
      queryParam: "",
      total: 0,
      dialogFormVisible: false,
      series: null,
      contests: [],
      form: {
        contest_id: null
      },
      rules: {
        contest_id: [
          {
            required: true,
            message: "请输入要添加的用户名列表",
            trigger: "change"
          }
        ]
      },
      tableData: []
    };
  },
  async mounted() {
    let id = this.$route.params.id;
    try {
      let res = await getSeries(id);
      this.series = res.data.series;
      this.fetchDataList();
      res = await getAllContests();
      this.contests = res.data.contests;
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
        let res = await getSeriesContestList(
          self.series.id,
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
        title: "提示",
        message:
          "每一行对应一个竞赛ID，若对应竞赛存在则加入系列赛，否则将忽略。",
        duration: 6000
      });
      this.$refs.form.clearValidate();
    },
    closeDialog() {
      this.$refs.form.resetFields();
    },
    submit() {
      const self = this;
      self.$refs["form"].validate(async valid => {
        if (valid) {
          try {
            let res;
            let id = self.$route.params.id;
            res = await addSeriesContest(id, self.form.contest_id);
            console.log(res);
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
    async handleDeleteSeriesContest(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除系列赛中的竞赛${row.name}吗?`, "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        });
        try {
          let res = await deleteSeriesContest(self.series.id, row.id);
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
  }
};
</script>

<style lang="scss" scoped>
</style>
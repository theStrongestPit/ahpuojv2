<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            .tags__wrapper
              .section__title 查找竞赛：
              .siderbar__searchbar__wrapper
                el-input(style="max-width:20em", placeholder="请输入竞赛名称", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                  el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
      .main
        h1.content__panel__title 竞赛列表
        el-table(:data="tableData", style="width: 100%", class="dataTable")
          el-table-column(width="90") 
            template(slot-scope="scope")
              span(v-if="scope.row.status==1", class="text-button text-button--success") 未开始
              span(v-if="scope.row.status==2",class="text-button text-button--primary") 进行中
              span(v-if="scope.row.status==3",class="text-button text-button--danger") 已结束
          el-table-column(label="名称", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'contest',params:{id:scope.row.id}}") {{scope.row.name}}
          el-table-column(label="模式", min-width="150")
            template(slot-scope="scope")
              span(:class="['text-button', scope.row.private == 0 ? 'text-button--success':'text-button--danger']") {{ scope.row.private == 0?"公开赛":"私有赛" }}
              span(:class="['text-button', scope.row.team_mode == 0 ? 'text-button--success':'text-button--primary']") {{ scope.row.team_mode == 0?"个人赛":"团体赛" }}
          el-table-column(label="开始时间", min-width="100") 
            template(slot-scope="scope")
              span(class="contestlist__time__tag") {{spliteDate(scope.row.start_time)}}&nbsp
              span(class="contestlist__time__tag") {{spliteTime(scope.row.start_time)}}
          el-table-column(label="结束时间", min-width="100")
            template(slot-scope="scope")
              span(class="contestlist__time__tag") {{spliteDate(scope.row.end_time)}}&nbsp
              span(class="contestlist__time__tag") {{spliteTime(scope.row.end_time)}}
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(screenWidth>960?',jumper':'')",:total="total",:small="!(screenWidth>960)")
</template>

<script>
import { getContestList } from "@/web-user/js/api/nologin.js";
export default {
  name: "",
  data() {
    return {
      currentPage: 1,
      perpage: 10,
      queryParam: "",
      tableData: [],
      total: 0,
      tags: []
    };
  },
  props: {
    screenWidth: {
      type: Number
    }
  },
  methods: {
    async fetchData() {
      console.log("??");
      const self = this;
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      try {
        let res = await getContestList(
          self.currentPage,
          self.perpage,
          self.queryParam
        );
        console.log(res);
        let data = res.data;
        self.tableData = data.data;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    spliteDate(dateTimeString) {
      return new String(dateTimeString).split(" ")[0];
    },
    spliteTime(dateTimeString) {
      return new String(dateTimeString).split(" ")[1];
    }
  },
  activated() {
    this.fetchData();
  }
};
</script>
<style lang="scss" scoped>
.contestlist__time__tag {
  display: inline-block;
}
</style>
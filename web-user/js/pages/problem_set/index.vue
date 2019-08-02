<template lang="pug">
  .content
    .content__main
      .problemset__siderbar(class="fr")
        ul.problemset__siderbar__list
          li
            el-button(size="mini",round,@click="handleSearchByResetConf(0)",style="margin-top:10px;") 重置
          li
            p 查找问题：
            .siderbar__searchbar__wrapper
              el-input(style="max-width:20em", placeholder="请输入问题名或ID", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
            .tags__wrapper
          li
            p 按难度检索：
            ul.button-list
                li
                  el-button(size="mini",:class="['tag__button',level==-1?'tag__button__active':'']",@click="handleSearchByLevel(-1)",round) 全部
                li
                  el-button(size="mini",:class="['tag__button',level==0?'tag__button__active':'']", @click="handleSearchByLevel(0)",round) 简单
                li
                  el-button(size="mini",:class="['tag__button',level==1?'tag__button__active':'']", @click="handleSearchByLevel(1)",round) 中等
                li 
                  el-button(size="mini",:class="['tag__button',level==2?'tag__button__active':'']", @click="handleSearchByLevel(2)",round) 困难
          li
            p 按标签检索：
            ul.button-list
              li
                el-button(size="mini",:class="['tag__button',tagId==-1?'tag__button__active':'']", @click="handleSearchByTag(-1)", round) 全部
              template(v-for="tag in tags")
                li
                  el-button(size="mini",:class="['tag__button',tagId==tag.id?'tag__button__active':'']", @click="handleSearchByTag(tag.id)", round) {{tag.name}}
      .problemset__main
        h1.content__panel__title 问题列表
        el-table(:data="tableData", style="width: 100%", class="dataTable", v-loading="loading")
          el-table-column(width="40")
            template(slot-scope="scope")
              svg-icon(name="ok",v-if="scope.row.status == 1") 
              svg-icon(name="wrong",v-else-if="scope.row.status == 2") 
          el-table-column(label="ID", prop="id", width="60")
          el-table-column(label="标题", min-width="160")
            template(slot-scope="scope")
              router-link(:to="{name:'problem',params:{id:scope.row.id}}") {{scope.row.title}}
          el-table-column(label="难度", min-width="60")
            template(slot-scope="scope")
              el-button(v-if="scope.row.level==0", size="mini",class="text-button text-button--success") 简单
              el-button(v-if="scope.row.level==1", size="mini",class="text-button text-button--warning") 中等
              el-button(v-if="scope.row.level==2", size="mini",class="text-button text-button--danger") 困难
          el-table-column(label="标签", min-width="160")
            template(slot-scope="scope")
              el-button(v-for="tag in scope.row.tags", :key="tag.id", size="mini",class="text-button text-button--success") {{tag.name}}
          el-table-column(label="通过率", min-width="80")
            template(slot-scope="scope") {{calcRate(scope.row)}}
          el-table-column(label="通过", prop="accepted", min-width="60")
          el-table-column(label="提交", prop="submit", min-width="60")
        el-pagination.tal(@current-change="fetchProblemList",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total")
</template>

<script>
import { getProblemList, getAllTags } from "@/web-user/js/api/nologin.js";
export default {
  name: "",
  data() {
    return {
      loading: false,
      currentPage: 1,
      perpage: 20,
      queryParam: "",
      tableData: [],
      total: 0,
      level: -1,
      tagId: -1,
      tags: []
    };
  },
  async mounted() {
    this.tagId = this.$store.getters.tagId;
    console.log(this.tagId);
    this.fetchProblemList();
    let res = await getAllTags();
    this.tags = res.data.tags;
  },
  methods: {
    async fetchProblemList() {
      const self = this;
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      self.loading = true;
      try {
        let res = await getProblemList(
          self.currentPage,
          self.perpage,
          self.queryParam,
          self.level,
          self.tagId
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
      this.$store.dispatch("resetTag");
    },
    handleSearchByResetConf() {
      this.loading = true;
      this.level = -1;
      this.tagId = -1;
      this.queryParam = "";
      this.fetchProblemList();
    },
    handleSearchByParam() {
      this.loading = true;
      this.fetchProblemList();
    },
    handleSearchByLevel(level) {
      this.loading = true;
      this.level = level;
      this.fetchProblemList();
    },
    handleSearchByTag(tagId) {
      this.loading = true;
      this.tagId = tagId;
      this.fetchProblemList();
    },
    calcRate(row) {
      let rate = row.submit == 0 ? 0 : row.accepted / row.submit;
      return Number(rate * 100).toFixed(2) + "%";
    }
  },
  activated() {
    console.log(this.$store.getters.tagId);
    if (this.$store.getters.tagId != -1) {
      this.tagId = this.$store.getters.tagId;
      this.fetchProblemList();
    }
    this.$store.dispatch("resetTag");
  }
};
</script>
<style lang="scss" scoped>
.problemset__main {
  background: $c15;
  margin-right: 250px;
  svg {
    padding-left: 10px;
    width: 14px;
    vertical-align: 0px;
    height: 14px;
  }
}
.problemset__siderbar {
  text-align: left;
  min-height: 600px;
  width: 240px;
  background: $c15;
  box-sizing: border-box;
  padding: 0.1rem;
  ul.problemset__siderbar__list {
    font-size: 20px;
    & > li {
      &:not(:last-child) {
        border-bottom: 1px solid $c13;
      }
      padding: 0.1rem 0;
      &:first-child {
        padding-top: 0;
      }
    }
  }

  ul.button-list {
    display: flex;
    justify-content: flex-start;
    align-content: flex-start;
    flex-wrap: wrap;
    li {
      flex: 0 1 auto;
    }
  }
}
</style>
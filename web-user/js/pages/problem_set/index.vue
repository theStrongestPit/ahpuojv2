<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            el-button(size="mini",round,@click="handleSearchByResetConf(0)") 重置
          li
            .section__title 查找问题：
            .siderbar__searchbar__wrapper
              el-input(style="max-width:20em", placeholder="请输入问题名或ID", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
            .tags__wrapper
          li
            .section__title 按难度检索：
            ul.button-list
                li
                  el-button(size="mini",round,:class="[level == -1?'is-active':'']",@click="handleSearchByLevel(-1)") 全部
                li
                  el-button(size="mini",round,:class="[level == 0?'is-active':'']",@click="handleSearchByLevel(0)") 简单
                li
                  el-button(size="mini",round,:class="[level == 1?'is-active':'']",@click="handleSearchByLevel(1)") 中等
                li 
                  el-button(size="mini",round,:class="[level == 2?'is-active':'']",@click="handleSearchByLevel(2)") 困难
          li
            .section__title 按标签检索：
            ul.button-list
              li
                el-button(size="mini",round,:class="[tagId == -1?'is-active':'']",@click="handleSearchByTag(-1)") 全部
              template(v-for="tag in tags")
                li
                  el-button(size="mini",round,:class="[tagId == tag.id?'is-active':'']",@click="handleSearchByTag(tag.id)") {{tag.name}}
      .main
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
          el-table-column(label="难度", min-width="60",align="center")
            template(slot-scope="scope")
              oj-tag(v-if="scope.row.level==0",type="success") 简单
              oj-tag(v-if="scope.row.level==1",type="warnning") 中等
              oj-tag(v-if="scope.row.level==2",type="danger") 困难
          el-table-column(label="标签", min-width="160",align="center")
            template(slot-scope="scope")
              oj-tag(v-for="tag in scope.row.tags", :key="tag.id") {{tag.name}}
          el-table-column(label="通过率", min-width="80",align="center")
            template(slot-scope="scope") {{calcRate(scope.row)}}
          el-table-column(label="通过", prop="accepted", min-width="60",align="center")
          el-table-column(label="提交", prop="submit", min-width="60",align="center")
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(screenWidth>960?',jumper':'')",:total="total",:small="!(screenWidth>960)")
</template>

<script>
import OjTag from '@/web-common/components/ojtag';
import {getProblemList, getAllTags} from '@/web-user/js/api/nologin.js';
export default {
  components: {OjTag},
  props: {
    screenWidth: {
      type: Number,
      default: 1920
    }
  },
  data() {
    return {
      loading: false,
      currentPage: 1,
      perpage: 20,
      queryParam: '',
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
    this.fetchData();
    let res = await getAllTags();
    this.tags = res.data.tags;
  },
  activated() {
    console.log(this.$store.getters.tagId);
    if (this.$store.getters.tagId != -1) {
      this.tagId = this.$store.getters.tagId;
      this.fetchData();
    }
    this.$store.dispatch('resetTag');
  },
  methods: {
    async fetchData() {
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
      this.$store.dispatch('resetTag');
    },
    handleSearchByResetConf() {
      this.loading = true;
      this.level = -1;
      this.tagId = -1;
      this.queryParam = '';
      this.fetchData();
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    handleSearchByLevel(level) {
      this.currentPage = 1;
      this.loading = true;
      this.level = level;
      this.fetchData();
    },
    handleSearchByTag(tagId) {
      this.currentPage = 1;
      this.loading = true;
      this.tagId = tagId;
      this.fetchData();
    },
    calcRate(row) {
      let rate = row.submit == 0 ? 0 : row.accepted / row.submit;
      return Number(rate * 100).toFixed(2) + '%';
    }
  }
};
</script>
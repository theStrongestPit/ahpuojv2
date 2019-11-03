<template lang="pug">
  .content
    .content__main
      .one-main
        h1.content__panel__title 排名
        el-table(:data="tableData", style="width: 100%", class="dataTable", v-loading="loading")
          el-table-column(label="排名", width="60")
            template(slot-scope="scope") 
              span {{ (currentPage-1) * 50 + scope.$index + 1}}
          el-table-column(label="用户",width="70")
            template(slot-scope="scope")
                router-link(:to="{name:'userinfo',params:{id:scope.row.id}}") 
                  .user__avatar__wrapper
                    img(:src="imgUrl(scope.row.avatar)",class="user__avatar")
          el-table-column
            template(slot-scope="scope")
              router-link(:to="{name:'userinfo',params:{id:scope.row.id}}") 
                span {{`${scope.row.nick}`}}                    
          el-table-column(label="通过率", width="80")
            template(slot-scope="scope") {{calcRate(scope.row)}}                
          el-table-column(label="解决",width="70",prop="solved")
          el-table-column(label="提交",width="70",prop="submit")

        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(screenWidth>960?',jumper':'')",:total="total",:small="!(screenWidth>960)")
</template>

<script>
import {getRankList} from '@/web-user/js/api/nologin.js';
import {setTimeout} from 'timers';
export default {
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
      perpage: 50,
      problemId: 0,
      tableData: [],
      total: 0
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    test(row) {
      console.log(row);
    },
    async fetchData() {
      const self = this;
      self.loading = true;
      try {
        let res = await getRankList(self.currentPage, self.perpage);
        console.log(res);
        setTimeout(() => {
          let data = res.data;
          self.tableData = data.data;
          self.total = data.total;
          self.loading = false;
        }, 200);
      } catch (err) {
        console.log(err);
      }
    },
    calcRate(row) {
      let rate = row.submit == 0 ? 0 : row.solved / row.submit;
      return Number(rate * 100).toFixed(2) + '%';
    }
  }
};
</script>

<style lang="scss" scoped>
.user__avatar__wrapper {
  img {
    width: 50px;
    height: 50px;
    border-radius: 25px;
  }
}
</style>
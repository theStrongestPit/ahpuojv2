<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
           el-button(size="mini",round,@click="handleSearchByResetConf()") 重置
           el-button(size="mini",round,@click="handleSearchMine(0)",:disabled="$store.getters.username.length===0") 我的记录
          li
            .section__title 按问题检索：
            .siderbar__searchbar__wrapper
              el-input(style="max-width:20em", :placeholder=" $route.name=='status'?'请输入问题名/ID':'请输入题号(如A)'", @keyup.enter.native="handleSearchByProblem", v-model="queryParam", maxlength="20", clearable)
                el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
          li
            .section__title 按用户检索：
            .siderbar__searchbar__wrapper
              el-input(style="max-width:20em", placeholder="请输入用户昵称", @keyup.enter.native="handleSearchByNick", v-model="nick", maxlength="20", clearable)
                el-button(slot="append" icon="el-icon-search", @click="handleSearchByNick")
          li
            .section__title 按语言检索：
            ul.button-list
              li
                el-button(size="mini",round,:class="[language == -1?'is-active':'']", @click="handleSearchByLanguage(-1)",) 全部
              template(v-for="item,index in langList")
                li(v-if="item.available")
                  el-button(size="mini",round,:class="[language == index?'is-active':'']",@click="handleSearchByLanguage(index)") {{item.name}}
          li
            .section__title 按结果检索：
            ul.button-list
              li
                el-button(size="mini",round,:class="[result==-1?'is-active':'']", @click="handleSearchByResult(-1)") 全部
              template(v-for="item in searchableResultList")
                li
                  el-button(size="mini",round,:class="[result==item.code?'is-active':'']", @click="handleSearchByResult(item.code)") {{item.name}}
      .main
        h1.content__panel__title 评测记录
        el-table(:data="tableData", style="width: 100%", class="dataTable", v-loading="loading")
          el-table-column(label="ID", prop="id", width="60")
          el-table-column(label="用户",min-width="70")
            template(slot-scope="scope")
                router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}") 
                  .user__avatar__wrapper
                    img(:src="imgUrl(scope.row.user.avatar)",class="user__avatar")
          el-table-column( min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}") 
                span {{`${scope.row.user.nick}`}}
          el-table-column(label="问题", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'problem',params:{id:scope.row.problem.id}}") {{ $route.name=="status"?`P${scope.row.problem.id} ${scope.row.problem.title}`:`${engNum(scope.row.num)} ${scope.row.problem.title}` }}
          el-table-column(label="评测状态", min-width="80")
            template(slot-scope="scope")
              router-link(:to="{name:'solution',params:{id:scope.row.id}}") 
                el-button(size="mini",:type="calcRerultType(scope.row.result)") {{ resultList[scope.row.result]?resultList[scope.row.result].name:"" }}
          el-table-column(label="语言", min-width="80")
            template(slot-scope="scope") 
              span {{ langList[scope.row.language].name}}
          el-table-column(label="时间", min-width="80")
            template(slot-scope="scope") 
              span {{ `${scope.row.time}ms`}}
          el-table-column(label="内存", min-width="80")
            template(slot-scope="scope") 
              span {{ `${scope.row.memory}KB`}}
          el-table-column(label="代码长度", min-width="80")
            template(slot-scope="scope") 
              span {{ calcCodeLength(scope.row.code_length)}}
          el-table-column(label="公开", min-width="60",v-if="$route.name=='status'")
            template(slot-scope="scope") 
              span {{ scope.row.public == 1?"是":"否"}}
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(screenWidth>960?',jumper':'')",:total="total",:small="!(screenWidth>960)")
</template>

<script>
import {getSolutionList, getLanguageList} from '@/web-user/js/api/nologin.js';
import {resultList} from '@/web-common/const';
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
      perpage: 20,
      tableData: [],
      queryParam: '',
      contestId: 0,
      contestPnum: -1,
      nick: '',
      language: -1,
      result: -1,
      total: 0,
      langList: [],
      resultList: [],
      timer: 0
    };
  },
  computed: {
    searchableResultList() {
      return this.resultList.filter((val, index, arr) => {
        return val.code >= 4 && val.code <= 11;
      });
    }
  },
  watch: {
    $route(to, from) {
      if (
        (from.name == 'contestStatus' && to.name == 'status') ||
        (to.name == 'contestStatus' && from.name == 'status')
      ) {
        this.$router.replace({name: 'refresh'});
      }
    }
  },
  async mounted() {
    let res = await getLanguageList();
    this.resultList = resultList;
    this.langList = res.data.languages;
  },
  beforeDestroy() {
    // 关闭定时器
    if (this.timer) {
      clearInterval(this.timer);
    }
  },
  activated() {
    if (this.$route.name == 'contestStatus') {
      this.contestId = this.$route.params.id;
    } else {
      this.contestId = 0;
    }
    // 如果bus中记录了搜索条件 获得bus中的搜索条件
    if ('' + this.$store.getters.solutionQueryParam) {
      this.queryParam = '' + this.$store.getters.solutionQueryParam;
    }
    if ('' + this.$store.getters.solutionUserNick) {
      this.nick = '' + this.$store.getters.solutionUserNick;
    }
    if (this.$store.getters.solutionLanguage != -1) {
      this.language = this.$store.getters.solutionLanguage;
    }
    if (this.$store.getters.solutionResult != -1) {
      this.result = this.$store.getters.solutionResult;
    }
    this.$store.dispatch('resetSolutionFilter');
    // 5s请求一次数据
    this.fetchData();
    this.timer = setInterval(() => {
      this.fetchData();
    }, 5000);
  },
  deactivated() {
    // 关闭定时器
    if (this.timer) {
      clearInterval(this.timer);
    }
  },

  methods: {
    async fetchData() {
      const self = this;
      try {
        let res = await getSolutionList(
          self.currentPage,
          self.perpage,
          self.queryParam,
          self.nick,
          self.language,
          self.result,
          self.contestId
        );
        let data = res.data;
        setTimeout(() => {
          self.tableData = data.data;
          self.total = data.total;
          this.loading = false;
        }, 200);
      } catch (err) {
        console.log(err);
      }
    },
    handleSearchByResetConf() {
      this.loading = true;
      this.queryParam = '';
      this.nick = '';
      this.language = -1;
      this.result = -1;
      this.fetchData();
    },
    handleSearchMine() {
      this.currentPage = 1;
      this.loading = true;
      this.nick = this.$store.getters.userNick;
      this.fetchData();
    },
    handleSearchByProblem() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    handleSearchByNick() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    handleSearchByLanguage(language) {
      this.currentPage = 1;
      this.loading = true;
      this.language = language;
      this.fetchData();
    },
    handleSearchByResult(result) {
      this.currentPage = 1;
      this.loading = true;
      this.result = result;
      this.fetchData();
    },
    handleSearchByTag(tagId) {
      this.currentPage = 1;
      this.loading = true;
      this.tagId = tagId;
      this.fetchData();
    },
    calcRate(row) {
      let rate = row.submit == 0 ? 0 : row.solved / row.submit;
      return Number(rate * 100).toFixed(2) + '%';
    },
    calcRerultType(result) {
      if (result == 4) {
        return 'success';
      } else {
        return 'danger';
      }
    },
    calcCodeLength(codeLength) {
      return Number(codeLength / 1000).toFixed(2) + 'KB';
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
<template lang="pug">
  .content
    title {{`${contest?contest.name:''} - AHPUOJ`}}
    .content__main
      .siderbar
        ul.siderbar__item__list
          li  
            .header 竞赛信息
          li  状态：
              template(v-if="contest")
                span.text-button.text-button--success(v-if="contest.status==1") 未开始
                span.text-button.text-button--primary(v-if="contest.status==2") 进行中
                span.text-button.text-button--danger(v-if="contest.status==3") 已结束
          li 模式：
              template(v-if="contest")          
                span.text-button(:class="[contest.private == 1 ? 'text-button--danger':'text-button--success']")  {{ contest.private == 1?"私有赛":"公开赛" }}
                span.text-button(:class="[contest.team_mode == 0 ? 'text-button--success':'text-button--primary']")  {{ contest.team_mode == 0?"个人赛":"团体赛" }}
          li 
            div 开始时间：
            p(v-if="contest") {{contest.start_time}}
            div.mt10 结束时间：
            p(v-if="contest")  {{contest.end_time}}
          li  持续时长：
            p(v-if="contest")  {{timeDiff}}
        .button__wrapper
          el-button(size="small",type="primary",@click="jumpToStatus") 记录
          el-button(size="small",type="primary",@click="jumpToRank") 排名
          el-button(v-if="contest&&contest.team_mode == 1",size="small",type="primary",@click="jumpToTeamRank") 团队排名
      .main
        h1 {{contest?contest.name:''}}
        .main__section(style="min-height:200px;") 
          h3 竞赛简介
          div(v-if="contest",v-html="contest.description")
        .main__section
          h3 问题列表
          el-table.dataTable(v-if="seeable",:data="contest.probleminfos", style="width: 100%")
            el-table-column(width="40")
              template(slot-scope="scope")
                svg-icon(name="ok",v-if="scope.row.status == 1") 
                svg-icon(name="wrong",v-else-if="scope.row.status == 2") 
            el-table-column(label="#", width="60")
              template(slot-scope="scope")
                span {{ engNum(scope.row.num) }}
            el-table-column(label="标题", min-width="180")
              template(slot-scope="scope")
                router-link(:to="{name:'contestProblem',params:{id:contest.id,num:scope.row.num}}", target="_blank") {{scope.row.title}}
          p(v-else) {{reason}}
</template>

<script>
import { getContest } from "@/web-user/js/api/nologin.js";
import CodeMirror from "@/web-common/components/codemirror.vue";
import { getLanguageList } from "@/web-user/js/api/nologin.js";
import { EventBus } from "@/web-common/eventbus";
import { submitJudgeCode } from "@/web-user/js/api/user.js";
export default {
  components: {
    CodeMirror
  },
  data() {
    return {
      seeable: false,
      reason: "",
      contest: null,
      langList: []
    };
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      console.log("initing");
      const self = this;
      let res = await getLanguageList();
      this.langList = res.data.languages;
      let id = self.$route.params.id;
      try {
        let res = await getContest(id);
        console.log(res);
        let data = res.data;
        self.contest = data.contest;
        self.seeable = data.seeable;
        self.reason = data.reason;
      } catch (err) {
        console.log(err);
        self.$router.replace({ name: "404Page" });
      }
    },
    jumpToStatus() {
      const self = this;
      let routerResolve = self.$router.resolve({
        name: "contestStatus",
        params: {
          id: self.contest.id
        }
      });
      window.open(routerResolve.href, "_blank");
    },
    jumpToRank() {
      const self = this;
      let routerResolve = self.$router.resolve({
        name: "contestRank",
        params: {
          id: self.contest.id
        }
      });
      window.open(routerResolve.href, "_blank");
    },
    jumpToTeamRank() {
      const self = this;
      let routerResolve = self.$router.resolve({
        name: "contestTeamRank",
        params: {
          id: self.contest.id
        }
      });
      window.open(routerResolve.href, "_blank");
    }
  },
  computed: {
    contestStatus() {
      let startDate = new Date(this.contest.start_time);
      let endDate = new Date(this.contest.end_time);
      let nowDate = new Date();
      if (nowDate.getTime() < startDate.getTime()) {
        // 未开始
        return 0;
      } else if (nowDate.getTime() > endDate.getTime()) {
        // 已结束
        return 2;
      } else {
        // 进行中
        return 1;
      }
    },
    timeDiff() {
      let startDate = new Date(this.contest.start_time);
      let endDate = new Date(this.contest.end_time);
      let dateDiff = endDate.getTime() - startDate.getTime();

      let days = Math.floor(dateDiff / (24 * 3600 * 1000));
      let left = dateDiff % (24 * 3600 * 1000);
      let hours = Math.floor(left / (3600 * 1000));
      left = left % (3600 * 1000);
      let minutes = Math.floor(left / (60 * 1000));
      left = left % (60 * 1000);
      let seconds = Math.round(left / 1000);

      let res = "";
      res += days ? `${days}天` : "";
      res += hours ? `${hours}小时` : "";
      res += minutes ? `${minutes}分钟` : "";
      res += seconds ? `${seconds}秒` : "";
      return res;
    }
  },
  beforeRouteUpdate(to, from, next) {
    console.log("beforeRouteUpdate!!");
    this.init();
    next();
  }
};
</script>
<style lang="scss" scoped>
</style>
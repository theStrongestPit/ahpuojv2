<template lang="pug">
  .content
    title {{user?`${user.nick}的个人空间`:''}}
    .content__main
      .userinfo__wrapper
        .userinfo__section
          h2 个人信息
          el-row.tac
            el-col(:span="12")
              .userinfo__avatart__wrapper
                img(:src="imgUrl(user?user.avatar:'')")
            el-col(:span="12")
              ul.infolist
                li
                  span 昵称 
                  span(v-if="user") {{user.nick}}
                li 
                  span 解决 
                  a(v-if="user",@click="jumpToSolved") {{user.solved}}
                li
                  span 提交 
                  a(v-if="user",@click="jumpToSubmit") {{user.submit}}
                li
                  span 排名 
                  span(v-if="user") {{user.rank}}
                li
                  span 注册时间 
                  span(v-if="user") {{user.created_at}}
        .userinfo__section
          h2 近期提交情况
          line-chart(:option="chartOption",:flag="renderFlag",:id="'chart'",style="width:100%;height:500px;")
        .userinfo__section
          h2 已解决的问题
          .problem__links(v-if="user")
            template(v-for="item,index in user.solved_problem_list")
              router-link(:to="{name:'problem',params:{id:item}}") {{item}} 
        .userinfo__section
          h2 未解决的问题
          .problem__links(v-if="user")
            template(v-for="item,index in user.unsolved_problem_list")
              router-link(:to="{name:'problem',params:{id:item}}") {{item}} 
</template>

<script>
import { getUserInfo } from "@/web-user/js/api/user.js";
import LineChart from "@/web-common/components/linechart.vue";
export default {
  name: "",
  data() {
    return {
      user: null,
      chartOption: {
        color: ["#ffdf25", "#36a9ce"],
        // title: {
        //   text: "123"
        // },
        tooltip: {},
        legend: {
          data: ["累计通过", "累计提交"]
        },
        xAxis: {
          type: "time"
        },
        yAxis: {},
        series: [
          {
            name: "累计通过",
            type: "line",
            data: []
          },
          {
            name: "累计提交",
            type: "line",
            data: []
          }
        ]
      },
      renderFlag: false
    };
  },
  components: {
    LineChart
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      try {
        let id = self.$route.params.id;
        let res = await getUserInfo(id);
        let data = res.data;
        self.user = data.userinfo;
        self.renderFlag = true;
        self.chartOption.series[0].data = self.user.recent_solved_statistic;
        self.chartOption.series[1].data = self.user.recent_submit_statistic;
      } catch (err) {
        console.log(err);
      }
    },
    jumpToSolved() {
      this.$store.dispatch("setSolutionFilter", {
        result: 4,
        nick: this.user.nick
      });
      this.$router.push({
        name: "status",
        params: {
          id: this.$route.params.id
        }
      });
    },
    jumpToSubmit() {
      this.$store.dispatch("setSolutionFilter", {
        nick: this.user.nick
      });
      this.$router.push({
        name: "status",
        params: {
          id: this.$route.params.id
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.userinfo__wrapper {
  padding: 0.15rem;
  text-align: left;
  background: $c15;
  .userinfo__section {
    min-height: 100px;
    padding: 0.15rem 0;
    &:not(:last-child) {
      border-bottom: 1px solid $c13;
      .userinfo__avatart__wrapper {
        img {
          height: 200px;
          width: 200px;
          border-radius: 100px;
          border: 1px solid $c14;
        }
      }
      ul.infolist {
        text-align: left;
        li {
          margin-top: 0.15rem;
          font-size: 20px;
          span {
            display: inline-block;
            &:first-child {
              width: 100px;
            }
          }
        }
      }
    }
    .problem__links {
      font-size: 24px;
      word-spacing: 16px;
    }
  }
}
</style>
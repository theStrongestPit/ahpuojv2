<template lang="pug">
  .content
    title {{user?`${user.nick}的个人空间`:''}}
    .content__main
      .one-main
        h1.content__panel__title 个人空间
        .main__section
          h3 用户信息
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
        .main__section
          h3 近期提交情况
          line-chart(:option="chartOption",:flag="renderFlag",:id="'chart'",style="width:100%;height:5rem;")
        .main__section
          h3 已解决的问题
          .problem__links(v-if="user")
            template(v-for="item,index in user.solved_problem_list")
              router-link(:to="{name:'problem',params:{id:item}}") {{item}} 
        .main__section
          h3 未解决的问题
          .problem__links(v-if="user")
            template(v-for="item,index in user.unsolved_problem_list")
              router-link(:to="{name:'problem',params:{id:item}}") {{item}} 
</template>

<script>
import {getUserInfo} from '@/web-user/js/api/user.js';
import LineChart from '@/web-common/components/linechart.vue';
export default {
  components: {
    LineChart
  },
  data() {
    return {
      user: null,
      chartOption: {
        color: ['#ffdf25', '#36a9ce'],
        // title: {
        //   text: "123"
        // },
        tooltip: {},
        legend: {
          data: ['累计通过', '累计提交']
        },
        xAxis: {
          type: 'time'
        },
        yAxis: {},
        series: [
          {
            name: '累计通过',
            type: 'line',
            data: []
          },
          {
            name: '累计提交',
            type: 'line',
            data: []
          }
        ]
      },
      renderFlag: false
    };
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
        console.log(self);
      } catch (err) {
        console.log(err);
      }
    },
    jumpToSolved() {
      this.$store.dispatch('setSolutionFilter', {
        result: 4,
        nick: this.user.nick
      });
      this.$router.push({
        name: 'status',
        params: {
          id: this.$route.params.id
        }
      });
    },
    jumpToSubmit() {
      this.$store.dispatch('setSolutionFilter', {
        nick: this.user.nick
      });
      this.$router.push({
        name: 'status',
        params: {
          id: this.$route.params.id
        }
      });
    }
  },
  beforeRouteUpdate(to, from, next) {
    console.log('beforeRouteUpdate!!');
    // 必须在下一次生命周期调用，因为当前路由参数还未改变，获取的路由参数还是之前的路由参数
    this.$nextTick(() => this.init());
    next();
  }
};
</script>

<style lang="scss" scoped>
.userinfo__avatart__wrapper {
  img {
    height: 2rem;
    width: 2rem;
    border-radius: 1rem;
    border: 0.01rem solid $--color-level14;
  }
}
ul.infolist {
  text-align: left;
  li {
    margin-top: 0.15rem;
    font-size: 0.16rem;
    span {
      display: inline-block;
      &:first-child {
        width: 1rem;
      }
    }
  }
}
.problem__links {
  font-size: 0.16rem;
  word-spacing: 0.16rem;
}
</style>
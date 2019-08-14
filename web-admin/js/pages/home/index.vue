<template lang="pug">
.admin-content
  .content__breadcrumb
    el-breadcrumb(separator="/")
      el-breadcrumb-item(:to="{name:`home`}") 首页
  .content__main
    .chart__wrapper
      h2.pl20.pt10 近期提交情况
      line-chart(:option="chartOption",:flag="renderFlag",:id="'chart'",style="width:100%;height:500px;")
</template>

<script>
import LineChart from "@/web-common/components/linechart.vue";
import { getSubmitStatistic } from "@/web-admin/js/api/admin.js";
export default {
  name: "",
  components: {
    LineChart
  },
  data() {
    return {
      chartOption: {
        color: ["#ffdf25"],
        // title: {
        //   text: "123"
        // },
        tooltip: {},
        legend: {
          data: ["累计提交"]
        },
        xAxis: {
          type: "time"
        },
        yAxis: {},
        series: [
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
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      try {
        let id = self.$route.params.id;
        let res = await getSubmitStatistic(id);
        let data = res.data;
        self.chartOption.series[0].data = data.recent_submit_statistic;
        self.renderFlag = true;
      } catch (err) {
        console.log(err);
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.chart__wrapper {
  background: $c15;
}
</style>
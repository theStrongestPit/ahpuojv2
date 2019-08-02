<template lang="pug">
  #mychart.chart(ref="chart",:id="id")
</template>
 
<script>
// 引入基本模板
import echarts from "echarts/lib/echarts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/bar";
import "echarts/lib/chart/radar";
// 引入提示框和图例组件
import "echarts/lib/component/tooltip";
import "echarts/lib/component/legend";

export default {
  name: "line-chart",
  props: {
    option: {
      type: Object
    },
    id: {
      type: String
    },
    flag: {
      type: Boolean
    }
  },
  data() {
    return {
      chartObj: null
    };
  },
  mounted() {
    console.log(this.flag);
  },
  methods: {
    init() {
      setTimeout(() => {
        this.chartObj = echarts.init(document.getElementById(this.id));
        this.chartObj.setOption(this.option);
      }, 500);
    }
  },
  watch: {
    flag(flag) {
      if (flag == true) {
        this.init();
      }
    }
  }
};
</script>
<style scoped>
.chart {
  height: 100%;
  width: 100%;
}
</style>
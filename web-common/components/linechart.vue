<template lang="pug">
  #mychart.chart(ref="chart",:id="id")
</template>
 
<script>
// 引入基本模板
import echarts from 'echarts/lib/echarts';
import 'echarts/lib/chart/line';
import 'echarts/lib/chart/bar';
import 'echarts/lib/chart/radar';
// 引入提示框和图例组件
import 'echarts/lib/component/tooltip';
import 'echarts/lib/component/legend';

export default {
  name: 'LineChart',
  props: {
    option: {
      type: Object,
      default() {
        return {};
      }
    },
    id: {
      type: String,
      default: 'line-chart'
    }
  },
  data() {
    return {
      chartObj: null
    };
  },
  watch: {
    option: {
      handler(val, oldval) {
        this.init();
      },
      deep: true
    }
  },
  mounted() {
    this.init();
  },
  methods: {
    init() {
      this.$nextTick(() => {
        this.chartObj = echarts.init(document.getElementById(this.id));
        this.chartObj.setOption(this.option);
      });
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
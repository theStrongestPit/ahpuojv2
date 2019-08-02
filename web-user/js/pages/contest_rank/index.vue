<template lang="pug">
  .content
    title {{contest?`竞赛&作业排名 -- ${contest.name} - AHPUOJ`:''}}
    .content__main
      .contest-ranklist__wrapper
        h1.content__panel__title {{contest?`竞赛&作业排名 -- ${contest.name}`:""}}
        el-table(v-if="seeable",:data="tableData", style="width: 100%", class="dataTable",:cell-style="cellStyle",:row-style="rowStyle")
          el-table-column(label="排名", type="index",min-width="40")
          el-table-column(label="用户名",min-width="160")
            template(slot-scope="scope")
                router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}")  {{scope.row.user.username}}
          el-table-column(label="昵称",min-width="160")
            template(slot-scope="scope")
                router-link(:to="{name:'userinfo',params:{id:scope.row.user.id}}")  {{scope.row.user.nick}}
          el-table-column( label="通过",min-width="70")
            template(slot-scope="scope")
             a.link(@click="jumpToContestStatus(scope.row)") {{scope.row.solved}}
          el-table-column(label="罚时", min-width="100")
            template(slot-scope="scope") {{secToTimeStr(scope.row.time)}}
          template(v-if="contest")
            template(v-for="count in contest.problem_count")
              el-table-column(min-width="100")
                template(slot="header" slot-scope="scope")
                  router-link(:to="{name:'contestProblem',params:{id:contest.id,num:scope.$index - problemColumnIOffset + 1}}") {{engNum(scope.$index - problemColumnIOffset + 1)}}
                template(slot-scope="scope") {{calcProblemStatus(scope.row,count)}}
        p(v-else) {{reason}}
</template>

<script>
import { getContestRankList } from "@/web-user/js/api/nologin.js";
export default {
  name: "",
  data() {
    return {
      problemColumnIOffset: 5,
      tableData: [],
      contest: null,
      timer: 0,
      seeable: false,
      reason: ""
    };
  },
  mounted() {
    this.fetctContestRankList();
    // 每隔1分钟拉取一次数据
    this.timer = setInterval(() => {
      this.fetctContestRankList();
    }, 60000);
  },
  methods: {
    async fetctContestRankList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getContestRankList(self.$route.params.id);
        console.log(res);
        let data = res.data;
        self.tableData = data.ranklist;
        self.seeable = data.seeable;
        self.reason = data.reason;
        self.contest = data.contest;
      } catch (err) {
        self.$router.replace({ name: "404Page" });
        console.log(err);
      }
    },
    handleSizeChange(val) {
      this.fetchNewList();
    },
    calcProblemStatus(row, index) {
      let res = "";
      if (row.ac_time[index - 1] > 0) {
        res += this.secToTimeStr(row.ac_time[index - 1]);
      }
      if (row.wa_count[index - 1] > 0) {
        res += `(-${row.wa_count[index - 1]})`;
      }
      return res;
    },
    cellStyle({ row, column, rowIndex, columnIndex }) {
      // 从题目的列开始计算 这一段算法照搬的hustoj的
      if (columnIndex >= this.problemColumnIOffset) {
        if (row.ac_time[columnIndex - this.problemColumnIOffset] > 0) {
          let aa =
            0x33 + row.wa_count[columnIndex - this.problemColumnIOffset] * 32;
          aa = aa > 0xaa ? 0xaa : aa;
          aa = aa.toString(16);
          let bgColor = aa + "ff" + aa;
          return `background:#${bgColor};`;
        } else if (row.wa_count[columnIndex - this.problemColumnIOffset] > 0) {
          let aa =
            0xaa - row.wa_count[columnIndex - this.problemColumnIOffset] * 10;
          aa = aa > 16 ? aa : 16;
          aa = aa.toString(16);
          let bgColor = "ff" + aa + aa;
          return `background:#${bgColor};`;
        }
      }
    },
    rowStyle({ row, rowIndex }) {
      if (row.user.username == this.$store.getters.username) {
        return `background: #f0f9eb;`;
      }
    },
    jumpToContestStatus(row) {
      this.$store.dispatch("setSolutionFilter", { nick: row.user.nick });
      this.$router.push({
        name: "contestStatus",
        params: { id: this.contest.id }
      });
    }
  },
  beforeDestroy() {
    // 关闭定时器
    if (this.timer) {
      clearInterval(this.timer);
    }
  }
};
</script>

<style lang="scss" scoped>
.contest-ranklist__wrapper {
  background: $c15;
  .link {
    padding: 0 0.1rem;
  }
}
</style>
<template lang="pug">
  .content
    title {{contest?`竞赛团队排名 -- ${contest.name} - AHPUOJ`:''}}
    .content__main
      .one-main
        el-button.fr.mr10.mt10(type="primary",@click="exportExcel") 下载excel
        h1.content__panel__title {{contest?`竞赛团队排名 -- ${contest.name}`:""}}
        el-table(v-if="seeable",:data="tableData", style="width: 100%", class="dataTable",:cell-style="cellStyle")
          el-table-column(label="排名", type="index",min-width="40")
          el-table-column(label="团队名",min-width="160",prop="team.name")
          el-table-column( label="总通过",min-width="70",prop="solved")
          el-table-column(label="总罚时", min-width="100")
            template(slot-scope="scope") {{secToTimeStr(scope.row.time)}}
          template(v-for="index in contest.problem_count")
            el-table-column(min-width="100")
              template(slot="header" slot-scope="scope")
                router-link(:to="{name:'contestProblem',params:{id:contest.id,num:scope.$index - problemColumnIOffset + 1}}") {{engNum(scope.$index - problemColumnIOffset + 1)}}
              template(slot-scope="scope") {{calcProblemStatus(scope.row,index)}}
        p(v-else) {{reason}}
</template>

<script>
import {getContestTeamRankList} from '@/web-user/js/api/nologin.js';
import FileSaver from 'file-saver';
import XLSX from 'xlsx';
export default {
  data() {
    return {
      problemColumnIOffset: 4,
      tableData: [],
      contest: null,
      timer: 0,
      seeable: false,
      reason: ''
    };
  },
  mounted() {
    this.fetctContestTeamRankList();
    // 每隔1分钟拉取一次数据
    this.timer = setInterval(() => {
      this.fetctContestTeamRankList();
    }, 60000);
  },
  beforeDestroy() {
    // 关闭定时器
    if (this.timer) {
      clearInterval(this.timer);
    }
  },
  methods: {
    async fetctContestTeamRankList() {
      const self = this;
      self.loading = true;
      try {
        let res = await getContestTeamRankList(self.$route.params.id);
        console.log(res);
        let data = res.data;
        self.tableData = data.teamranklist;
        self.seeable = data.seeable;
        self.reason = data.reason;
        self.contest = data.contest;
      } catch (err) {
        self.$router.replace({name: '404Page'});
        console.log(err);
      }
    },
    calcProblemStatus(row, index) {
      let res = '';
      if (row.ac_time[index - 1] > 0) {
        res += this.secToTimeStr(row.ac_time[index - 1]);
        res += `(${row.ac_count[index - 1]})`;
      }
      if (row.wa_count[index - 1] > 0) {
        res += `(-${row.wa_count[index - 1]})`;
      }
      return res;
    },
    cellStyle({row, column, rowIndex, columnIndex}) {
      // 从题目的列开始计算 这一段算法照搬的hustoj的
      if (columnIndex >= this.problemColumnIOffset) {
        if (row.ac_time[columnIndex - this.problemColumnIOffset] > 0) {
          let aa =
            0x33 + row.wa_count[columnIndex - this.problemColumnIOffset] * 32;
          aa = aa > 0xaa ? 0xaa : aa;
          aa = aa.toString(16);
          let bgColor = aa + 'ff' + aa;
          return `background:#${bgColor};`;
        } else if (row.wa_count[columnIndex - this.problemColumnIOffset] > 0) {
          let aa =
            0xaa - row.wa_count[columnIndex - this.problemColumnIOffset] * 10;
          aa = aa > 16 ? aa : 16;
          aa = aa.toString(16);
          let bgColor = 'ff' + aa + aa;
          return `background:#${bgColor};`;
        }
      }
    },
    jumpToContestStatus(row) {
      this.$store.dispatch('setSolutionFilter', {nick: row.user.nick});
      this.$router.push({
        name: 'contestStatus',
        params: {id: this.contest.id}
      });
    }
  },

  exportExcel() {
    /* generate workbook object from table */
    let wb = XLSX.utils.table_to_book(document.querySelector('#ranktable'));
    /* get binary string as output */
    let wbout = XLSX.write(wb, {
      bookType: 'xlsx',
      bookSST: true,
      type: 'array'
    });
    try {
      FileSaver.saveAs(
        new Blob([wbout], {type: 'application/octet-stream'}),
        `${this.contest.name}团队排名.xlsx`
      );
    } catch (e) {
      if (typeof console !== 'undefined') console.log(e, wbout);
    }
    return wbout;
  }
};
</script>

<style lang="scss" scoped>
.link {
  padding: 0 0.1rem;
}
</style>
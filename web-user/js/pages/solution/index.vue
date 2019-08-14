<template lang="pug">
  .content
    title {{solution?`S${solution.id} 评测详情 - AHPUOJ`:''}}
    .content__main
      .siderbar
        ul.siderbar__item__list
          li 
            .header 评测信息
          li(class="userinfo__wrapper align__center")
            div
              router-link(:to="{name:'userinfo',params:{id:solution?solution.user.id:0}}")
                img(:src="imgUrl(solution?solution.user.avatar:'')")
            div
              router-link(:to="{name:'userinfo',params:{id:solution?solution.user.id:0}}") {{solution?solution.user.nick:""}}
          li
            div.mt10
              strong 代码
              span(class="fr",v-if="solution") {{langList[solution.language] + " " + solution.code_length+"kb"  }}
            div.mt10
              strong 提交时间
              span(class="fr",,v-if="solution") {{solution.in_date}}
      .main
        h1.content__panel__title(style="padding-left:0;") 评测详情
        .main__section
          h3 问题
              p(v-if="solution")
                router-link(v-if="solution.contest_id==0",:to="{name:'problem',params:{id:solution.problem.id}}") {{`P${solution.problem.id}  ${solution.problem.title}`}}
                router-link(v-else,:to="{name:'contestProblem',params:{id:solution.contest_id,num:solution.num}}") {{ `C${solution.contest_id}  ${engNum(solution.num)} ${solution.problem.title}`}}
        .main__section
          h3 结果
          p(v-if="solution") {{resultList[solution.result]?resultList[solution.result].name:""}}
        .main__section
          h3 运行信息
          p(v-if="solution") {{ `用时${solution.time}ms \\ 内存${solution.memory}kb`}}
        .main__section          
          h3 编译信息
          p(v-if="solution") {{solution.compile_info?solution.compile_info:"没有编译信息"}}
        .main__section               
          h3 错误信息
          p(v-if="solution") {{renderWrongInfo}}
        // 非比赛模式下 代码提交者可以下载样例
        template(v-if="solution && solution.runtime_info && $store.getters.userId==solution.user.id && solution.contest_id == 0  && solution.result >= 5 && solution.result <= 8")
          .main__section
            h3 测试点数据下载
            el-button(size="mini",type="success",@click="handleDownloadDataFile(wrongFileName+'.in')") {{wrongFileName+".in"}}
            el-button(size="mini",type="success",@click="handleDownloadDataFile(wrongFileName+'.out')") {{wrongFileName+".out"}}
        .main__section(v-if="solution && solution.contest_id == 0 && solution.result == 4 && $store.getters.userId == solution.user.id")
          h3 公开代码
          p 公开你的源码，用你的智慧帮助其他的人解决问题！
          p.mt10 当前状态 
            span.text-button(:class="[solution.public == 0 ? 'text-button--danger':'text-button--success']") {{solution.public == 0 ? "不公开":"公开"}}
          el-button.mt10(:type="solution.public == 1?'danger':'primary'", @click="handleToggleSolutionStatus") {{solution.public == 1?'隐藏代码':'公开代码'}}
        .main__section
          h3 代码
          code-mirror(v-if="seeable",:code.sync="solution.source",:language="solution.language",:readonly="true")
          p(v-else) 你没有查看这份代码的权限


</template>

<script>
import clipboard from "clipboard";
import { getSolution } from "@/web-user/js/api/nologin.js";
import CodeMirror from "@/web-common/components/codemirror.vue";
import { downloadDatafile } from "@/web-user/js/api/user.js";
import { EventBus } from "@/web-common/eventbus";
import {
  submitJudgeCode,
  toggleSolutionStatus
} from "@/web-user/js/api/user.js";
import { resultList } from "@/web-common/const";
import { langList } from "@/web-common/const";
export default {
  components: {
    CodeMirror
  },
  data() {
    return {
      solution: null,
      langList: [],
      seeable: false,
      resultList: [],
      timer: 0
    };
  },
  mounted() {
    this.resultList = resultList;
    this.langList = langList;
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      let id = self.$route.params.id;
      try {
        let res;
        res = await getSolution(id);
        let data = res.data;
        self.seeable = data.seeable;
        self.solution = data.solution;

        // 需要重复询问
        let loading;
        if (self.solution.result < 4) {
          loading = this.$loading({
            lock: true,
            text: "评测中",
            spinner: "el-icon-loading",
            background: "rgba(0, 0, 0, 0.7)"
          });
          self.timer = setInterval(async () => {
            res = await getSolution(id);
            let data = res.data;
            self.seeable = data.seeable;
            self.solution = data.solution;
            if (self.solution.result >= 4) {
              clearInterval(self.timer);
              loading.close();
            }
          }, 2000);
        }

        console.log(self.solution);
      } catch (err) {
        console.log(err);
        self.$router.replace({ name: "404Page" });
      }
    },
    handleSearchTag(tagId) {
      this.$store.dispatch("setTag", tagId);
      this.$router.push({ name: "problemSet" });
    },
    async handleDownloadDataFile(filename) {
      try {
        let res = await downloadDatafile(
          this.solution.problem.id,
          this.solution.id,
          filename
        );
        let url = window.URL.createObjectURL(new Blob([res.data]));
        let downloadElement = document.createElement("a");
        downloadElement.style.display = "none";
        downloadElement.href = url;
        downloadElement.setAttribute("download", filename);
        document.body.appendChild(downloadElement);
        downloadElement.click();
        document.body.removeChild(downloadElement); //下载完成移除元素
        window.URL.revokeObjectURL(url); //释放掉blob对象
      } catch (err) {
        console.log(err);
      }
    },
    async handleToggleSolutionStatus() {
      const self = this;
      try {
        let res = await toggleSolutionStatus(self.solution.id);
        self.solution.public = !self.solution.public;
        this.$message({
          message: res.data.message,
          type: "success"
        });
      } catch (err) {
        console.log(err);
        this.$message({
          message: err.response.data.message,
          type: "err"
        });
      }
    }
  },
  computed: {
    wrongFileName() {
      return this.solution.runtime_info.substring(
        0,
        this.solution.runtime_info.lastIndexOf(".")
      );
    },
    renderWrongInfo() {
      if (!this.solution.runtime_info) {
        return "没有错误信息";
      }
      if (this.solution.result >= 5 && this.solution.result <= 8) {
        return "测试样例" + this.wrongFileName + "处发生了错误";
      } else {
        return solution.runtime_info;
      }
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
.userinfo__wrapper {
  img {
    width: 80px;
    height: 80px;
    border-radius: 40px;
  }
}
</style>
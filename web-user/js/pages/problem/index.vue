<template lang="pug">
  .content
    title {{problem?`${problemTitle} - AHPUOJ`:''}}
    .content__main
      .siderbar
        ul.siderbar__item__list
          li 
            .header 题目信息
          li 
            ul.submitinfo__list
              li
                p 通过
                p 
                  span(v-if="problem") {{problem.accepted}}
              li
                p 提交
                p
                  span(v-if="problem") {{problem.submit}}
          li.problem__infos
            div
              span(v-if="problem") {{`时间限制：${problem.time_limit}S`}}
            div
              span(v-if="problem") {{`内存限制：${problem.memory_limit}MB`}}
            div(v-if="$route.name == 'problem'") 难度：
              template(v-if="problem")
                span.text-button.text-button--success(v-if="problem.level==0")  简单
                span.text-button.text-button--warning(v-if="problem.level==1", size="mini") 中等
                span.text-button.text-button--danger(v-if="problem.level==2", size="mini") 困难
            li(v-if="$route.name == 'problem'") 
              .mb10 题目标签：
              ul(class="button-list",v-if="problem")
                template(v-for="tag in problem.tags")
                  li
                    el-button(size="mini",type="success",class="tag-button",@click="handleSearchTag(tag.id)") {{tag.name}}
        
        .button__wrapper
          el-button(size="small",type="primary",@click="handleChangeView") {{ problemView?'提交':'返回'}}
          el-button(size="small",type="primary",v-if="$route.name=='problem'",,@click="jumpToIssues") 讨论版
          el-button(size="small",type="primary",@click="jumpToSolutions") 记录
      .main
        h1.content__panel__title {{problemTitle}}
        transition(name="fade",mode="out-in",:duration="{ enter: 500, leave: 0 }")
          .problem_view(v-if="problemView",key="problem")
            .main__section
              h3 问题描述
              div(v-if="problem",v-html="problem.description")
            .main__section
              h3 输入
              div(v-if="problem",v-html="problem.input")
            .main__section
              h3 输出
              div(v-if="problem",v-html="problem.output")
            .main__section
              h3 样例输入
                span(@click="handleCopyInput",ref="copyInputBtn",data-clipboard-action="copy",data-clipboard-target="#input_content",class="copyBtn") 复制
              div(v-if="problem",id="input_content",style="white-space: pre-line;")  {{problem.sample_input}}
            .main__section
              h3 样例输出
                span(@click="handleCopyOutput",ref="copyOutputBtn",data-clipboard-action="copy",data-clipboard-target="#output_content",class="copyBtn")  复制
              div(v-if="problem",id="output_content",style="white-space: pre-line;") {{problem.sample_output}}
            .main__section
              h3 提示
              div(v-if="problem",v-html="problem.hint")        
          .submit_view(v-else,key="submit")
            el-form.mt5(ref="form",:model="form",class="submit__form romove__clearfix")
              el-form-item.clearfix
                el-select(v-model="form.language", placeholder="请选择")
                  template(v-for="item,index in langList")
                    el-option(v-if="item.available==true", :key="item.name", :label="item.name", :value="index")
                el-button.ml10.mt15(type="primary",size="mini",@click="dialogVisible = true") 编辑器快捷键指南
              el-form-item(style="height:550px;")
                code-mirror(:code.sync="form.source",:language="form.language")
              el-form-item
                el-row
                  el-col(:span="12")
                    span 输入
                    el-input(type="textarea",resize="none",:rows="5",v-model="testRunForm.input_text")
                  el-col(:span="12")
                    span 输出
                    el-input(type="textarea",resize="none",:rows="5",v-model="outputText",disabled,cursor="text")
              el-form-item
                el-button(type="success",@click="submitToTestRun",:disabled="testrunDisabled") {{testrunButtonText}}
                el-button(type="primary",@click="submitToJudge") 提交
    el-dialog.tal(title="快捷键指南",:visible.sync="dialogVisible",width="800px",:close-on-click-modal="false")
      p 编辑器使用sublime风格的快捷键，部分快捷键可能与系统快捷键有冲突
      p "Shift-Tab": "减少缩进",
      p "Shift-Ctrl-K": "删除当前行",
      p "Alt-Q": "wrapLines",
      p "Ctrl-T": "transposeChars",
      p "Alt-Left": "向左单位性地移动光标",
      p "Alt-Right": "向右单位性地移动光标",
      p "Ctrl-Up": "向上移动卷轴",
      p "Ctrl-Down": "向下移动卷轴",
      p "Ctrl-L": "选中当前行",
      p "Shift-Ctrl-L": "先选中多行，再按下快捷键，会在每行行尾插入光标，即可同时编辑这些行",
      p "Ctrl-Enter": "向下插入新行",
      p "Shift-Ctrl-Enter": "向上插入新行",
      p "Ctrl-D": "选中光标所占的文本，继续操作则会选中下一个相同的文本",
      p "Shift-Ctrl-Space": "选择代码块内的内容（继续选择父代码块）",
      p "Shift-Ctrl-M": "选择括号内的内容（继续选择父括号）",
      p "Ctrl-M": "光标移动至括号内结束或开始的位置",
      p "Shift-Ctrl-Up": "将光标所在行和上一行代码互换",
      p "Shift-Ctrl-Down": "将光标所在行和下一行代码互换",
      p "Ctrl-F": "打开编辑器内置搜索工具，支持正则表达式",
      p "Ctrl-H": "打开编辑器内置替换工具，支持正则表达式",
      p "Ctrl-/": "单行注释",
      p "Ctrl-J": "合并选中的多行代码为一行",
      p "Shift-Ctrl-D": "复制光标所在整行，插入到下一行",
      p "F9": "对选中的行进行排序",
      p "Ctrl-F9": "对选中的行进行非大小写敏感排序",
      p "Ctrl-K Ctrl-K": "删除当前行光标右侧内容",
      p "Ctrl-K Ctrl-Backspace": "删除当前行光标左侧内容",
      p "Ctrl-K Ctrl-U": "将光标所在单词转换为大写形式",
      p "Ctrl-K Ctrl-L": "将光标所在单词转换为小写形式",
      p "Ctrl-K Ctrl-0": "展开全部代码",
      p "Ctrl-K Ctrl-J": "展开全部代码",
      p "Ctrl-Alt-Up": "向上添加多行光标，可同时编辑多行",
      p "Ctrl-Alt-Down": "向下添加多行光标，可同时编辑多行",
      p "Ctrl-F3": "向下寻找匹配项",
      p "Shift-Ctrl-F3": "向上寻找匹配项",
      p "Alt-F3": "寻找全部匹配项",
      p "Shift-Ctrl-[": "折叠代码",
      p "Shift-Ctrl-]": "展开代码",
      p "F3": "寻找下一个匹配项",
      p "Shift-F3": "向上寻找下一个匹配项",
      span(slot="footer" class="dialog-footer")
        el-button(type="primary",@click="dialogVisible = false") 确 定
</template>

<script>
import clipboard from "clipboard";
import { testRunInterval } from "@/web-common/const";
import { getProblem, getContestProblem } from "@/web-user/js/api/nologin.js";
import CodeMirror from "@/web-common/components/codemirror.vue";
import { getLanguageList } from "@/web-user/js/api/nologin.js";
import { EventBus } from "@/web-common/eventbus";
import {
  submitJudgeCode,
  submitTestRunCode,
  getLatestSource,
  getLatestContestSource
} from "@/web-user/js/api/user.js";
import { setInterval, clearInterval } from "timers";
export default {
  components: {
    CodeMirror
  },
  data() {
    return {
      testrunDisabled: false,
      testrunButtonText: "测试运行",
      outputText: "",
      dialogVisible: false,
      problemView: true,
      problem: null,
      copyInputBtn: null,
      copyOutputBtn: null,
      langList: [],
      isFetchedLatestSource: false,
      form: {
        language: 0,
        problem_id: 0,
        contest_id: 0,
        num: 0,
        source: ""
      },
      testRunForm: {
        language: 0,
        problem_id: 0,
        source: "",
        input_text: ""
      }
    };
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      const self = this;
      let res = await getLanguageList();
      this.langList = res.data.languages;
      let id = parseInt(self.$route.params.id);
      try {
        let res;
        // 如果是普通题目路由
        if (self.$route.name == "problem") {
          res = await getProblem(id);
          let data = res.data;
          self.problem = data.problem;
          self.form.problem_id = self.problem.id;
          self.testRunForm.input_text = data.problem.sample_input;
          self.outputText = "结果应为\n" + data.problem.sample_output;
        } else {
          // 如果是比赛题目路由
          let num = parseInt(self.$route.params.num);
          res = await getContestProblem(id, num);
          console.log(res);

          let data = res.data;
          self.problem = data.problem;
          self.form.problem_id = self.problem.id;
          self.form.contest_id = id;
          self.form.num = num;
          self.testRunForm.input_text = data.problem.sample_input;
        }
        this.copyInputBtn = new clipboard(this.$refs.copyInputBtn);
        this.copyOutputBtn = new clipboard(this.$refs.copyOutputBtn);
      } catch (err) {
        console.log(err);
        self.$router.replace({ name: "404Page" });
      }
    },
    handleCopyInput() {
      const self = this;
      let clipboard = self.copyInputBtn;
      clipboard.on("success", () => {
        self.$message({
          message: "复制成功",
          type: "success"
        });
      });
      clipboard.on("error", () => {
        self.$message({
          message: "复制失败，请手动复制",
          type: "error"
        });
      });
    },
    handleCopyOutput() {
      const self = this;
      let clipboard = self.copyOutputBtn;
      clipboard.on("success", () => {
        self.$message({
          message: "复制成功",
          type: "success"
        });
      });
      clipboard.on("error", () => {
        self.$message({
          message: "复制失败，请手动复制",
          type: "error"
        });
      });
    },
    handleSearchTag(tagId) {
      this.$store.dispatch("setTag", tagId);
      this.$router.push({ name: "problemSet" });
    },
    async handleChangeView() {
      if (this.problemView == true) {
        // 切换到提交视图，需要检查是否登录
        if (this.$store.getters.username.length === 0) {
          EventBus.$emit("goLogin");
          this.$message({
            message: "请登录后提交评测",
            type: "info"
          });
        } else {
          // 拉取最近一次提交的代码
          let id = parseInt(this.$route.params.id);
          if (!this.isFetchedLatestSource) {
            let res;
            if (this.$route.name == "problem") {
              res = await getLatestSource(id);
            } else {
              let num = parseInt(this.$route.params.num);
              res = await getLatestContestSource(id, num);
            }
            this.form.language = res.data.sourcecode.language;
            this.form.source = res.data.sourcecode.source;
            this.isFetchedLatestSource = true;
          }
          this.problemView = false;
        }
      } else {
        this.problemView = true;
      }
    },
    async submitToTestRun() {
      const self = this;
      if (self.form.source.length == 0) {
        self.$message({
          message: "代码不能为空",
          type: "error"
        });
        return;
      }
      self.testRunForm.source = self.form.source;
      self.testRunForm.language = self.form.language;
      self.outputText = "正在评测中，耐心请等待.......";

      // 短暂时间内无法重复提交评测
      let countDown = testRunInterval;
      self.testrunDisabled = true;
      self.testrunButtonText = `测试运行（${countDown}）`;
      let t = setInterval(function() {
        if (countDown == 0) {
          clearInterval(t);
          self.testrunDisabled = false;
          self.testrunButtonText = `测试运行`;
        } else {
          countDown--;
          self.testrunButtonText = `测试运行（${countDown}）`;
        }
      }, 1000);

      try {
        let res = await submitTestRunCode(self.testRunForm);
        self.$message({
          message: res.data.message,
          type: "success"
        });
        self.outputText = res.data.custom_output;
      } catch (err) {
        self.$message({
          message: err.response.data.message,
          type: "error"
        });
        console.log(err);
      }
    },
    async submitToJudge() {
      const self = this;
      if (self.form.source.length == 0) {
        self.$message({
          message: "代码不能为空",
          type: "error"
        });
        return;
      }
      try {
        let res = await submitJudgeCode(self.form);
        self.$message({
          message: res.data.message,
          type: "success"
        });
        self.$router.push({
          name: "solution",
          params: { id: res.data.solution.id }
        });
      } catch (err) {
        self.$message({
          message: err.response.data.message,
          type: "error"
        });
        console.log(err);
      }
    },
    jumpToIssues() {
      let routerResolve = this.$router.resolve({
        name: "problemIssueList",
        params: {
          id: this.problem.id
        }
      });
      window.open(routerResolve.href, "_blank");
    },
    jumpToSolutions() {
      if (this.$route.name == "problem") {
        this.$store.dispatch("setSolutionFilter", {
          queryParam: this.problem.id
        });
        this.$router.push({
          name: "status"
        });
      } else {
        let num = Number.parseInt(this.$route.params.num);
        this.$store.dispatch("setSolutionFilter", {
          queryParam: this.engNum(num)
        });
        this.$router.push({
          name: "contestStatus",
          params: {
            id: this.$route.params.id
          }
        });
      }
    }
  },
  computed: {
    problemTitle() {
      if (this.$route.name == "problem" && this.problem != null) {
        return `P${this.problem.id}  ${this.problem.title}`;
      } else if (this.$route.name == "contestProblem" && this.problem != null) {
        let num = Number.parseInt(this.$route.params.num);
        let engNum = this.engNum(num);
        return `C${this.$route.params.id}  ${engNum} ${this.problem.title}`;
      }
    }
  }
};
</script>
<style lang="scss" scoped>
.siderbar__item__list {
  .submitinfo__list {
    display: flex;
    li {
      text-align: center;
      width: 50%;
      flex: 0 1 auto;
      p {
        font-size: 18px;
      }
    }
  }
  .problem__infos {
    div {
      margin: 0.1rem 0;
    }
  }
}
</style>
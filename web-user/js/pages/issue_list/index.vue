<template lang="pug">
  .content
    title {{contest?`问题P${$route.params.id}的讨论版 - AHPUOJ`:''}}
    .content__main
      .issues__wrapper
        .link.fr
          router-link(v-if="$route.name=='problemIssueList'",:to="{name:'problem',params:{id:$route.params.id}}") {{`转到问题`}}
        h1.content__panel__title {{$route.name=="issueList"?"讨论版 总版":`问题P${$route.params.id}的讨论版`}}
        template(v-for="item in data")
          .issue__box
            .issue__userinfo__wrapper
              ul
                li.issue__user__avatar
                  router-link(:to="{name:'userinfo',params:{id:item.user.id}}")
                    img(:src="imgUrl(item.user.avatar)")
                li.issue__user__name.ell 
                  router-link(:to="{name:'userinfo',params:{id:item.user.id}}") {{item.user.nick}}
            .issue__content     
              router-link(:to="{name:'issue',params:{id:item.id}}",target="_blank") {{item.title}}
              .issue__addon.tar 
                router-link(v-if="item.problem.id>0",:to="{name:'problem',params:{id:item.problem.id}}") {{`In P${item.problem.id} ${item.problem.title}`}}
                p(v-else) 总版
                p {{item.reply_count}}个回复 最后回复时间 {{item.updated_at}}
        el-pagination.tal(@current-change="fetchIssueList",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total")

        h1.content__panel__title 发表新讨论
        .post__box__wrapper
          .post__box(v-if="$store.getters.username")
            el-input(placeholder="请输入讨论标题",v-model="issueForm.title",:autofocus="true")
            tinymce-editor.mt10(v-model="replyForm.content",:height="300")
            el-button.mt10(type="primary",@click="postIssue") 发表
          a(v-else,@click="goLogin") 请登陆后发表讨论
</template>

<script>
import { getIssueList } from "@/web-user/js/api/nologin.js";
import TinymceEditor from "@/web-common/components/tinymce_editor.vue";
import { EventBus } from "@/web-common/eventbus";
import { postIssue, replyToIssue } from "@/web-user/js/api/user.js";
export default {
  components: {
    TinymceEditor
  },
  name: "",
  data() {
    return {
      currentPage: 1,
      perpage: 20,
      problemId: 0,
      data: [],
      total: 0,
      issueForm: {
        title: "",
        problem_id: 0
      },
      active: false,
      replyForm: {
        content: "",
        reply_id: 0,
        reply_user_id: 0
      }
    };
  },
  mounted() {
    this.fetchIssueList();
  },
  methods: {
    async fetchIssueList() {
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      const self = this;
      if (self.$route.name == "issueList") {
        self.problemId = 0;
      } else {
        self.problemId = self.$route.params.id;
      }
      try {
        let res = await getIssueList(
          self.problemId,
          self.currentPage,
          self.perpage
        );
        console.log(res);
        let data = res.data;
        self.data = data.data;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    },
    goLogin() {
      EventBus.$emit("goLogin");
    },
    async postIssue() {
      // 标题和内容不能为空
      const self = this;
      if (self.issueForm.title == "") {
        this.$message({
          message: "标题不能为空",
          type: "error"
        });
        return;
      }
      if (self.issueForm.title.length > 20) {
        this.$message({
          message: "标题长度限制在20个字符以内",
          type: "error"
        });
        return;
      }
      if (self.replyForm.content == "") {
        this.$message({
          message: "内容不能为空",
          type: "error"
        });
        return;
      }
      try {
        let res;
        if (self.$route.name == "problemIssueList") {
          self.issueForm.problem_id = parseInt(self.$route.params.id);
        } else {
          self.issueForm.problem_id = 0;
        }
        res = await postIssue(self.issueForm);
        let issueId = res.data.issue.id;
        res = await replyToIssue(issueId, self.replyForm);
        self.$message({
          message: "发布讨论主题成功",
          type: "success"
        });
        self.issueForm.title = "";
        self.replyForm.content = "";
        self.fetchIssueList();
      } catch (err) {
        self.$message({
          message: err.response.data.message,
          type: "error"
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.issues__wrapper {
  .link {
    font-size: 0.24rem;
    line-height: 0.5rem;
    padding-right: 0.2rem;
  }
  h1.content__panel__title {
    background: $c15;
  }
  .issue__box {
    background: $c15;
    position: relative;
    margin-top: 0.2rem;
    padding: 0.2rem;
    border: 1px solid $c13;
    .issue__userinfo__wrapper {
      text-align: center;
      word-wrap: break-word;
      word-break: break-all;
      img {
        width: 60px;
        height: 60px;
        border-radius: 30px;
        border: 1px solid $c12;
        box-shadow: 1px 1px 1px 1px $c12;
      }
      float: left;
      width: 100px;
    }
    &:last-child {
      border-bottom: 1px solid $c13;
    }
    .issue__content {
      a {
        font-size: 24px;
      }
      box-sizing: border-box;
      padding: 0.2rem 0 0 0.2rem;
      margin-left: 100px;
      min-height: 80px;
      text-align: left;
      font-size: 0.16rem;
    }
    .issue__addon {
      a {
        font-size: 0.16rem;
      }
      position: absolute;
      bottom: 10px;
      right: 10px;
    }
  }
  .post__box__wrapper {
    min-height: 100px;
    text-align: left;
    background: $c15;
    padding: 0.2rem;
  }
}
</style>
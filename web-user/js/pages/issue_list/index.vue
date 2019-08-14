<template lang="pug">
  .content
    title {{$route.name=="issueList"?`讨论版 总版 - AHPUOJ`:`问题P${$route.params.id}的讨论版 - AHPUOJ`}}
    .content__main
      .one-main(v-if="issueEnable==true")
        .link.fr
          router-link(v-if="$route.name=='problemIssueList'",:to="{name:'problem',params:{id:$route.params.id}}") {{`转到问题`}}
        h1.content__panel__title {{$route.name=="issueList"?"讨论版 总版":`问题P${$route.params.id}的讨论版`}}
        .issue__box__list
          template(v-for="item in data")
            .issue__box
              .issue__userinfo__wrapper
                ul
                  li.issue__user__avatar
                    router-link(:to="{name:'userinfo',params:{id:item.user.id}}")
                      img(:src="imgUrl(item.user.avatar)")
                  li.issue__user__name.ell 
                    router-link(:to="{name:'userinfo',params:{id:item.user.id}}") {{item.user.nick}}
              .issue__content(:class="item.is_deleted == 1?'issue-content--deleted':''")     
                router-link(:to="{name:'issue',params:{id:item.id}}",target="_blank") {{item.title}}
                .issue__addon.mt10
                  el-button.fl(v-if="$store.getters.userRole=='admin'",:type="item.is_deleted == 0?'danger':'success'",size="mini",@click="toggleIssueStatus(item.id)") {{item.is_deleted == 0 ? "删除":"恢复"}}
                  .issue__addon__info.tar
                    router-link(v-if="item.problem.id>0",:to="{name:'problem',params:{id:item.problem.id}}") {{`In P${item.problem.id} ${item.problem.title}`}}
                    p(v-else) 总版
                    p.text-muted {{item.reply_count}}条回复 最后回复时间 {{item.updated_at}}
        el-pagination.tal.mt20(@current-change="fetchIssueList",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total")
        .mt30
        h1.content__panel__title 发表新回复
        .post__box__wrapper
          .post__box(v-if="$store.getters.username")
            el-input(placeholder="请输入讨论标题",v-model="issueForm.title",:autofocus="true")
            tinymce-editor.mt10(v-model="replyForm.content",:height="300")
            el-button.mt10(type="primary",@click="postIssue") 发表
          a(v-else,@click="goLogin") 请登陆后发表讨论
      div(v-else-if="issueEnable==false")
        p 讨论版功能已经被管理员关闭
      div(v-else)
</template>

<script>
import { getIssueList } from "@/web-user/js/api/nologin.js";
import TinymceEditor from "@/web-common/components/tinymce_editor.vue";
import { EventBus } from "@/web-common/eventbus";
import { postIssue, replyToIssue } from "@/web-user/js/api/user.js";
import { toggleIssueStatus } from "@/web-user/js/api/admin.js";
export default {
  components: {
    TinymceEditor
  },
  name: "",
  data() {
    return {
      issueEnable: null,
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
        let data = res.data;

        self.issueEnable = data.issue_enable;
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
    },
    async toggleIssueStatus(issueId) {
      let self = this;
      try {
        let res = await toggleIssueStatus(issueId);
        self.$message({
          message: "变更主题状态成功",
          type: "success"
        });
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
.link {
  font-size: 0.24rem;
  line-height: 0.5rem;
  padding-right: 0.2rem;
}
.issue__box {
  background: $c15;
  position: relative;
  padding: 0.1rem;
  &:not(:last-of-type) {
    border-bottom: 1px solid $c13;
  }
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
  .issue__user__name {
    font-size: 14px;
  }
  .issue__content {
    a {
      font-size: 18px;
    }
    box-sizing: border-box;
    padding: 0.1rem 0 0 0.3rem;
    margin-left: 100px;
    min-height: 80px;
    text-align: left;
    font-size: 0.16rem;
  }
  .issue-content--deleted {
    // background-color: #f5f7fa;
    text-decoration: line-through;
  }
  .issue__addon {
    .issue__addon__info {
      a {
        font-size: 0.16rem;
      }
      position: absolute;
      bottom: 10px;
      right: 10px;
    }
  }
}
</style>
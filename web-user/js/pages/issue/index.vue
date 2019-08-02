<template lang="pug">
  .content
    title {{issue?`${issue.title} - AHPUOJ`:''}}
    .content__main
      .replys__wrapper
        .link.fr(v-if="issue")
          span 板块&nbsp;
          router-link(v-if="issue.problem.id == 0",:to="{name:'issueList'}") {{`总板`}}
          router-link(v-else,:to="{name:'problemIssueList',params:{id:issue.problem.id}}") {{`问题P ${issue.problem.title}`}}
        h1.content__panel__title {{issue?issue.title:''}}
        template(v-for="item,index in replys")
          .reply__box
            .reply__userinfo__wrapper
              ul
                li.reply__user__avatar
                  router-link(:to="{name:'userinfo',params:{id:item.user.id}}")
                    img(:src="imgUrl(item.user.avatar)")
                li.reply__user__name.ell 
                  router-link(:to="{name:'userinfo',params:{id:item.user.id}}") {{item.user.nick}}
            .reply__content
              div(v-html="item.content")
              .reply__addon.clearfix  
                p.fr
                  span  {{item.updated_at}}&nbsp;
                  span {{item.reply_count}}个回复
                el-button(type="primary",size="mini",@click="handleReplyToReply(item)") 回复
            .subreplys__wrapper.mt10
              template(v-for="subitem,subindex in item.sub_replys")
                .subreply__box
                  .subreply__userinfo__wrapper
                    ul
                      li.reply__user__avatar
                        router-link(:to="{name:'userinfo',params:{id:subitem.user.id}}")
                          img(:src="imgUrl(subitem.user.avatar)")
                      li.reply__user__name.ell 
                        router-link(:to="{name:'userinfo',params:{id:subitem.user.id}}") {{subitem.user.nick}}                  
                  .subreply__content
                    div(v-html="calcSubReply(subitem)")
                    .reply__addon.clearfix  
                      p.fr
                        span  {{item.updated_at}}&nbsp;
                        span {{item.reply_count}}个回复
        el-pagination.tal.mt10.mb10(@current-change="fetchIssue",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total",style="background:#fff;")

        h1.content__panel__title 发表新回复
        .post__box__wrapper
          .post__box(v-if="$store.getters.username")
            tinymce-editor.mt10(v-model="replyContent",:height="300")
            el-button.mt10(type="primary",@click="reply(1)") 发表
          a(v-else,@click="goLogin") 请登陆后发表讨论

    el-dialog(title="回复内容", :visible.sync="dialogFormVisible", @closed="closeDialog", @opened="openDialog", width="800px",:close-on-click-modal="false")
      tinymce-editor.mt10(v-model="subReplyContent",:height="300")
      .dialog-footer(slot="footer")
        el-button(@click="cancel") 取消
        el-button(type="primary", native-type="submit", @click="reply(2)") 确定

</template>

<script>
import { getIssue } from "@/web-user/js/api/nologin.js";
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
      dialogFormVisible: false,
      perpage: 20,
      issue: null,
      replys: [],
      total: 0,
      replyContent: "",
      subReplyContent: "",
      replyForm: {
        content: "",
        reply_id: 0,
        reply_user_id: 0
      }
    };
  },
  mounted() {
    this.fetchIssue();
  },
  methods: {
    async fetchIssue() {
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      const self = this;
      try {
        let res = await getIssue(
          self.$route.params.id,
          self.currentPage,
          self.perpage
        );
        console.log(res);
        let data = res.data;
        self.issue = data.issue;
        self.replys = data.replys;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    },
    openDialog() {
      console.log("open");
    },
    closeDialog() {
      this.subReplyContent = "";
    },
    cancel() {
      this.dialogFormVisible = false;
    },
    goLogin() {
      EventBus.$emit("goLogin");
    },
    calcSubReply(subReply) {
      if (subReply.reply_user_nick) {
        return `回复${subReply.reply_user_nick}: ${subReply.content}`;
      }
      return subReply.content;
    },
    // way = 1 表示回复主题 way = 2 表示回复某个回复
    async reply(way) {
      console.log(way);
      // 内容不能为空
      const self = this;
      self.replyForm.content =
        way == 1 ? self.replyContent : self.subReplyContent;
      if (self.replyForm.content == "") {
        this.$message({
          message: "内容不能为空",
          type: "error"
        });
        return;
      }
      try {
        let res;
        // 对主题的回复
        if (way == 1) {
          console.log("way1");
          self.replyForm.reply_id = 0;
          self.replyForm.reply_user_id = 0;
        }
        res = await replyToIssue(self.issue.id, self.replyForm);
        self.$message({
          message: res.data.message,
          type: "success"
        });
        self.replyForm.content = "";
        self.replyContent = "";
        self.dialogFormVisible = false;
        self.fetchIssue();
      } catch (err) {
        self.$message({
          message: err.response.data.message,
          type: "error"
        });
      }
    },
    handleReplyToReply(reply) {
      console.log(reply);
      this.dialogFormVisible = true;
      this.replyForm.reply_id = reply.id;
      this.replyForm.reply_user_id = reply.user.id;
    }
  }
};
</script>

<style lang="scss" scoped>
.replys__wrapper {
  .link {
    font-size: 0.24rem;
    line-height: 0.5rem;
    padding-right: 0.2rem;
  }
  h1.content__panel__title {
    background: $c15;
  }
  .reply__box {
    background: $c15;
    position: relative;
    margin-top: 0.2rem;
    padding: 0.2rem;
    border: 1px solid $c13;
    .reply__userinfo__wrapper {
      text-align: center;
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
    .reply__content {
      position: relative;
      box-sizing: border-box;
      padding: 0.1rem 0.2rem 0.2rem 0.1rem;
      margin-left: 100px;
      min-height: 100px;
      text-align: left;
      font-size: 20px;
    }
    .reply__addon {
      font-size: 16px;
      position: absolute;
      width: 100%;
      left: 0px;
      bottom: 0px;
    }
    .subreplys__wrapper {
      .subreply__box {
        padding: 0.1rem;
        margin-left: 100px;
        border: 1px solid $c13;
        .subreply__userinfo__wrapper {
          text-align: center;
          padding-top: 20px;
          img {
            width: 40px;
            height: 40px;
            border-radius: 20px;
            border: 1px solid $c12;
            box-shadow: 1px 1px 1px 1px $c12;
          }
          float: left;
          width: 100px;
        }
        .subreply__content {
          position: relative;
          box-sizing: border-box;
          padding: 0.15rem;
          margin-left: 100px;
          min-height: 100px;
          text-align: left;
          font-size: 20px;
        }
      }
    }
  }
  .post__box__wrapper {
    min-height: 100px;
    text-align: left;
    background: $c15;
    padding: 1rem;
  }
}
</style>
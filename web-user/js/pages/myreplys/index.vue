<template lang="pug">
  .content
    .content__main
      .replys__wrapper
        h1.content__panel__title 回复列表
        div(v-if="replys&&replys.length==0")
          .reply__box
            p 还没有人回复你
        div(v-else)
          template(v-for="item,index in replys")
            .reply__box
              router-link(:to="{name:'userinfo',params:{id:item?item.user.id:0}}") {{item.user.username}}
              span 在帖子
              router-link(:to="{name:'issue',params:{id:item?item.issue_id:0}}") {{item.issue_title}}
              span 中回复了你
              br
              .reply__content(v-html="calcContent(item.content)")
        el-pagination.tal.mt10.mb10(@current-change="fetchMyReplys",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total",style="background:#fff;")
</template>

<script>
import TinymceEditor from "@/web-common/components/tinymce_editor.vue";
import { EventBus } from "@/web-common/eventbus";
import { getMyReplys } from "@/web-user/js/api/user.js";
export default {
  components: {
    TinymceEditor
  },
  data() {
    return {
      currentPage: 1,
      dialogFormVisible: false,
      perpage: 20,
      replys: [],
      total: 0
    };
  },
  mounted() {
    this.fetchMyReplys();
  },
  methods: {
    async fetchMyReplys(resetScroll) {
      if (resetScroll != false) {
        window.pageYOffset = 0;
        document.documentElement.scrollTop = 0;
        document.body.scrollTop = 0;
      }
      const self = this;
      try {
        let res = await getMyReplys(self.currentPage, self.perpage);
        let data = res.data;
        console.log(res);
        self.replys = data.replys;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    },
    calcContent(content) {
      return content.length <= 100 ? content : content.substr(0, 100) + "...";
    }
  }
};
</script>

<style lang="scss" scoped>
.replys__wrapper {
  h1.content__panel__title {
    background: $c15;
  }
  .reply__box {
    background: $c15;
    position: relative;
    margin-top: 0.2rem;
    padding: 0.2rem;
    border: 1px solid $c13;
    text-align: left;
    &:last-child {
      border-bottom: 1px solid $c13;
    }
    .reply__content {
      background: $c14;
      padding: 0.1rem;
      position: relative;
      box-sizing: border-box;
      min-height: 60px;
      text-align: left;
      font-size: 16px;
    }
  }
}
</style>
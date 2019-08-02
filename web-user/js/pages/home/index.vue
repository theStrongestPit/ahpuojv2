<template lang="pug">
  .content
    .content__main
      .home__wrapper
        p(class="welcome__title") 欢迎使用AHPUOJ
        .carousel__wrapper
          el-carousel(trigger="click",height="400px", indicator-position="outside")
            el-carousel-item
              img(:src="imgUrl('/static/images/acm.jpg')")
        template(v-for="item in newList")
          .new__box
            .new__title {{item.title}}
            .new__content(v-html="item.content")
            .new__time {{item.updated_at}}
        el-pagination.tal(@current-change="fetchNewList",:current-page.sync="currentPage",background,
        :page-size="perpage",layout="prev, pager, next,jumper",:total="total")
</template>

<script>
import { getNewList } from "@/web-user/js/api/nologin.js";
export default {
  name: "",
  data() {
    return {
      currentPage: 1,
      perpage: 5,
      newList: [],
      total: 0
    };
  },
  mounted() {
    this.fetchNewList();
  },
  methods: {
    async fetchNewList() {
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      const self = this;
      try {
        let res = await getNewList(self.currentPage, self.perpage);
        console.log(res);
        let data = res.data;
        self.newList = data.data;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.home__wrapper {
  background: $c15;
  .welcome__title {
    padding-top: 0.2rem;
    font-size: 40px;
  }
  .carousel__wrapper {
    padding: 0 2rem;
    .el-carousel {
      overflow: hidden;
    }
  }
  .new__box {
    position: relative;
    padding: 0.3rem 0.5rem;
    border-top: 1px solid $c13;
    &:last-child {
      border-bottom: 1px solid $c13;
    }
    .new__title {
      text-align: left;
      font-size: 30px;
      color: $c0;
    }
    .new__content {
      min-height: 200px;
      text-align: left;
      font-size: 24px;
    }
    .new__time {
      position: absolute;
      bottom: 5px;
      right: 5px;
    }
  }
}
</style>
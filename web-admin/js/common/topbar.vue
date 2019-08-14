<template lang="pug">
.topbar__wrapper
  .topbar
    .topbar__title AHPUOJ后台管理系统
    .topbar__userinfo.fr(v-if="screenWidth > 960") {{`欢迎您，${$store.getters.userNick}`}}
    .topbar__mobile_nav(v-else)
      .mobile-humber(@click="toggleMobileNav")
        a(href="#", :class="{active:showMobileNav}")
          span(class="line")
          span(class="line")
          span(class="line")
  transition(name="slide-fade")
    el-menu.tal(:default-active="defaultActive",@select="toggleMobileNav",class="topbar__mobile__nav__menu",background-color="#545c64",
    text-color="#fff",active-text-color="#ffd04b",:router="true", v-if="showMobileNav && screenWidth <= 960")
      template(v-for="item in showItems")
        template(v-if="item.meta.issub === false")
          el-menu-item(:index="item.children[0].name", :route="{name:item.children[0].name}", class="submenu-title-noDropdown", :key="item.children[0].name")
            svg-icon(:name="item.children[0].meta.icon")
            span {{item.children[0].meta.title}}
        template(v-else)
          el-submenu(:index="item.meta.title", :key="item.meta.title")
            template(slot="title")
              svg-icon(:name="item.meta.icon")
              span {{item.meta.title}}
            template(v-for="children in item.children")
              el-menu-item(:index="children.name", :route="{name:children.name}", :key="children.name") {{children.meta.title}}
      el-menu-item(class="submenu-title-noDropdown",@click.native="jumpToFront")
        svg-icon(name="logout")
        span 返回前台
</template>

<script>
import routes from "@/web-admin/js/routes";
export default {
  name: "topbar",
  data() {
    return {
      showMobileNav: false
    };
  },
  props: {
    screenWidth: {
      type: Number
    }
  },
  methods: {
    toggleMobileNav() {
      this.showMobileNav = !this.showMobileNav;
    },
    jumpToFront() {
      window.location.href = "/";
    }
  },
  computed: {
    showItems() {
      // 过滤一级导航
      return routes.filter(function(route, index) {
        if (route.meta.hidden === false && route.children) {
          // 过滤二级导航
          route.children = route.children.filter(function(item, index) {
            return item.meta.hidden === false;
          });
          return route.meta.hidden === false && route.children;
        }
      });
    },
    defaultActive() {
      return this.$route.name;
    }
  }
};
</script>

<style lang="scss" scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.5s ease;
  max-height: 1000px;
}
.slide-fade-enter, .slide-fade-leave-to
/* .slide-fade-leave-active for below version 2.1.8 */ {
  max-height: 0px;
  opacity: 1;
}
.topbar__mobile__nav__menu {
  overflow: hidden;
  position: absolute !important;
  top: 80px;
  width: 100%;
  z-index: 1000;
}
.topbar {
  position: relative;
  height: 80px;
  background: $plblue;
  .topbar__title {
    position: absolute;
    left: 20px;
    top: 0;
    height: 100%;
    line-height: 80px;
    text-align: center;
    color: white;
    font-size: 36px;
  }
  .topbar__userinfo {
    line-height: 75px;
    font-size: 20px;
    padding-right: 20px;
  }
  .topbar__mobile_nav {
    position: relative;
    height: 100%;
    .mobile-humber {
      width: 100px;
      float: right;
      height: 100%;
      box-sizing: border-box;
      padding: 15px;
      a {
        height: 100%;
        width: 80%;
        margin: 0 auto;
        display: block;
        position: relative;
        transition: all 0.3s;
        .line {
          position: absolute;
          width: 100%;
          margin: 0 auto;
          height: 2px;
          display: block;
          background: $c15;
          transition: all 0.3s;
          &:nth-child(1) {
            top: 0;
          }
          &:nth-child(2) {
            top: 50%;
          }
          &:nth-child(3) {
            top: 100%;
          }
        }
      }
      .active {
        .line {
          position: absolute;
          width: 100%;
          margin: 0 auto;
          height: 2px;
          display: block;
          background: $c15;
          transition: all 0.3s;
          &:nth-child(1) {
            top: 0;
            transform: translateY(25px) rotate(45deg);
          }
          &:nth-child(2) {
            display: none;
          }
          &:nth-child(3) {
            top: 100%;
            transform: translateY(-25px) rotate(-45deg);
          }
        }
      }
    }
  }
}
</style>
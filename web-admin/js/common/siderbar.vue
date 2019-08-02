<template lang="pug">
  .siderbar
    el-menu(:default-active="defaultActive",class="el-menu-vertical-demo siderbar__ul",background-color="#545c64",
    text-color="#fff",active-text-color="#ffd04b",:router="true")
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
  name: "siderbar",
  data() {
    return {};
  },
  methods: {
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
.siderbar__ul {
  text-align: left;
}
.siderbar {
  position: fixed;
  top: 80px;
  left: 0px;
  height: 100%;
  width: 248px;
  background: #535b64;
  border-right: solid 2px #e6e6e6;
  .el-menu {
    border-right: none !important;
  }
}
</style>

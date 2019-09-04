<template lang="pug">
.topbar__wrapper
  .topbar__title(v-if="screenWidth > 960")
    router-link(:to="{name:'index'}") AHPUOJ
  .topbar__mobile_nav(v-else)
    .mobile-humber(@click="toggleMobileNav")
      a(href="#", :class="{active:showMobileNav}")
        span.line
        span.line
        span.line
  .topbar__nav(v-if="screenWidth > 960")
    ul.topbar__nav__bar.clearfix
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'problemSet'}")
          svg-icon(name="problem") 
          span 问题集
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'issueList'}")
          svg-icon(name="chat") 
          span 讨论区
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'status'}",@click.native="resetSolutionFilter")
          svg-icon(name="server") 
          span 评测机
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'contestList'}")
          svg-icon(name="champion") 
          span 竞赛
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'seriesList'}")
          svg-icon(name="tournament") 
          span 系列赛
      li.topbar__nav__item.topbar__section
        router-link(:to="{name:'ranklist'}")
          svg-icon(name="ranking") 
          span 排名
  .topbar__right__section
    .topbar__userinfo__wrapper(v-if="$store.getters.token && $store.getters.userDefunct == 0",@mouseleave="showDropDownMenu=false")
      .username-wrapper.ell.tar
        span.username  {{$store.getters.userNick}}
      img(:src="imgUrl($store.getters.userAvatar)",@mouseover="showDropDownMenu = true")
      el-collapse-transition
        ul.topbar__userinfo__dropdown(v-if="showDropDownMenu")
          li 
            router-link(:to="{name:'userinfo',params:{id:$store.getters.userId}}") 个人空间
          li 
            router-link(:to="{name:'account'}") 账号设置
          li
            router-link(:to="{name:'myreplys'}") 查看回复
          li(v-if="$store.getters.userRole=='admin'") 
            a(href="/admin") 后台管理
          li(@click="handleLogout") 
            a 登出

    .topbar__login__wrapper(v-else, @click="handleLogin")
        a  
          svg-icon(name="login") 
          span 登录

  transition(name="slide-fade")
      el-menu(@select="toggleMobileNav",class="topbar__mobile__nav__menu",background-color="#111144",
      text-color="#fff",active-text-color="#ffd04b",:router="true", v-if="showMobileNav && screenWidth <= 960")
        el-menu-item(:index="1", :route="{name:'index'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="dashboard")
          span 首页
        el-menu-item(:index="1", :route="{name:'problemSet'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="problem")
          span 问题集
        el-menu-item(:index="1", :route="{name:'issueList'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="chat")
          span 讨论区
        el-menu-item(:index="1", :route="{name:'status'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="server")
          span 评测机
        el-menu-item(:index="1", :route="{name:'contestList'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="champion")
          span 竞赛
        el-menu-item(:index="1", :route="{name:'seriesList'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="tournament")
          span 系列赛
        el-menu-item(:index="1", :route="{name:'ranklist'}", class="submenu-title-noDropdown")
          svg-icon.m__svg(name="ranking")
          span 排名

  el-dialog(:visible.sync="dialogFormVisible",width="5rem",:close-on-click-modal="false",custom-class="auth__dialog__wrapper")
    .auth__dialog
      .title(slot="header") 
        span(:class="[method=='login'?'active':'']",@click="method='login'") &nbsp登 录&nbsp
        span(:class="[method=='register'?'active':'']", @click="method='register'") &nbsp注 册&nbsp
      el-form(v-if="method=='login'", ref="loginForm", :model="loginForm", :rules="loginRules", key="login")
        el-form-item(prop="username")
          el-input(v-model="loginForm.username", placeholder="请输入用户名")
            svg-icon(slot="prefix", name="user", class="auth__input__prefix")
        el-form-item(prop="password")
          el-input(v-model="loginForm.password",type="password", placeholder="请输入密码")
              svg-icon(slot="prefix", name="password", class="auth__input__prefix")
        el-form-item
          el-button(type="primary",style="width:100%;", @click="submitLogin") 登录
        router-link(:to="{name:'findpass'}",target="_blank") 忘记密码？点击找回
      el-form(v-if="method=='register'", ref="registerForm", :model="registerForm", :rules="registerRules", key="register")
        el-form-item(prop="email")
          el-input(v-model="registerForm.email", placeholder="请输入邮箱")
            svg-icon(slot="prefix", name="email", class="auth__input__prefix")
        el-form-item(prop="username")
          el-input(v-model="registerForm.username", placeholder="请输入用户名")
            svg-icon(slot="prefix", name="user", class="auth__input__prefix")
        el-form-item(prop="nick")
          el-input(v-model="registerForm.nick", placeholder="请输入昵称")
            svg-icon(slot="prefix", name="user", class="auth__input__prefix")
        el-form-item(prop="password")
          el-input(v-model="registerForm.password",type="password", placeholder="请输入密码")
            svg-icon(slot="prefix", name="password", class="auth__input__prefix")
        el-form-item(prop="confirmpassword")
          el-input(v-model="registerForm.confirmpassword",type="password", placeholder="请确认密码")
            svg-icon(slot="prefix", name="password", class="auth__input__prefix")
        el-form-item
          el-button(type="primary",style="width:100%;",@click="submitRegister") 注册
</template>

<script>
import { login, register } from "@/web-user/js/api/auth.js";
import { EventBus } from "@/web-common/eventbus";
export default {
  data() {
    var validatePassword = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请输入密码"));
      } else {
        if (this.registerForm.confirmpassword !== "") {
          this.$refs.registerForm.validateField("confirmpassword");
        }
        callback();
      }
    };
    var validateConfirmPassword = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请再次输入密码"));
      } else if (value !== this.registerForm.password) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    };
    return {
      showMobileNav: false,
      showDropDownMenu: false,
      dialogFormVisible: false,
      method: "login",
      loginForm: {
        username: "",
        password: ""
      },
      registerForm: {
        email: "",
        username: "",
        nick: "",
        password: "",
        confirmpassword: ""
      },
      loginRules: {
        username: [
          {
            required: true,
            message: "请输入用户名",
            trigger: "blur"
          }
        ],
        password: [
          {
            required: true,
            message: "请输入密码",
            trigger: "blur"
          }
        ]
      },
      registerRules: {
        email: [
          {
            required: true,
            message: "请输入邮箱地址",
            trigger: "blur"
          },
          {
            type: "email",
            message: "请输入正确的邮箱地址",
            trigger: "blur"
          }
        ],
        username: [
          {
            required: true,
            message: "请输入用户名",
            trigger: "blur"
          },
          {
            max: 20,
            message: "用户名不能超过20个字符",
            trigger: "blur"
          },
          {
            pattern: /^[a-zA-Z0-9]+$/,
            message: "用户名只能包含英文和数字字符",
            trigger: "blur"
          }
        ],
        nick: [
          {
            required: true,
            message: "请输入用户昵称",
            trigger: "blur"
          },
          {
            max: 20,
            message: "昵称不能超过20个字符",
            trigger: "blur"
          }
        ],
        password: [
          {
            validator: validatePassword,
            trigger: "blur"
          },
          {
            min: 6,
            message: "密码不能少于6个字符",
            trigger: "blur"
          },
          {
            max: 20,
            message: "密码不能超过20个字符",
            trigger: "blur"
          }
        ],
        confirmpassword: [
          {
            validator: validateConfirmPassword,
            trigger: "blur"
          },
          {
            min: 6,
            message: "密码不能少于6个字符",
            trigger: "blur"
          },
          {
            max: 20,
            message: "密码不能超过20个字符",
            trigger: "blur"
          }
        ]
      },
      showMobileNav: false
    };
  },
  props: {
    screenWidth: {
      type: Number
    }
  },
  mounted() {
    console.log(this.$store);
    EventBus.$on("goLogin", () => {
      this.handleLogin();
    });
  },
  methods: {
    toggleMobileNav() {
      this.showMobileNav = !this.showMobileNav;
    },
    handleLogin() {
      this.dialogFormVisible = true;
    },
    handleLogout() {
      this.showDropDownMenu = false;
      setTimeout(() => {
        this.$router.push({ name: "index" });
        this.$store.dispatch(`Logout`);
        this.$message({
          message: "登出成功",
          type: "success"
        });
      }, 500);
    },
    submitLogin() {
      const self = this;
      self.$refs["loginForm"].validate(async valid => {
        if (valid) {
          try {
            let res = await self.$store.dispatch("Login", self.loginForm);
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.dialogFormVisible = false;
            self.showDropDownMenu = false;
            await self.$store.dispatch("GetUserInfo");
            self.$router.replace({ name: "refresh" });
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
        } else {
          return false;
        }
      });
    },
    async submitRegister() {
      const self = this;
      self.$refs.registerForm.validate(async valid => {
        if (valid) {
          try {
            let res = await self.$store.dispatch("Register", self.registerForm);
            console.log(res);
            self.$message({
              message: res.data.message,
              type: "success"
            });
            self.dialogFormVisible = false;
            self.showDropDownMenu = false;
            await self.$store.dispatch("GetUserInfo");
          } catch (err) {
            console.log(err);
            self.$message({
              message: err.response.data.message,
              type: "error"
            });
          }
        } else {
          return false;
        }
      });
    },
    resetSolutionFilter() {
      this.$store.dispatch("resetSolutionFilter");
    }
  }
};
</script>

<style lang="scss" scoped>
$mibile-nav-height: 50px;

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.5s ease;
  max-height: 10rem;
}
.slide-fade-enter,
.slide-fade-leave-to {
  max-height: 0rem;
  opacity: 1;
}

.topbar__mobile__nav__menu {
  text-align: left;
  overflow: hidden;
  position: absolute !important;
  top: $mibile-nav-height;
  width: 100%;
  z-index: 1000;
}

.topbar__wrapper {
  position: relative;
  background: $pddblue;
  border-bottom: 1px solid $pdblue;
  height: 100px;
  @media screen and (max-width: 960px) {
    height: $mibile-nav-height;
  }
  .topbar__title {
    text-indent: 10px;
    a {
      cursor: pointer;
    }
    height: 100%;
    line-height: 100px;
    float: left;
    width: 150px;
    font-size: 30px;
    @media screen and (min-width: 1280px) {
      width: 200px;
      font-size: 40px;
    }
    text-align: center;
    a {
      cursor: pointer;
      color: $pblue;
    }
  }

  .topbar__nav {
    color: $c9;
    position: relative;
    margin-left: 150px;
    @media screen and (min-width: 1280px) {
      margin-left: 200px !important;
      .topbar__section {
        a {
          font-size: 22px !important;
        }
        svg {
          height: 30px !important;
        }
      }
    }
    line-height: 100px;
    .topbar__nav__bar {
      display: flex;
      justify-content: space-around;
      margin-right: 200px;
      & > .topbar__nav__item {
        flex: 0 1 auto;
      }
      li.topbar__nav__item {
        display: inline-block;
        height: 100%;
        position: relative;
        cursor: pointer;
      }
    }
    .topbar__section {
      width: 200px;
      a {
        font-size: 18px;
        display: block;
        color: $c9;
      }
      svg {
        height: 24px;
        margin-right: 0.05rem;
      }
      display: block;
      height: 100%;
      cursor: pointer;
      transition: all 0.2s ease-out;
      &:hover {
        color: $pblue;
        a {
          color: $pblue;
        }
        &::before {
          background: $pblue;
          transform: translate3d(0, 0, 0) scaleX(0.8);
        }
      }
      &::before {
        content: "";
        left: 0rem;
        bottom: 0.1rem;
        position: absolute;
        background: $pblue;
        height: 0.02rem;
        width: 100%;
        transition: all 0.2s ease;
        transform: translate3d(0, 0, 0) scaleX(0);
      }
      &.active {
        color: $pblue;
        &::before {
          background: $pblue;
          transform: translate3d(0, 0, 0) scaleX(0.8);
        }
      }
    }
  }
  .topbar__right__section {
    line-height: 100px;
    height: 100%;
    font-size: 0.22rem;
    width: 200px;
    position: absolute;
    right: 0.4rem;
    top: 0;
    a {
      color: $c9;
      svg {
        height: 30px !important;
      }
    }
    .topbar__userinfo__wrapper {
      position: relative;
      .username-wrapper {
        box-sizing: border-box;
        float: left;
        padding: 0 20px;
        height: 100px;
        width: 140px;
        line-height: 100px;
        color: $c12;
        font-size: 16px;
      }
      img {
        cursor: pointer;
        margin: 20px 0 0 0;
        height: 60px;
        border-radius: 30px;
      }
      .topbar__userinfo__dropdown {
        position: absolute;
        width: 150px;
        top: 100px;
        right: -35px;
        z-index: 100;
        background: $c15;
        border: 1px solid $c12;
        font-size: 0.18rem;
        transition: all 0.3s;
        li {
          cursor: pointer;
          color: $c3;
          height: 50px;
          line-height: 50px;
          border-bottom: 0.01rem solid $c13;
          transition: all 0.3s;
          a {
            font-size: 16px;
          }
          &:hover {
            color: $pblue;
          }
        }
      }
    }
    .topbar__login__wrapper {
      margin-right: -20px;
    }
    @media screen and (max-width: 960px) {
      height: $mibile-nav-height;
      line-height: $mibile-nav-height;
      .topbar__userinfo__wrapper {
        .username-wrapper {
          height: $mibile-nav-height;
          line-height: $mibile-nav-height;
        }
        img {
          margin: 5px 0 0 0;
          height: 40px;
          border-radius: 20px;
        }
        .topbar__userinfo__dropdown {
          top: $mibile-nav-height;
        }
      }
    }
  }
}

.auth__dialog {
  border-radius: 10px;
  .title {
    text-align: left;
    margin-bottom: 0.1rem;
    padding-left: 0.1rem;
    color: $pblue;
    font-size: 24px;
    span {
      cursor: pointer;
      color: $c10;
      border-bottom: 4px solid $c12;
    }
    span.active {
      color: $pblue;
      border-bottom: 4px solid $pblue;
    }
  }
  .auth__input__prefix {
    height: 100%;
    line-height: 100%;
    width: 20px;
  }
}

.topbar__mobile_nav {
  position: relative;
  height: 100%;
  float: left;
  .mobile-humber {
    margin-left: 10px;
    width: 40px;
    float: left;
    height: 100%;
    box-sizing: border-box;
    padding: 10px 0;
    a {
      height: 100%;
      width: 100%;
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
          transform: translateY(15px) rotate(45deg);
        }
        &:nth-child(2) {
          display: none;
        }
        &:nth-child(3) {
          top: 100%;
          transform: translateY(-15px) rotate(-45deg);
        }
      }
    }
  }
}
</style>
import Vue from 'vue/dist/vue.esm.js'
import routes from '@/web-user/js/routes'
import store from '@/web-user/js/store'
import { Message } from 'element-ui'
import Cookies from "js-cookie"
import VueRouter from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

Vue.use(VueRouter)

const router = new VueRouter({
    mode: 'history',
    base: '/',
    routes
})

router.beforeEach(async (to, from, next) => {
    console.log("npstart")
    NProgress.start() // 进度条开始
    if (Cookies.get("access-token")) { // 如果有token
        console.log("have token")
        if (store.getters.username.length === 0) {
            try {
                let res = await store.dispatch('GetUserInfo') // 拉取用户信息
                let data = res.data.user;
                console.log("get user info")
                next()
            } catch (err) {
                console.log(err)
                store.dispatch('Logout').then(() => {
                    Message.error('登录超时,请重新登陆')
                })
                next()
            }
        } else {
            next()
        }
    } else {
        next() // 没有登录 不影响正常浏览
    }
})

router.afterEach(() => {
    NProgress.done() // 进度条结束
})

export default router
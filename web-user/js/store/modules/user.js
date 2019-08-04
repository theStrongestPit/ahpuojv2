import {
    login,
    register
} from '@/web-user/js/api/auth'
import {
    getUser
} from '@/web-user/js/api/user'
import Cookies from "js-cookie"

const user = {
    state: {
        id: 0,
        username: '',
        nick: '',
        role: '',
        avatar: '',
        submit: 0,
        solved: 0,
        defunct: 1, // defunct = 1表示被封禁的状态
        token: Cookies.get("access-token")
    },
    mutations: {
        SET_TOKEN: (state, token) => {
            state.token = token
        },
        SET_USER: (state, user) => {
            state.id = user.id
            state.username = user.username
            state.nick = user.nick
            state.role = user.role
            state.avatar = user.avatar
            state.submit = user.submit
            state.solved = user.solved
            state.defunct = user.defunct
        },
        SET_AVATAR: (state, path) => {
            state.avatar = path
        },
    },
    actions: {
        Login({
            commit
        }, loginForm) {
            return new Promise(async (resolve, reject) => {
                try {
                    let res = await login(loginForm)
                    let token = Cookies.get("access-token");
                    let role = 'user';
                    commit('SET_TOKEN', token)
                    resolve(res)
                } catch (err) {
                    reject(err)
                }
            })
        },
        Register({
            commit
        }, registerForm) {
            return new Promise(async (resolve, reject) => {
                try {
                    let res = await register(registerForm)
                    let token = Cookies.get("access-token");
                    let role = 'user';
                    commit('SET_TOKEN', token)
                    resolve(res)
                } catch (err) {
                    reject(err)
                }
            })
        },
        // 刷新token
        RefreshToken({
            commit
        }, token) {
            commit('SET_TOKEN', token);
        },
        // 获取用户姓名权限头像信息
        GetUserInfo({
            commit,
            state
        }) {
            return new Promise(async (resolve, reject) => {
                try {
                    let res = await getUser()
                    const user = res.data.user
                    commit('SET_USER', user)
                    resolve(res)
                } catch (err) {
                    reject(err)
                }
            })
        },
        // 更新用户信息
        UpdateUserInfo({
            commit
        }, user) {
            commit('SET_USER', user)
        },
        // 更新用户头像
        UpdateUserAvatar({
            commit
        }, path) {
            commit('SET_AVATAR', path)
        },
        // 登出
        Logout({
            commit,
            state
        }) {
            commit('SET_TOKEN', '')
            commit('SET_USER', {
                username: "",
                nick: "",
                role: "",
                avatar: "",
                submit: 0,
                solved: 0,
                defunct: 1,
            })
            Cookies.remove("access-token")
        }
    }
}

export default user
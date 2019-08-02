const layout = () => import('@/web-user/js/pages/Layout.vue')

export default [{
    path: '/',
    component: layout,
    children: [{
            path: '',
            name: 'index',
            meta: {
                keepAlive: true,
                title: '首页 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/home'),
        },
        {
            path: 'findpass',
            name: 'findpass',
            meta: {
                keepAlive: true,
                title: '找回密码 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/findpass'),
        },
        {
            path: 'resetpass',
            name: 'resetpass',
            meta: {
                keepAlive: true,
                title: '重设密码 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/resetpass'),
        },
        {
            path: 'problemset',
            name: 'problemSet',
            meta: {
                keepAlive: true,
                title: '问题集 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/problem_set'),
        },
        {
            path: 'issues',
            name: 'issueList',
            meta: {
                keepAlive: false,
                title: '讨论区 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/issue_list'),
        },
        {
            path: 'issue/:id',
            name: 'issue',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/issue'),
        },
        {
            path: 'problem/:id/issues',
            name: 'problemIssueList',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/issue_list'),
        },
        {
            path: 'status',
            name: 'status',
            meta: {
                keepAlive: true,
                title: '评测机 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/status'),
        },
        {
            path: 'contest/:id/status',
            name: 'contestStatus',
            meta: {
                keepAlive: true,
                title: '评测机 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/status'),
        },
        {
            path: 'contest/:id/rank',
            name: 'contestRank',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/contest_rank'),
        },
        {
            path: 'contest/:id/teamrank',
            name: 'contestTeamRank',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/contest_team_rank'),
        },
        {
            path: 'problem/:id',
            name: 'problem',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/problem'),
        },
        {
            path: 'contest/:id/problem/:num',
            name: 'contestProblem',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/problem'),
        },
        {
            path: 'contests',
            name: 'contestList',
            meta: {
                keepAlive: true,
                title: '竞赛&作业 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/contest_list'),
        },
        {
            path: 'ranklist',
            name: 'ranklist',
            meta: {
                keepAlive: true,
                title: '排名 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/ranklist'),
        },
        {
            path: 'contest/:id',
            name: 'contest',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/contest'),
        },
        {
            path: 'solution/:id',
            name: 'solution',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/solution'),
        },
        {
            path: 'account',
            name: 'account',
            meta: {
                keepAlive: false,
                title: '账号设置 - AHPUOJ',
            },
            component: () => import('@/web-user/js/pages/account_setting'),
        },
        {
            path: 'userinfo/:id',
            name: 'userinfo',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/userinfo'),
        },
        // hack方法 只刷新路由
        {
            path: 'refresh',
            name: 'refresh',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/refresh'),
        },

        // 404路由
        {
            path: '404',
            name: '404Page',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/404'),
        },
        {
            path: '*',
            name: '404',
            meta: {
                keepAlive: false,
            },
            component: () => import('@/web-user/js/pages/404'),
        },
    ]
}, ]
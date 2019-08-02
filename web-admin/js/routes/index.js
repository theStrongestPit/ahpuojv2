const layout = () => import('@/web-admin/js/pages/Layout.vue')


export default [{
        path: '/admin',
        meta: {
            issub: false,
            hidden: false,
        },
        component: layout,
        children: [{
            path: '',
            name: 'home',
            meta: {
                icon: 'dashboard',
                title: '首页',
                hidden: false,
            },
            component: () => import('@/web-admin/js/pages/home')
        }]
    },
    {
        path: '/admin/new',
        component: layout,
        meta: {
            icon: 'file-text',
            title: "新闻管理",
            issub: true,
            hidden: false,
        },
        children: [{
                path: 'add',
                name: 'adminAddNew',
                meta: {
                    title: '添加新闻',
                    keepAlive: false,
                    hidden: false,
                },
                component: () => import('@/web-admin/js/pages/new'),
            },
            {
                path: ':id/edit',
                name: 'adminEditNew',
                meta: {
                    title: '编辑新闻',
                    keepAlive: false,
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/new'),
            },
            {
                path: 'list',
                name: 'adminNewlist',
                meta: {
                    title: '新闻列表',
                    hidden: false,
                    keepAlive: true,
                },
                component: () => import('@/web-admin/js/pages/new_list'),
            },
        ]
    },
    {
        path: '/admin/user',
        component: layout,
        meta: {
            icon: 'user',
            title: "用户管理",
            issub: true,
            hidden: false,
        },
        children: [{
            path: 'list',
            name: 'adminUserList',
            meta: {
                title: '用户列表',
                hidden: false,
                keepAlive: true,
            },
            component: () => import('@/web-admin/js/pages/user_list'),
        }, ]
    },
    {
        path: '/admin/team',
        component: layout,
        meta: {
            issub: false,
            hidden: false,
        },
        children: [{
                path: 'list',
                name: 'adminTeamList',
                meta: {
                    icon: "team",
                    title: '团队列表',
                    hidden: false,
                    keepAlive: true,
                },
                component: () => import('@/web-admin/js/pages/team_list'),
            },
            {
                path: ':id/manage',
                name: 'adminTeamManage',
                meta: {
                    title: '团队管理',
                    keepAlive: false,
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/team_manage'),
            },
        ]
    },
    {
        path: '/admin/problem',
        component: layout,
        meta: {
            icon: 'problem',
            title: "问题管理",
            issub: true,
            hidden: false,
        },
        children: [{
                path: 'add',
                name: 'adminAddproblem',
                meta: {
                    title: '添加问题',
                    keepAlive: false,
                    hidden: false,
                },
                component: () => import('@/web-admin/js/pages/problem'),
            },
            {
                path: ':id/edit',
                name: 'adminEditProblem',
                meta: {
                    title: '编辑问题',
                    keepAlive: false,
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/problem'),
            },
            {
                path: ':id/data',
                name: 'adminProblemData',
                meta: {
                    title: '编辑问题数据',
                    keepAlive: false,
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/problem_data'),
            },
            {
                path: 'import',
                name: 'adminImportProblem',
                meta: {
                    title: '导入问题',
                    hidden: false,
                },
                component: () => import('@/web-admin/js/pages/import_problem'),
            },
            {
                path: 'rejudge',
                name: 'adminRejudgeProblem',
                meta: {
                    title: '重判问题',
                    hidden: false,
                },
                component: () => import('@/web-admin/js/pages/rejudge_problem'),
            },
            {
                path: 'reassign',
                name: 'adminReassignProblem',
                meta: {
                    title: '重排问题',
                    hidden: false,
                },
                component: () => import('@/web-admin/js/pages/reassign_problem'),
            },
            {
                path: 'list',
                name: 'adminProblemList',
                meta: {
                    title: '问题列表',
                    hidden: false,
                    keepAlive: true,
                },
                component: () => import('@/web-admin/js/pages/problem_list'),
            },
        ]
    },
    {
        path: '/admin/tag',
        component: layout,
        meta: {
            hidden: false,
            issub: false,
        },
        children: [{
            path: '',
            name: 'adminTaglist',
            meta: {
                icon: 'tag',
                title: '标签列表',
                hidden: false,
                keepAlive: true,
            },
            component: () => import('@/web-admin/js/pages/tag_list'),
        }, ]
    },
    {
        path: '/admin/contest',
        component: layout,
        meta: {
            icon: 'champion',
            title: "竞赛&作业管理",
            issub: true,
            hidden: false,
        },
        children: [{
                path: 'add',
                name: 'adminAddContest',
                meta: {
                    title: '添加竞赛&作业',
                    hidden: false,
                    keepAlive: false,
                },
                component: () => import('@/web-admin/js/pages/contest'),
            },
            {
                path: ':id/edit',
                name: 'adminEditContest',
                meta: {
                    title: '编辑竞赛&作业',
                    hidden: true,
                    keepAlive: false,
                },
                component: () => import('@/web-admin/js/pages/contest'),
            },
            {
                path: ':id/manage',
                name: 'adminContestManage',
                meta: {
                    title: '竞赛&作业人员管理',
                    hidden: true,
                    keepAlive: false,
                },
                component: () => import('@/web-admin/js/pages/contest_manage'),
            },
            {
                path: ':id/teammanage',
                name: 'adminContestTeamManage',
                meta: {
                    title: '竞赛&作业(团队赛)人员管理',
                    hidden: true,
                    keepAlive: false,
                },
                component: () => import('@/web-admin/js/pages/contest_team_manage'),
            },
            {
                path: 'list',
                name: 'adminContestList',
                meta: {
                    title: '竞赛&作业列表',
                    hidden: false,
                    keepAlive: false,
                },
                component: () => import('@/web-admin/js/pages/contest_list'),
            },
            {
                path: 'series',
                name: 'adminSeriesList',
                meta: {
                    title: '系列赛列表',
                    hidden: false,
                    keepAlive: true,
                },
                component: () => import('@/web-admin/js/pages/series_list'),
            },
            {
                path: ':id/seriermanage',
                name: 'adminSeriesManage',
                meta: {
                    title: '系列赛管理',
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/series_manage'),
            },
        ]
    },
    {
        path: '/admin/generator',
        component: layout,
        meta: {
            issub: false,
            hidden: false,
        },
        children: [{
            path: '',
            name: 'adminAccountGenerator',
            meta: {
                icon: "add_user",
                title: "账号生成器",
                hidden: false,
            },
            component: () => import('@/web-admin/js/pages/account_generator'),
        }]
    },
    // 404路由
    {
        path: '/admin/*',
        component: layout,
        meta: {
            hidden: true,
        },
        children: [{
                path: '/admin/404',
                name: 'admin404Page',
                meta: {
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/404'),
            },
            {
                path: '',
                name: 'admin404',
                meta: {
                    hidden: true,
                },
                component: () => import('@/web-admin/js/pages/404'),
            },
        ]
    },
]
import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            redirect: '/projects'
        },
        {
            path: '/',
            component: resolve => require(['../components/page/home/Index.vue'], resolve),
            meta: { title: '首页' },
            children:[
                {
                    path: '/projects',
                    component: resolve => require(['../components/page/project/Index.vue'], resolve),
                    meta: { title: '项目列表' }
                },
                {
                    path: '/dbs/:pid',
                    component: resolve => require(['../components/page/db/Index.vue'], resolve),
                    meta: { title: '数据库列表' },
                    props: true
                },
                {
                    path: '/dashboard',
                    component: resolve => require(['../components/page/Dashboard.vue'], resolve),
                    meta: { title: '系统首页' }
                },
                {
                    path: '/404',
                    component: resolve => require(['../components/page/404.vue'], resolve),
                    meta: { title: '404' }
                },
                {
                    path: '/403',
                    component: resolve => require(['../components/page/403.vue'], resolve),
                    meta: { title: '403' }
                }
            ]
        },
        {
            path: '/login',
            component: resolve => require(['../components/page/Login.vue'], resolve)
        },
        {
            path: '*',
            redirect: '/404'
        }
    ]
})

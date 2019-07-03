<template>
    <el-container>
        <el-header style="border-bottom:solid 1px #e6e6e6;">
            <el-container style="display:block;">
                <h1 style="line-height:60px;float:left;">LOGO</h1>
                <el-menu :default-active="activeIndex" style="float:right;" mode="horizontal" @select="handleSelect">
                    <el-menu-item>
                        <el-input style="width:260px;" placeholder="请输入内容"></el-input>
                    </el-menu-item>
                    <el-menu-item index="1">项目管理</el-menu-item>
                    <el-menu-item index="2">工作台</el-menu-item>
                    <el-menu-item index="3" disabled>消息中心</el-menu-item>
                    <el-submenu index="4">
                        <template slot="title">我的账户</template>
                        <el-menu-item index="4-1">个人信息</el-menu-item>
                        <el-menu-item index="4-2">修改密码</el-menu-item>
                        <el-menu-item index="4-3">退出登录</el-menu-item>
                    </el-submenu>
                </el-menu>
            </el-container>
        </el-header>
        <el-main style="min-height:600px;max-height:800px;">
            <transition name="move" mode="out-in">
                <router-view></router-view>
            </transition>
        </el-main>
        <el-footer>Footer</el-footer>
    </el-container>
</template>

<script>
    import bus from '../common/bus';
    export default {
        data() {
            return {
                activeIndex:'1'
            }
        },
        methods:{
            handleSelect() {
            }
        }
    }
</script>

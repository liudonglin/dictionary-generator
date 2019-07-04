<template>
    <el-container style="height: 100%;">
        <el-header style="border-bottom:solid 1px #e6e6e6;">
            <el-container style="display:block;">
                <h1 style="line-height:60px;float:left;">LOGO</h1>
                <el-menu :default-active="activeIndex" style="float:right;" mode="horizontal" @select="handleSelect">
                    <el-menu-item>
                        <el-input style="width:260px;" placeholder="输入后回车键查询" v-model="search_word" @keyup.enter.native="searchEnterFun"></el-input>
                    </el-menu-item>
                    <el-menu-item index="1">项目管理</el-menu-item>
                    <el-menu-item index="2">工作台</el-menu-item>
                    <el-menu-item index="3">用户管理</el-menu-item>
                    <el-submenu index="4">
                        <template slot="title">我的账户</template>
                        <el-menu-item index="4-1">个人信息</el-menu-item>
                        <el-menu-item index="4-2">修改密码</el-menu-item>
                        <el-menu-item index="4-3">退出登录</el-menu-item>
                    </el-submenu>
                </el-menu>
            </el-container>
        </el-header>
        <el-main style="">
            <transition name="move" mode="out-in">
                <router-view></router-view>
            </transition>
        </el-main>
        <el-footer>
            <div style="padding: 20px;">
                <el-link href="https://github.com/liudonglin/dictionary-generator" target="_blank">GitHub</el-link>
            </div>
        </el-footer>
    </el-container>
</template>

<script>
    import bus from '../../common/bus';
    export default {
        data() {
            return {
                activeIndex:'1',
                search_word:''
            }
        },
        methods:{
            handleSelect() {
            },
            searchEnterFun(e) {
                let keyCode = window.event? e.keyCode:e.which;
                let val = this.search_word;
                if(keyCode == 13) {
                     bus.$emit('searchWordChange', val);
                 }
            }
        }
    }
</script>

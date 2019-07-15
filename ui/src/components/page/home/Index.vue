<template>
    <el-container style="height: 100%;">
        <el-header style="border-bottom:solid 1px #e6e6e6;">
            <el-container style="display:block;">
                <h1 style="line-height:60px;float:left;">LOGO</h1>
                <el-menu :default-active="onRoutes" style="float:right;" mode="horizontal" @select="handleSelect" unique-opened router>
                    <el-menu-item>
                        <el-input style="width:260px;" placeholder="输入后回车键查询" v-model="search_word" @keyup.enter.native="searchEnterFun"></el-input>
                    </el-menu-item>
                    <el-menu-item index="/projects">项目管理</el-menu-item>
                    <el-menu-item index="/code">代码生成</el-menu-item>
                    <el-menu-item index="/usermgt">用户管理</el-menu-item>
                    <el-submenu index="/account">
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
            <div style="height:100%;">
                
                <el-link :underline="false" style="margin-left:5px;margin-right:30px;float:right;margin-top:4px;" href="https://hub.docker.com/r/liudonglin/dictionary-generator" target="_blank">DockerHub</el-link>
                <img class="logo-img" src="../../../assets/img/docker-hub-logo.png"/>
                
                <el-link style="margin-left:5px;margin-right:30px;float:right;margin-top:4px;" :underline="false"  href="https://github.com/liudonglin/dictionary-generator" target="_blank">GitHub</el-link>
                <svg xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet" viewBox="2 2 20 20" style="width:25px;height:25px;float:right;"><path d="M12 2.4c-5.412 0-9.8 4.407-9.8 9.843 0 4.35 2.808 8.038 6.702 9.34.49.09.67-.214.67-.475 0-.233-.01-.853-.014-1.674-2.726.595-3.3-1.32-3.3-1.32-.447-1.137-1.09-1.44-1.09-1.44-.89-.61.068-.598.068-.598.983.07 1.5 1.016 1.5 1.016.875 1.504 2.295 1.068 2.853.816.088-.635.34-1.068.62-1.314-2.175-.25-4.463-1.094-4.463-4.865 0-1.075.382-1.955 1.01-2.642-.102-.25-.438-1.25.095-2.605 0 0 .823-.265 2.695 1.01.782-.22 1.62-.33 2.453-.333.833.004 1.67.113 2.453.332 1.87-1.273 2.693-1.01 2.693-1.01.534 1.355.197 2.357.098 2.606.628.688 1.007 1.567 1.007 2.64 0 3.782-2.29 4.614-4.473 4.858.352.303.664.904.664 1.822 0 1.315-.01 2.377-.01 2.7 0 .263.176.57.673.474 3.89-1.305 6.696-4.99 6.696-9.34C21.8 6.808 17.41 2.4 12 2.4z"></path></svg>
            
                <el-link style="margin-right:30px;float:right;margin-top:4px;" :underline="false"  href="https://element.eleme.cn/#/zh-CN" target="_blank">Element UI</el-link>
                <svg style="width:30px;height:30px;float:right;" viewBox="0 0 38 48" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
                    <g id="v2.2.0" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
                        <g id="首页-默认效果-copy-2" transform="translate(-70.000000, -19.000000)" fill="#409EFF">
                            <path d="M212.135441,45.1578077 Z M103.416502,46.2175511 C103.407008,47.7945682 102.587841,48.1463474 102.587841,48.1463474 C102.587841,48.1463474 88.4520478,56.3145087 87.5329956,56.8262751 C86.622336,57.2171561 86.0136703,56.8262751 86.0136703,56.8262751 C86.0136703,56.8262751 71.2225706,48.2479572 70.6826962,47.8708444 C70.1426843,47.4937316 70.1301643,46.9063774 70.1301643,46.9063774 C70.1301643,46.9063774 70.1450232,29.9199174 70.1301643,29.1333364 C70.1153053,28.3466177 71.0969576,27.7555461 71.0969576,27.7555461 L85.8755373,19.2134387 C86.7853714,18.7332016 87.6711286,19.2134387 87.6711286,19.2134387 C87.6711286,19.2134387 100.726623,26.8020937 102.173442,27.6177257 C103.591507,28.2919577 103.416502,29.6843424 103.416502,29.6843424 C103.416502,29.6843424 103.425307,44.7519192 103.416502,46.2175511 L103.416502,46.2175511 Z M97.5160105,29.1691339 C94.4900173,27.4274483 87.3615842,23.2971063 87.3615842,23.2971063 C87.3615842,23.2971063 86.6661039,22.9212326 85.9517747,23.2971063 L74.3487406,29.98215 C74.3487406,29.98215 73.5780025,30.4449014 73.589697,31.0604805 C73.6013915,31.6760597 73.589697,44.9698708 73.589697,44.9698708 C73.589697,44.9698708 73.5994654,45.4295932 74.0234951,45.7246472 C74.4473873,46.0197011 86.0601898,52.7332451 86.0601898,52.7332451 C86.0601898,52.7332451 86.538152,53.0391759 87.2530315,52.7332451 C87.9746526,52.3327262 99.0730335,45.9402582 99.0730335,45.9402582 C99.0730335,45.9402582 99.7162325,45.6648929 99.7236619,44.4307056 C99.7257257,44.0747959 99.7265512,42.6932881 99.7266887,40.9575229 L86.6600502,48.8752394 L86.6600502,45.846221 C86.6600502,44.6021205 87.6231289,43.7809811 87.6231289,43.7809811 L99.1807607,36.8165793 C99.6168978,36.3611251 99.7068768,35.6315447 99.7254505,35.3556287 C99.7250378,34.0904627 99.7244874,32.9852841 99.7240747,32.2851678 L86.6600502,40.2012321 L86.6600502,37.034531 C86.6600502,35.7904305 87.4855462,35.2446564 87.4855462,35.2446564 L97.5160105,29.1691339 Z" id="Shape-Copy"></path>
                        </g>
                    </g>
                </svg>
            
            </div>
        </el-footer>
    </el-container>
</template>

<script>
    import bus from '../../common/bus';
    export default {
        data() {
            return {
                search_word:''
            }
        },
        computed:{
            onRoutes() {
                let path = this.$route.path
                if (path.indexOf('/dbs/')==0) {
                    path = "/projects"
                }
                return path
            }
        },
        methods:{
            handleSelect() {
            },
            searchEnterFun(e) {
                let keyCode = window.event? e.keyCode:e.which;
                let val = this.search_word;
                if(keyCode == 13) {
                     bus.$emit('headerSearchWordChange', val);
                 }
            },
        }
    }
</script>

<style>
.logo-img{
    border-radius:50%;
    width:26px;
    height:26px;
    float:right;
}
</style>


<template>
    <div style="height:100%;">
        <el-container style="height:100%;">
            <el-aside width="300px"  style="border-right:solid 1px #e6e6e6;">
                <v-sidebar></v-sidebar>
            </el-aside>
            <el-main>
                <v-codePage></v-codePage>
            </el-main>
        </el-container> 
    </div>
</template>

<script>
    import vSidebar from './TableAside.vue';
    import vCodePage from './CodePage.vue';
    export default {
        data(){
            return {
            }
        },
        components:{
            vSidebar,vCodePage
        },
    }
</script>
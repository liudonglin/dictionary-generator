<template>
    <div style="height:100%;">
        <el-container style="height:100%;">
            <el-aside width="300px"  style="border-right:solid 1px #e6e6e6;">
                <v-db-sidebar :pid="pid"></v-db-sidebar>
            </el-aside>
            <el-main>
                <v-table-list :pid="pid"></v-table-list>
            </el-main>
        </el-container> 
    </div>
</template>

<script>
    import vDbSidebar from './DBAside.vue';
    import vTableList from './TableList.vue';
    import bus from '../../common/bus';
    export default {
        props: {
            pid: {
                type: String,
                default: '0'
            }
        },
        data(){
            return {}
        },
        components:{
            vDbSidebar,vTableList
        },
        created() {
        },
        methods: {
        },
        watch: {    
            '$route' (to, from) {
            }
        }
    }
</script>

        
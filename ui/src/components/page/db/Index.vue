<template>
    <div style="height:100%;">
        <el-container style="height:100%;">
            <el-aside width="300px"  style="border-right:solid 1px #e6e6e6;">
                <v-db-sidebar :pid="pid" :pname="projectInfo.name" ></v-db-sidebar>
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
            return {
                loadProjectUrl: '/api/project/load',
                projectInfo:{},
            }
        },
        components:{
            vDbSidebar,vTableList
        },
        created() {
            this.loadProject();
        },
        methods: {
            loadProject() {
                let pid = parseInt(this.pid)
                this.$axios.post(this.loadProjectUrl, { id:pid }).then(result=>{
                    if (result.success) {
                        this.projectInfo = result.data
                    }
                })
            },
        },
    }
</script>

        
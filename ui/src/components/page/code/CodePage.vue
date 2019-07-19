<template>
<div style="height: 100%;" v-loading="codeLoading">
    <el-tabs type="card" v-model="selectTempleteID" @tab-click="handleTempleteChange">
        <el-tab-pane v-for="templete in codeTempletes" :key="templete.name" :label="templete.name" :name="templete.id+''"></el-tab-pane>
    </el-tabs>
    <mavon-editor v-model="content" :toolbarsFlag="false" :shortCut="false" :subfield="false" defaultOpen="edit" placeholder="请选择数据表以生成代码"/>
</div>
</template>

<script>
import bus from '../../common/bus';
import { mavonEditor } from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'
export default {
    data(){
        return {
            loadTempleteUrl:"/api/templete/load",
            listTempleteUrl:"/api/templete/list",
            codeLoading:false,
            codeTempletes:[],
            tableId:0,
            projectId:0,
            selectTempleteID:"0",
            content:"",
        }
    },
    components: {
        mavonEditor
    },
    created() {
        bus.$on('codePageSelectTableChange', val => {
            this.tableId = parseInt(val)
            this.getCodeContent()
        })

        bus.$on('codePageSelectProjectChange', val => {
            this.projectId = parseInt(val)
            this.getTempletes()
            this.content = ""
            this.tableId = 0
        })
    },
    methods: {
        getTempletes() {
            this.codeLoading = true
            this.$axios.post(this.listTempleteUrl, {pid:this.projectId, index:0, size:999999, order_by:"templete_created DESC" }).then(result=>{
                if (result.success) {
                    this.codeTempletes = result.data.list
                    if(this.codeTempletes.length>0){
                        this.selectTempleteID = this.codeTempletes[0].id+""
                    }
                }
                this.codeLoading = false
            })
        },
        getCodeContent() {
            if (this.tableId==0){
                return
            }
            this.codeLoading = true
            let templeteId = parseInt(this.selectTempleteID)
            this.$axios.post(this.loadTempleteUrl, { templete_id:templeteId,tid:this.tableId }).then(result=>{
                if (result.success) {
                    this.content = result.data
                }
                this.codeLoading = false
            })
        },
        handleTempleteChange(tab) {
            this.getCodeContent()
        },
    }
}
</script>
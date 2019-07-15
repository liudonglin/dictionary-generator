<template>
<div style="height: 100%;">
    <el-tabs type="card" v-model="selectTempleteID">
        <el-tab-pane v-for="templete in codeTempletes" :key="templete.name" :label="templete.name" :name="templete.id+''"></el-tab-pane>
    </el-tabs>
    <mavon-editor v-model="content" :editable="false" :toolbarsFlag="false" :subfield="false" defaultOpen="preview"/>
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
            codeTempletes:[],
            tableId:0,
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

        this.getTempletes()
    },
    methods: {
        getTempletes() {
            this.$axios.post(this.listTempleteUrl, {index:0, size:999999, order_by:"templete_created DESC" }).then(result=>{
                if (result.success) {
                    this.codeTempletes = result.data.list
                    if(this.codeTempletes.length>0){
                        this.selectTempleteID = this.codeTempletes[0].id+""
                    }
                }
            })
        },
        getCodeContent() {
            let templeteId = parseInt(this.selectTempleteID)
            this.$axios.post(this.loadTempleteUrl, { templete_id:templeteId,tid:this.tableId }).then(result=>{
                if (result.success) {
                    this.content = result.data
                }
            })
        }
    }
}
</script>
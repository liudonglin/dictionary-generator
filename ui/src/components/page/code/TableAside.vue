<template>
<div>
    <ul style="padding-right:20px;overflow-x:hidden;">

        <li style="list-style:none;">
            <div class="db-box">
                <el-select v-model="selectProjectId" placeholder="请选择" @change="handleProjectChange">
                    <el-option v-for="item in projects" :key="item.id" :label="item.name" :value="item.id">
                    </el-option>
                </el-select>
            </div>
        </li>

        <el-collapse accordion @change="handleDBCollapseChange">
            <el-collapse-item v-for="db in dbs" :key="db.id" :title="db.name" :name="db.id">
                <div class="db-item" v-for="table in tables" :key="table.id" :class="selectTableId==table.id?'active':''" >
                    <span @click="handleTableSelectChange(table)">{{table.name}}</span>
                </div>
            </el-collapse-item>
        </el-collapse>
        
    </ul>

</div>
</template>

<script>
import bus from '../../common/bus';
export default {
    data() {
        return {
            listProjectUrl: '/api/project/list',
            listDBtUrl: '/api/database/list',
            listTabletUrl: '/api/table/list',
            projects:[{name:"请选择项目",id:0}],
            selectProjectId:0,
            dbs:[],
            selectDBId:0,
            cacheTable: new Map(),
            tables:[],
            selectTableId:0,
        }
    },
     created() {
        this.searchProjects();
    },
    methods: {
        searchProjects() {
            this.$axios.post(this.listProjectUrl, { index:0, size:999999, order_by:"project_created DESC" }).then(result=>{
                if (result.success) {
                    this.projects = result.data.list
                }
            })
        },
        handleProjectChange() {
            if (this.selectProjectId==0){
                return
            }
            this.$axios.post(this.listDBtUrl, {pid:this.selectProjectId, index:0, size:999999, order_by:"database_created DESC" }).then(result=>{
                if (result.success) {
                    this.dbs = result.data.list
                }
            })
            this.selectDBId = 0
            this.selectTableId = 0
            bus.$emit('codePageSelectProjectChange', this.selectProjectId);
        },
        handleDBCollapseChange(activeName) {
            let did = parseInt(activeName)
            let tables = this.cacheTable.get(did)

            if(tables!=null) {
                this.tables = tables
                return
            }

            this.$axios.post(this.listTabletUrl, {did:did, index:0, size:999999, order_by:"table_created DESC" }).then(result=>{
                if (result.success) {
                    this.tables = result.data.list
                    this.cacheTable.set(did,result.data.list)
                }
            })
        },
        handleTableSelectChange(table) {
            this.selectDBId = table.did
            this.selectTableId = table.id
            bus.$emit('codePageSelectTableChange', this.selectTableId);
        },
    }
}
</script>

<style>
.db-box {
    width: 100%;
    height: 60px;
    padding-top:10px; 
}

.db-item {
    color: #909399;
    font-size: 14px;
    cursor: pointer;
    padding-left: 10px;
    padding-top:10px;
    padding-bottom:10px;
}

.active {
    border-bottom: 2px solid #409EFF;
    color:#000000;
}

.db-item:hover {
    color: #303133;
    transition:color 0.3s linear, right .3s ease;
}
</style>

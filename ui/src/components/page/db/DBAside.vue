<template>
<div>
    <ul style="padding-right:20px;overflow-x:hidden;">

        <li style="list-style:none;">
            <div class="db-box">
                <el-button type="primary" icon="el-icon-plus" circle class="db-btn" title="新增"></el-button>
                <el-button type="primary" icon="el-icon-connection" circle class="db-btn" @click="openConnForm" title="数据库倒入"></el-button>
                <el-button type="primary" icon="el-icon-upload" circle class="db-btn"  title="excel导入"></el-button>
            </div>
        </li>

        <li style="list-style:none;" v-for="db in dbDatas" :key="db.id">
            <div class="db-box db-item">{{db.name}}
                <div class="db-item-del">删除</div>
            </div>
        </li>
        
    </ul>

    <!-- 连接弹出框 -->
    <el-dialog title="连接" :visible.sync="editConnVisible" width="50%" @close="closeConnForm('connForm')">
        
        <el-table :data="connection_list" :class="dbInfos.length>0?'hide':''">
            <el-table-column property="name" label="名称">
                <template slot-scope="scope">
                    <el-input v-model="scope.row.name" disabled></el-input>
                </template>
            </el-table-column>
            <el-table-column property="host" label="域名">
                <template slot-scope="scope">
                    <el-input v-model="scope.row.host" disabled></el-input>
                </template>
            </el-table-column>
            <el-table-column align="right" width="80">
                <template slot-scope="scope">
                    <el-button size="mini" type="primary" @click="setDefaultConn(scope.row.id)" :disabled="scope.row.is_default">默认</el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-tree :data="dbInfos" show-checkbox :props="dbTreeProps" v-loading="treeLoading" @check-change="handleTreeCheckChange">
            <span class="custom-tree-node" slot-scope="{ node,data }">
                <span>{{ node.label }}</span>
                <span>{{ data.comment }}</span>
            </span>
        </el-tree>
        
        <span slot="footer" class="dialog-footer">
            <el-button @click="closeConnForm">取 消</el-button>
            <el-button @click="loadConn" type="primary" :class="dbInfos.length>0?'hide':''">连 接</el-button>
            <el-button @click="saveConn" type="primary" :class="dbInfos.length==0?'hide':''">确 定</el-button>
        </span>
    </el-dialog>
</div>
</template>

<script>
    import bus from '../../common/bus';
import { debuglog } from 'util';
    const default_conn = "default_project_conn_"
    export default {
        props: {
            pid: {
                type: String,
                default: '0'
            }
        },
        data() {
            return {
                listConnsUrl:'/api/conn/loadpid',
                loadConnUrl:'/api/dbimport/loaddb',
                saveConnUrl:'/api/dbimport/savedbs',
                listDBUrl:'/api/database/list',
                search_word:'',
                dbDatas:[],
                connection_list:[],
                editConnVisible: false,
                treeLoading: false,
                dbInfos:[], 
                dbSelectInfo:new Map(),
                dbTreeProps:{
                    children: 'tables',
                    label: 'name'
                },
            }
        },
        created() {
            this.search();
        },
        methods: {
            search() {
                this.$axios.post(this.listDBUrl, { name:this.search_word, pid:this.pid, index:0, size:999999, order_by:"database_created DESC" }).then(result=>{
                    if (result.success) {
                        this.dbDatas = result.data.list
                    }
                })
            },
            openConnForm() {
                this.$axios.post(this.listConnsUrl, this.pid).then(result=>{
                    if (result.success) {
                        if (result.data==null||result.data.length==0){
                            return
                        }

                        let def = localStorage.getItem(default_conn + this.pid);
                        result.data.forEach((item, index) => {
                            if (def==item.id) {
                                item.is_default = true
                            } else {
                                item.is_default = false
                            }
                        });

                        this.connection_list = result.data
                    }
                })

                this.editConnVisible=true
            },
            setDefaultConn(cid){
                localStorage.setItem(default_conn + this.pid,cid)
                this.connection_list.forEach((item, index) => {
                    if (cid==item.id) {
                        item.is_default = true
                    } else {
                        item.is_default = false
                    }
                });
            },
            closeConnForm() {
                this.editConnVisible = false
                this.dbInfos=[];
                this.dbSelectInfo.clear();
            },
            loadConn() {
                let conn = null
                this.connection_list.forEach((item, index) => {
                    if (item.is_default) {
                        conn = item
                    } 
                });

                if (conn==null) {
                    this.$message('请设置默认链接');
                    return
                }

                this.treeLoading = true
                
                this.$axios.post(this.loadConnUrl, conn).then(result=>{
                    if (result.success) {
                        this.dbInfos = result.data
                    }
                    this.treeLoading = false
                })
            },
            saveConn() {
                if (this.dbSelectInfo.size==0){
                    return
                }

                this.treeLoading = true

                let topDBs = new Map()
                let _pid = this.pid
                this.dbSelectInfo.forEach(function(value,key) {
                    topDBs.set(value.db_name,{ name: value.db_name, pid: _pid, tables: [] })
                })
                this.dbSelectInfo.forEach(function(value,key){
                    let topDB = topDBs.get(value.db_name)
                    topDB.tables.push(value)
                })

                let postDB = []
                topDBs.forEach(function(value,key){
                    postDB.push(value)
                })

                this.$axios.post(this.saveConnUrl, postDB).then(result=>{
                    if (result.success) {
                        this.search()
                        this.closeConnForm()
                    }
                    this.treeLoading = false
                })
                
            },
            handleTreeCheckChange(data, checked, indeterminate) {
                let key = ''

                if (data.tables==null) {
                    // 选择数据表

                    key = data.db_name+'_'+data.name
                    if (checked) {
                        this.dbSelectInfo.set(key,data)
                    } else {
                        this.dbSelectInfo.delete(key)
                    }

                } else {
                    // 选择数据库
                    
                    data.tables.forEach((item, index) => {
                        key = item.db_name+'_'+item.name
                        if (checked) {
                            this.dbSelectInfo.set(key,item)
                        } else {
                            this.dbSelectInfo.delete(key)
                        }
                    });
                }
            },
        }
    }
</script>

<style scoped>

.db-btn {
    margin-left: 20px;
    margin-bottom: 20px;
}

.db-box {
    width: 100%;
    height: 60px;
    border-bottom:solid 1px #e6e6e6;
    line-height: 60px;
    
}

.db-item {
    color: #909399;
    font-size: 14px;
    cursor: pointer;
    padding-left: 10px;
    transition:color 0.3s linear;
    z-index:70;
}

.db-item-del {
    display: block;
    width:90px;
    height: 60px;
    font-size: 14px;
    line-height: 60px;
    cursor: pointer;
    z-index:80;
    right: -290px;
    position: relative;
    top: -60px;
    color: #fff;
    background-color: #F56C6C;
    transition:all 0.3s linear;
    overflow:hidden;
    box-sizing: border-box;
    padding-left: 20px;
}

.db-item:hover {
    color: #303133;
    transition:color 0.3s linear, right .3s ease;
}

.db-item:hover .db-item-del{
    right: -200px;
}

.hide {
    display: none;
}

</style>
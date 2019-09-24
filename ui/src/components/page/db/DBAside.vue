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
            <div class="db-box">
                <el-button type="primary" icon="el-icon-plus" circle class="db-btn" @click="editDBVisible=true" title="新增"></el-button>
                <el-button type="primary" icon="el-icon-connection" circle class="db-btn" @click="openConnForm" title="数据库倒入"></el-button>
                <el-button type="primary" icon="el-icon-download" circle class="db-btn" @click="downloadProject" title="excel导出"></el-button>
            </div>
        </li>

        <li style="list-style:none;" v-for="db in dbDatas" :key="db.id">
            <div class="db-box db-item" :class="selectDBId==db.id?'active':''">
                <span @click="handleDBSelectChange(db.id)">{{db.name}}</span>
                <div class="db-item-action">
                    <div class="edit">
                        <i class="el-icon-edit" @click="editDBForm(db)">编辑</i>
                    </div>
                    <div class="del">
                        <i class="el-icon-delete" @click="deleteDB(db)">删除</i>
                    </div>
                </div>
            </div>
        </li>
        
    </ul>

    <!-- 编辑DB弹出框 -->
        <el-dialog title="编辑" :visible.sync="editDBVisible" width="40%" @close="closeDBForm('form')">
            <el-form ref="form" :model="form" :rules="rules" label-width="100px">
                <el-form-item label="项目名称:">
                    <el-input v-model="selectProjectName" maxlength="40" show-word-limit :disabled="true"></el-input>
                </el-form-item>
                <el-form-item label="数据库名称:" prop="name">
                    <el-input v-model="form.name" maxlength="40" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="标 题:">
                   <el-input v-model="form.title" maxlength="40" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="描述信息:">
                    <el-input type="textarea" placeholder="请输入内容" v-model="form.description"
                        maxlength="200" show-word-limit :autosize="{ minRows: 4, maxRows: 8}">
                    </el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="closeDBForm('form')">取 消</el-button>
                <el-button @click="saveDB('form')" type="primary">确 定</el-button>
            </span>
        </el-dialog>

    <!-- 连接弹出框 -->
    <el-dialog title="连接" :visible.sync="importConnVisible" width="50%" @close="closeConnForm('connForm')">
        
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
    import { Base64 } from 'js-base64';
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
                listProjectUrl: '/api/project/list',
                listConnsUrl:'/api/conn/loadpid',
                loadConnUrl:'/api/dbimport/loaddb',
                saveConnUrl:'/api/dbimport/savedbs',
                listDBUrl:'/api/database/list',
                saveDBUrl: '/api/database/save',
                deleteDBUrl:'/api/database/delete',
                downloadProjectUrl:'/api/export/project',
                search_word:'',
                dbDatas:[],
                selectDBId:0,
                connection_list:[],
                importConnVisible: false,
                editDBVisible: false,
                treeLoading: false,
                dbInfos:[], 

                selectProjectId:0,
                selectProjectName:"",
                projects:[],

                dbTreeSelectInfo:new Map(),
                dbTreeProps:{
                    children: 'tables',
                    label: 'name'
                },
                form: {
                    id: 0,
                    pid: parseInt(this.pid),
                    name: '',
                    title: '',
                    description: ''
                },
                rules: {
                    name: [
                        { required: true, message: '请输入名称', trigger: 'blur' }
                    ],
                    db_name: [
                        { required: true, message: '请输入数据库名称', trigger: 'blur' }
                    ],
                    title:[
                        { required: true, message: '请填写描述信息', trigger: 'blur' }
                    ]
                }
            }
        },
        created() {
            this.searchDBs();
            this.searchProjects();
        },
        methods: {
            searchProjects() {
                this.$axios.post(this.listProjectUrl, { index:0, size:999999, order_by:"project_created DESC" }).then(result=>{
                    if (result.success) {
                        this.projects = result.data.list
                        this.selectProjectId = this.form.pid;
                        let selectProjectId = this.selectProjectId
                        this.selectProjectName = this.projects.find(function(item){
                            return item.id == selectProjectId;
                        }).name
                    }
                })
            },
            searchDBs() {
                let pid = parseInt(this.pid)
                this.$axios.post(this.listDBUrl, { name:this.search_word, pid:pid, index:0, size:999999, order_by:"database_created DESC" }).then(result=>{
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

                        if(def==null) {
                            def = result.data[0].id
                            localStorage.setItem(default_conn + this.pid,def)
                        }

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

                this.importConnVisible=true
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
                this.importConnVisible = false
                this.dbInfos=[];
                this.dbTreeSelectInfo.clear();
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
                if (this.dbTreeSelectInfo.size==0){
                    return
                }

                this.treeLoading = true

                let topDBs = new Map()
                let _pid = parseInt(this.pid)
                this.dbTreeSelectInfo.forEach(function(value,key) {
                    topDBs.set(value.db_name,{ name: value.db_name, pid: _pid, tables: [] })
                })
                this.dbTreeSelectInfo.forEach(function(value,key){
                    let topDB = topDBs.get(value.db_name)
                    topDB.tables.push(value)
                })

                let postDB = []
                topDBs.forEach(function(value,key){
                    postDB.push(value)
                })

                this.$axios.post(this.saveConnUrl, postDB).then(result=>{
                    if (result.success) {
                        this.searchDBs()
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
                        this.dbTreeSelectInfo.set(key,data)
                    } else {
                        this.dbTreeSelectInfo.delete(key)
                    }

                } else {
                    // 选择数据库
                    
                    data.tables.forEach((item, index) => {
                        key = item.db_name+'_'+item.name
                        if (checked) {
                            this.dbTreeSelectInfo.set(key,item)
                        } else {
                            this.dbTreeSelectInfo.delete(key)
                        }
                    });
                }
            },
            closeDBForm(formName) {
                this.editDBVisible = false;
                this.form.id=0;
                this.form.name="";
                this.form.title="";
                this.form.description="";
                this.$refs[formName].resetFields();
            },
            saveDB(formName) {
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.$axios.post(this.saveDBUrl, this.form).then(result=>{
                            if (result.success) {
                                this.form.id=result.data
                                this.searchDBs()
                                this.closeDBForm(formName)
                            }
                        })
                    }
                });
            },
            editDBForm(db){
                this.form.id=db.id;
                this.form.name=db.name;
                this.form.title=db.title;
                this.form.description=db.description;
                this.editDBVisible = true;
            },
            deleteDB(db) {
                this.$confirm(`确定要删除数据库 : ${db.name}`, '提示信息', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.$axios.post(this.deleteDBUrl, db.id).then(result=>{
                        if (result.success) {
                            // 遍历删除页面上对应列的数据
                            let index =0;
                            this.dbDatas.forEach(function(_db,i) {
                                if (db.id==_db.id) {
                                    index = i;
                                }
                            })
                            this.dbDatas.splice(index, 1);
                            bus.$emit('dbMgtSelectDBChange', this.selectDBId);
                        }
                    })
                })
            },
            handleProjectChange() {
                let selectProjectId = this.selectProjectId
                this.selectProjectName = this.projects.find(function(item){
                    return item.id == selectProjectId;
                }).name
                this.$router.push('/dbs/'+selectProjectId);
            },
            handleDBSelectChange(dbid) {
                this.selectDBId = dbid;
                bus.$emit('dbMgtSelectDBChange', this.selectDBId);
            },
            downloadProject(){
                let pid = parseInt(this.pid)
                this.$axios.post(this.downloadProjectUrl,pid,{responseType: 'blob'}).then(result=>{
                    if (!result.data) {
                        return
                    }

                    let filename = result.headers['content-disposition'].split(";")[1].split("filename=")[1];
                    if (!filename) {
                        filename = "excel.xlsx";
                    }
                    filename = Base64.decode(filename)

                    if (typeof window.navigator.msSaveBlob !== 'undefined') {
                        // IE workaround for "HTML7007: One or more blob URLs were 
                        // revoked by closing the blob for which they were created. 
                        // These URLs will no longer resolve as the data backing 
                        // the URL has been freed."
                        window.navigator.msSaveBlob(blob, filename);
                        return
                    }

                    let url = window.URL.createObjectURL(result.data)
                    let link = document.createElement('a')
                    link.style.display = 'none'
                    link.href = url
                    link.setAttribute('download', filename)

                    document.body.appendChild(link)
                    link.click()
                    document.body.removeChild(link);
                    window.URL.revokeObjectURL(url);
                })
            }
        },
        watch: {    
            '$route' (to, from) {   
                this.searchDBs()
            }
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

.db-item.active {
    border-bottom: 2px solid #409EFF;
    color:#000000;
}

.db-item-action {
    display: block;
    width:140px;
    height: 60px;
    font-size: 14px;
    line-height: 60px;
    z-index:80;
    right: -290px;
    position: relative;
    top: -60px;
    color: #fff;
    transition:all 0.3s linear;
    overflow:hidden;
    box-sizing: border-box;
}

.db-item:hover {
    color: #303133;
    transition:color 0.3s linear, right .3s ease;
}

.db-item:hover .db-item-action{
    right: -150px;
}

.db-item-action .edit{
    width:50%;
    height: 100%;
    display: block;
    background-color: #409EFF;
    float:left;
    text-align:center;
}

.db-item-action .del{
    width:50%;
    height: 100%;
    display: block;
    background-color: #F56C6C;
    float:right;
    text-align:center;
}

.hide {
    display: none;
}

</style>
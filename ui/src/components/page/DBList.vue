<template>
    <div>
        <div class="container">
            <div class="handle-box">
                <el-row class="mb20">
                    <el-form :inline="true" >
                        <el-col :span="12">
                        <el-form-item label="项目名称:">
                            <el-input v-model="project.name" :disabled="true"></el-input>
                        </el-form-item>
                        </el-col>
                        <el-form-item label="数据库:">
                            <el-input v-model="project.data_base" :disabled="true"></el-input>
                        </el-form-item>
                    </el-form>
                </el-row>
                <el-input v-model="search_word" placeholder="数据库名称" class="handle-input mr10"></el-input>
                <el-button icon="el-icon-search" circle @click="search" title="查询"></el-button>
                <el-button type="primary" icon="el-icon-plus" circle @click="editDBVisible=true" title="新增"></el-button>
                <el-button type="primary" icon="el-icon-connection" circle @click="editConnVisible=true" title="数据库倒入"></el-button>
                <el-button type="primary" icon="el-icon-upload" circle @click="uploadDB" title="excel导入"></el-button>
            </div>

            <el-collapse @change="handleDBCollapseChange" accordion v-loading="collapseLoading">
                <el-collapse-item v-for="db in dbData" :key="db.id" :name="db.id">
                    <template slot="title">
                        <el-col :span="8">
                            <span>数据库 : {{db.name}}</span>
                        </el-col>
                        <el-col :span="2" :offset="14">
                            <el-button type="primary" icon="el-icon-edit" circle title="新增"></el-button>
                            <el-button type="danger" icon="el-icon-delete" circle></el-button>
                        </el-col>
                    </template>

                    <div v-for="table in db.tables" :key="table.id" class="mb20">
                        <el-form :inline="true">
                            <el-form-item label="表 名:">
                                <el-input v-model="table.name" ></el-input>
                            </el-form-item>
                            <el-form-item label="描 述:">
                                <el-input v-model="table.description"></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" >保存</el-button>
                            </el-form-item>
                        </el-form>

                        <el-table :data="table.columns">
                            <el-table-column prop="name" label="列名" sortable >
                            </el-table-column>
                            <el-table-column prop="title" label="描述" sortable >
                            </el-table-column>
                            <el-table-column prop="data_type" label="数据类型" sortable >
                            </el-table-column>
                            <el-table-column prop="name" label="项目名称" sortable >
                            </el-table-column>
                            <el-table-column label="操作" width="200" align="center">
                                <template>
                                    <el-button type="text" icon="el-icon-edit" >编辑</el-button>
                                    <el-button type="text" icon="el-icon-delete" class="red" >删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-collapse-item>
            </el-collapse>
            <div class="pagination">
                <el-pagination background @current-change="search" layout="prev, pager, next" 
                :total="pageTotal" :page-size="pageSize" :current-page.sync="pageCurrent" :hide-on-single-page="true" >
                </el-pagination>
            </div>

        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editDBVisible" width="40%" @close="closeDBForm('form')">
            <el-form ref="form" :model="form" :rules="rules" label-width="100px">
                <el-form-item label="项目名称:">
                    <el-input v-model="project.name" maxlength="20" show-word-limit :disabled="true"></el-input>
                </el-form-item>
                <el-form-item label="数据库名称:" prop="name">
                    <el-input v-model="form.name" maxlength="20" show-word-limit></el-input>
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
        <el-dialog title="连接" :visible.sync="editConnVisible" width="50%" @close="closeConnForm('connForm')">
            <el-form ref="connForm" :model="connForm" :rules="rules" label-width="100px" :class="connInfo.length>0?'hide':''">
                <el-form-item label="数据库类型:">
                    <el-input v-model="project.data_base" :disabled="true"></el-input>
                </el-form-item>
                <el-form-item label="HostPort:" prop="host_port">
                    <el-input v-model="connForm.host_port"></el-input>
                </el-form-item>
                <el-form-item label="登录账户:" prop="user">
                    <el-input v-model="connForm.user"></el-input>
                </el-form-item>
                <el-form-item label="登录密码:" prop="password">
                    <el-input v-model="connForm.password" show-password></el-input>
                </el-form-item>
            </el-form>

            <el-tree :data="connInfo" show-checkbox :props="dbTreeProps" v-loading="treeLoading" @check-change="handleTreeCheckChange">
                <span class="custom-tree-node" slot-scope="{ node,data }">
                    <span>{{ node.label }}</span>
                    <span>{{ data.comment }}</span>
                </span>
            </el-tree>
            
            <span slot="footer" class="dialog-footer">
                <el-button @click="loadConn('connForm')" type="primary" :class="connInfo.length>0?'hide':''">连 接</el-button>
                <el-button @click="closeConnForm('connForm')" :class="connInfo.length==0?'hide':''">取 消</el-button>
                <el-button @click="saveConn" type="primary" :class="connInfo.length==0?'hide':''">确 定</el-button>
            </span>
        </el-dialog>

    </div>
</template>

<script>
    import bus from '../common/bus';
    export default {
        props: {
            pid: {
                type: String,
                default: '0'
            }
        },
        data() {
            return {
                loadConnUrl:'/api/conn/loaddb',
                saveConnUrl:'/api/conn/savedbs',
                loadProjectUrl: '/api//project/load',
                saveDBUrl: '/api/database/save',
                listDBUrl:'/api/database/list',
                loadDBUrl:'/api/database/load',
                loading: false,
                treeLoading:false,
                collapseLoading:false,
                pageTotal:0,
                pageSize:10,
                pageCurrent:1,
                search_word: '',
                editDBVisible: false,
                editConnVisible: false,
                project:{},
                name: '',
                dbData:[],
                connInfo:[], 
                connSelectInfo:new Map(),
                dbTreeProps:{
                    children: 'tables',
                    label: 'name'
                },
                form: {
                    id: 0,
                    pid: parseInt(this.pid),
                    name: '',
                    description: ''
                },
                connForm: {
                    host_port:'',
                    user:'',
                    password:''
                },
                rules: {
                    name: [
                        { required: true, message: '请输入项目名称', trigger: 'blur' }
                    ],
                    db_name: [
                        { required: true, message: '请输入数据库名称', trigger: 'blur' }
                    ],
                    host_port: [
                        { required: true, message: '请输入host:port', trigger: 'blur' }
                    ],
                    user: [
                        { required: true, message: '请输入登录账户', trigger: 'blur' }
                    ],
                    password: [
                        { required: true, message: '请输入登录密码', trigger: 'blur' }
                    ],
                }
            }
        },
        created() {
            this.search();
            this.loadProject();
        },
        computed:{
            
        },
        methods: {
            loadProject(){
                this.$axios.post(this.loadProjectUrl, { id:this.form.pid }).then(result=>{
                    if (result.success) {
                        this.project = result.data
                    }
                })
            },
            search() {
                this.loading=true
                this.$axios.post(this.listDBUrl, { name:this.search_word, pid:this.form.pid, index:this.pageCurrent-1, size:this.pageSize, order_by:"database_created DESC" }).then(result=>{
                    if (result.success) {
                        this.dbData = result.data.list
                        this.pageTotal = result.data.total
                    }
                    this.loading=false
                })
            },
            closeDBForm(formName) {
                this.editDBVisible = false;
                this.form.id=0;
                this.form.name="";
                this.form.description="";
                this.$refs[formName].resetFields();
            },
            saveDB(formName) {
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.$axios.post(this.saveDBUrl, this.form).then(result=>{
                            if (result.success) {
                                this.form.id=result.data
                                this.search()
                            }
                        })
                    }
                });
            },
            closeConnForm(formName){
                this.editConnVisible = false
                this.connForm.host_port="";
                this.connForm.user="";
                this.connForm.password="";
                this.connInfo=[];
                this.connSelectInfo.clear();
                this.$refs[formName].resetFields();
            },
            loadConn(formName){
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.treeLoading = true
                        this.connForm.data_base = this.project.data_base
                        this.$axios.post(this.loadConnUrl, this.connForm).then(result=>{
                            if (result.success) {
                                this.connInfo = result.data
                            }
                            this.treeLoading = false
                        })
                    }
                });
            },
            handleDBCollapseChange(val){
                if (val!='') {
                    let _db ={}
                    this.dbData.forEach(function(value,index){
                        if (value.id==val){
                            _db = value;
                        }
                    })

                    if (_db.tables==null) {
                        this.collapseLoading = true
                        this.$axios.post(this.loadDBUrl, {id:_db.id}).then(result=>{
                            if (result.success) {
                                _db.tables = result.data.tables
                            }
                            this.collapseLoading = false
                        })
                    }
                }
            },
            handleTreeCheckChange(data, checked, indeterminate) {
                let key = ''
                let _checked = false

                if (data.tables!=null) {
                    key = data.name
                    _checked = indeterminate
                } else {
                    key = data.db_name+'_'+data.name
                    _checked = checked
                }
                
                if (_checked) {
                    this.connSelectInfo.set(key,data)
                } else {
                    this.connSelectInfo.delete(key)
                }
            },
            saveConn(){
                if (this.connSelectInfo.size==0){
                    return
                }

                this.treeLoading = true

                let topDBs = new Map()
                let _pid = this.form.pid
                this.connSelectInfo.forEach(function(value,key){
                    if (value.tables!=null) {
                        topDBs.set(key,{ name: value.name, pid: _pid, tables: [] })
                    }
                })
                this.connSelectInfo.forEach(function(value,key){
                    if (value.tables==null) {
                        let topDB = topDBs.get(value.db_name)
                        topDB.tables.push(value)
                    }
                })

                let postDB = []
                topDBs.forEach(function(value,key){
                    postDB.push(value)
                })

                this.$axios.post(this.saveConnUrl, postDB).then(result=>{
                    if (result.success) {
                        this.search()
                    }
                    this.treeLoading = false
                })
            },
            uploadDB(){

            }
        }
    }
</script>

<style scoped>
    .table{
        width: 100%;
        font-size: 14px;
    }

    .hide{
        display: none;
    }

    .mr30{
        margin-right: 30px;
    }

    .ml30{
        margin-left: 30px;
    }

    .mb20{
        margin-bottom: 20px;
    }

    .handle-box {
        margin-bottom: 20px;
    }
    .handle-input {
        width: 300px;
        display: inline-block;
    }
    .mr10{
        margin-right: 10px;
    }
    .red{
        color: #ff0000;
    }

    .table-expand {
        font-size: 0;
    }

    .table-expand label {
        width: 90px;
        color: #99a9bf;
    }

    .table-expand .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 50%;
    }
</style>
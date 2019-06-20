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
                <!-- <el-button type="primary" icon="el-icon-upload" circle @click="uploadDB" title="excel导入"></el-button> -->
            </div>

            <el-collapse @change="handleDBCollapseChange" accordion v-loading="collapseLoading">
                <el-collapse-item v-for="db in dbData" :key="db.id" :name="db.id">
                    <template slot="title">
                        <el-col :span="8">
                            <span>数据库 : {{db.name}}</span>
                        </el-col>
                        <el-col :span="2" :offset="14">
                            <el-button @click="editDBForm(db)" type="primary" icon="el-icon-edit" circle title="编辑数据库"></el-button>
                            <el-button @click="deleteDB(db)" type="danger" icon="el-icon-delete" circle title="删除数据库"></el-button>
                        </el-col>
                    </template>

                    <div v-for="table in db.tables" :key="table.id" class="mb20">
                        <el-form :inline="true" :model="table" :rules="rules" ref="tableform">
                            <el-form-item label="表 名:" prop="name">
                                <el-input v-model="table.name" maxlength="20" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="描 述:">
                                <el-input v-model="table.description" maxlength="20" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" icon="el-icon-check" @click="saveTable(table)" circle title="保存"></el-button>
                                <el-button type="danger" icon="el-icon-delete" @click="deleteTable(table)" circle title="删除"></el-button>
                            </el-form-item>
                        </el-form>

                        <el-table :data="table.columns" border>
                            <el-table-column prop="pk" label="主键" width="80">
                                <template slot-scope="scope">
                                    <i class="el-icon-success" v-if="scope.row.pk"></i>
                                </template>
                            </el-table-column>
                            <el-table-column prop="ai" label="自增" width="80" :formatter="aiFormatter">
                            </el-table-column>
                            <el-table-column prop="name" label="列名" width="160">
                            </el-table-column>
                            <el-table-column prop="data_type" label="数据类型" width="100" >
                            </el-table-column>
                            <el-table-column prop="length" label="数据长度" width="120" >
                            </el-table-column>
                            <el-table-column prop="null" label="可空" width="80" :formatter="nullFormatter">
                            </el-table-column>
                            <el-table-column prop="index" label="索引列" width="80" :formatter="indexFormatter">
                            </el-table-column>
                            <el-table-column label="枚举" width="120">
                                <template slot-scope="scope">
                                    <div>
                                        <el-popover placement="top-start" width="240"
                                        trigger="manual" v-if="scope.row.enum_list!=null&&scope.row.enum_list.length>0"
                                        v-model="scope.row.enum_visible">
                                        <p>{{scope.row.title}}</p>
                                        <el-table :data="scope.row.enum_list">
                                            <el-table-column width="80" property="key" label="字段"></el-table-column>
                                            <el-table-column width="80" property="value" label="值"></el-table-column>
                                            <el-table-column width="80" property="des" label="描述"></el-table-column>
                                        </el-table>
                                        <el-button type="text" icon="el-icon-message"
                                            slot="reference" @click="handleEnumVisibleChange(scope.row)">展示</el-button>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column prop="title" label="描述" >
                            </el-table-column>
                            <el-table-column label="操作" width="160" align="center">
                                <template slot-scope="scope">
                                    <el-button type="text" icon="el-icon-edit" @click="openColumnForm(scope.row)">编辑</el-button>
                                    <el-button type="text" icon="el-icon-delete" @click="deleteColumn(scope.row)" class="red" >删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>

                    <div class="button-center">
                        <el-button @click="db.tables.push({ id:0, did:db.id, name:'' })"
                        class="mb20 mt20" type="primary" size="medium" icon="el-icon-plus" circle title="新增数据表">
                        </el-button>
                    </div>
                </el-collapse-item>
            </el-collapse>
            <div class="pagination">
                <el-pagination background @current-change="search" layout="prev, pager, next" 
                :total="pageTotal" :page-size="pageSize" :current-page.sync="pageCurrent" :hide-on-single-page="true" >
                </el-pagination>
            </div>

        </div>

        <!-- 编辑DB弹出框 -->
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

        <!-- 编辑Column弹出框 -->
        <el-dialog title="编辑" :visible.sync="editColumnVisible" width="45%" @close="closeColumnForm('columnform')">
            <el-form ref="columnform" :model="columnForm" :rules="rules" label-width="100px">
                <el-form-item label="列名:" prop="name">
                    <el-input v-model="columnForm.name" maxlength="20" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="数据类型:" prop="data_type">
                    <el-select v-model="columnForm.data_type" placeholder="请选择">
                        <el-option
                        v-for="item in mysql_data_types"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="主键:">
                    <el-switch v-model="columnForm.pk" active-text="YES" inactive-text="NO">
                    </el-switch>
                </el-form-item>
                <el-form-item label="自增:">
                    <el-switch v-model="columnForm.ai" active-text="YES" inactive-text="NO">
                    </el-switch>
                </el-form-item>
                <el-form-item label="可空:">
                    <el-switch v-model="columnForm.null" active-text="YES" inactive-text="NO">
                    </el-switch>
                </el-form-item>
                <el-form-item label="索引列:">
                    <el-switch v-model="columnForm.index" active-text="YES" inactive-text="NO">
                    </el-switch>
                </el-form-item>
                <el-form-item label="长度:">
                    <el-input v-model="columnForm.length" maxlength="20" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="描述:" prop="title">
                    <el-input v-model="columnForm.title" maxlength="20" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="枚举:">
                    <el-table :data="columnForm.enum_list">
                        <el-table-column property="key" label="字段">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.key" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="value" label="值">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.value" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="des" label="描述">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.des" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column align="right" width="80">
                            <template slot="header">
                                <el-button @click="columnForm.enum_list.push({key:'', value:'' ,des:'' })" size="mini" type="primary">添加</el-button>
                            </template>
                            <template slot-scope="scope">
                                <el-button @click="columnForm.enum_list.splice(scope.$index, 1);" size="mini" type="text" icon="el-icon-delete" class="red">删除</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-form-item>
                
                <el-form-item label="备注信息:">
                    <el-input type="textarea" placeholder="请输入内容" v-model="columnForm.description"
                        maxlength="200" show-word-limit :autosize="{ minRows: 4, maxRows: 8}">
                    </el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="closeColumnForm('columnform')">取 消</el-button>
                <el-button @click="saveColumnForm('columnform')" type="primary">确 定</el-button>
            </span>
        </el-dialog>

        <!-- 连接弹出框 -->
        <el-dialog title="连接" :visible.sync="editConnVisible" width="50%" @close="closeConnForm('connForm')">
            <el-form ref="connForm" :model="connForm" :rules="rules" label-width="100px" :class="connInfo.length>0?'hide':''">
                <el-form-item label="数据库类型:">
                    <el-input v-model="project.data_base" :disabled="true"></el-input>
                </el-form-item>
                <el-form-item label="Host:" prop="host">
                    <el-input v-model="connForm.host"></el-input>
                </el-form-item>
                <el-form-item label="Port:" prop="port">
                    <el-input v-model="connForm.port"></el-input>
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
    import { debuglog } from 'util';
    export default {
        props: {
            pid: {
                type: String,
                default: '0'
            }
        },
        data() {
            return {
                loadConnUrl:'/api/dbimport/loaddb',
                saveConnUrl:'/api/dbimport/savedbs',
                loadProjectUrl: '/api//project/load',
                saveDBUrl: '/api/database/save',
                listDBUrl:'/api/database/list',
                loadDBUrl:'/api/database/load',
                deleteDBUrl:'/api/database/delete',
                saveColumnUrl: '/api/column/save',
                deleteColumnUrl:'/api/column/delete',
                saveTableUrl: '/api/table/save',
                deleteTableUrl:'/api/table/delete',
                loading: false,
                treeLoading:false,
                collapseLoading:false,
                pageTotal:0,
                pageSize:10,
                pageCurrent:1,
                search_word: '',
                editDBVisible: false,
                editColumnVisible: false,
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
                    host:'',
                    port:'',
                    user:'',
                    password:''
                },
                columnForm:{
                    enum_list:[]
                },
                mysql_data_types:[
                    {label:"int",value:"int"},
                    {label:"varchar",value:"varchar"},
                    {label:"bit",value:"bit"},
                    {label:"timestamp",value:"timestamp"},
                    {label:"longtext",value:"longtext"},
                    {label:"tinyint",value:"tinyint"},
                    {label:"datetime",value:"datetime"},
                    {label:"bigint",value:"bigint"},
                    {label:"decimal",value:"decimal"},
                    {label:"float",value:"float"},
                    {label:"date",value:"date"},
                    {label:"text",value:"text"},
                    {label:"double",value:"double"},
                    {label:"char",value:"char"},
                    {label:"time",value:"time"},
                ],
                rules: {
                    name: [
                        { required: true, message: '请输入名称', trigger: 'blur' }
                    ],
                    db_name: [
                        { required: true, message: '请输入数据库名称', trigger: 'blur' }
                    ],
                    host: [
                        { required: true, message: '请输入host', trigger: 'blur' }
                    ],
                    port: [
                        { required: true, message: '请输入port', trigger: 'blur' }
                    ],
                    user: [
                        { required: true, message: '请输入登录账户', trigger: 'blur' }
                    ],
                    password: [
                        { required: true, message: '请输入登录密码', trigger: 'blur' }
                    ],
                    data_type:[
                        { required: true, message: '请选择数据类型', trigger: 'blur' }
                    ],
                    title:[
                        { required: true, message: '请填写描述信息', trigger: 'blur' }
                    ]
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
            editDBForm(db){
                this.form.id=db.id;
                this.form.name=db.name;
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
                            this.dbData.forEach(function(_db,i) {
                                if (db.id==_db.id) {
                                    index = i;
                                }
                            })
                            this.dbData.splice(index, 1);
                        }
                    })
                })
            },
            closeConnForm(formName){
                this.editConnVisible = false
                this.connForm.host="";
                this.connForm.port="";
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
                                _db.tables = result.data.tables.map(table => {
                                    table.columns = table.columns.map(col => 
                                    {
                                        col.enum_list = []
                                        if (col.enum!=null&&col.enum!='') {
                                            let items = col.enum.split(";")
                                            items.forEach(function(item){
                                                if (item!=null&&item!=''){
                                                    let kvd = item.split(":")
                                                    if (kvd.length>2){
                                                        col.enum_list.push({key:kvd[0], value:kvd[1] ,des:kvd[2] });
                                                    }
                                                }
                                            })
                                        }
                                        return col
                                    })
                                    return table
                                })
                            }
                            this.collapseLoading = false
                        })
                    }
                }
            },
            handleTreeCheckChange(data, checked, indeterminate) {
                let key = ''
                if (data.tables==null) {
                    key = data.db_name+'_'+data.name
                } else {
                    return
                }
                if (checked) {
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
                this.connSelectInfo.forEach(function(value,key) {
                    topDBs.set(value.db_name,{ name: value.db_name, pid: _pid, tables: [] })
                })
                this.connSelectInfo.forEach(function(value,key){
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
                    }
                    this.treeLoading = false
                })
            },
            handleEnumVisibleChange(row) {
                if (row.enum_visible){
                    row.enum_visible = false
                } else {
                    row.enum_visible = true
                }
            },
            nullFormatter(row, column) {
                if (row.null==true){
                    return "YES"
                }
                return "NO"
            },
            aiFormatter(row, column) {
                if (row.ai==true){
                    return "YES"
                }
                return "NO"
            },
            indexFormatter(row, column) {
                if (row.index==true){
                    return "YES"
                }
                return "NO"
            },
            saveTable(table){
                if (table.name==null||table.name==''){
                    return
                }
                this.$axios.post(this.saveTableUrl, table).then(result=>{
                    if (result.success) {
                        table.id = result.data
                    }
                })
            },
            deleteTable(table){
                this.$confirm(`所属列也将全部删除,确定要删除表 : ${table.name}`, '提示信息', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.$axios.post(this.deleteTableUrl, table.id).then(result=>{
                        if (result.success) {
                            // 遍历删除页面上对应列的数据
                            this.dbData.forEach(function(db) {
                                if (db.id==table.did) {
                                    let index =0;
                                    db.tables.forEach(function(tab,i) {
                                        if (table.id==tab.id) {
                                            index = i;
                                        }
                                    })
                                    db.tables.splice(index, 1);
                                }
                            })
                        }
                    })
                })
            },
            openColumnForm(col) {
                this.columnForm = JSON.parse(JSON.stringify(col))
                this.old_data = col
                this.editColumnVisible = true
            },
            closeColumnForm(formName) {
                this.columnForm = {
                    enum_list:[]
                }
                this.editColumnVisible = false
                this.$refs[formName].resetFields();
            },
            saveColumnForm(formName) {
                this.$refs[formName].validate((valid) => {
                    if (!valid) {
                        return
                    }
                    let enumStr = ''
                    this.columnForm.enum_list.forEach(function(item) {
                        //{key:'', value:'' ,des:'' }
                        enumStr += `${item.key}:${item.value}:${item.des};`
                    })
                    this.columnForm.enum = enumStr
                    let columnForm = this.columnForm
                    this.$axios.post(this.saveColumnUrl, this.columnForm).then(result=>{
                        if (result.success) {
                            // 遍历修改页面上对应列的数据
                            this.dbData.forEach(function(db){
                                if (db.tables!=null){
                                    db.tables.forEach(function(table){
                                        if (table.id==columnForm.tid){
                                            table.columns.forEach(function(column){
                                                if (column.id==columnForm.id){
                                                    column.name = columnForm.name
                                                    column.data_type = columnForm.data_type
                                                    column.pk = columnForm.pk
                                                    column.index = columnForm.index
                                                    column.ai = columnForm.ai
                                                    column.length = columnForm.length
                                                    column.null = columnForm.null
                                                    column.enum = columnForm.enum
                                                    column.enum_list = columnForm.enum_list
                                                    column.title = columnForm.title
                                                    column.description = columnForm.description
                                                }
                                            })
                                        }
                                    })
                                }
                            })
                        }
                    })
                })
            },
            deleteColumn(column){
                this.$confirm(`确定要删除列 : ${column.name}`, '提示信息', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.$axios.post(this.deleteColumnUrl, column.id).then(result=>{
                        if (result.success) {
                            // 遍历删除页面上对应列的数据
                            this.dbData.forEach(function(db){
                                if (db.tables!=null){
                                    db.tables.forEach(function(table){
                                        if (table.id==column.tid){
                                            let index = -1
                                            table.columns.forEach(function(col,i){
                                                if (col.id==column.id){
                                                    index = i;
                                                }
                                            })
                                            table.columns.splice(index, 1);
                                        }
                                    })
                                }
                            })
                        }
                    })
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

    .mt20{
        margin-top: 20px;
    }

    .button-center{
        text-align:center
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
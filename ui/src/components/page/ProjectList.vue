<template>
    <div>
        <el-row :gutter="20">
            <el-col :span="6" v-for="pro in data" :key="pro.id">
                <div class="project-box">
                    <div class="body" :class="{mysql:pro.data_base=='mysql',mssql:pro.data_base=='mssql'}">
                        <el-button type="danger" icon="el-icon-delete" circle class="cbtn" @click="del(pro)"></el-button>
                        <el-button icon="el-icon-edit" circle class="cbtn" @click="edit(pro)"></el-button>
                    </div>
                    <div class="title">
                        <el-button type="text" @click="editDB(pro.id)" >{{pro.name}}</el-button>
                    </div>
                </div>
            </el-col>

            <el-col :span="6">
                <div class="project-box add" @click="add"></div>
            </el-col>
        </el-row>

        <!-- <div class="pagination">
            <el-pagination background @current-change="handlePageChange" layout="prev, pager, next" 
            :total="pageTotal" :page-size="pageSize" :current-page.sync="pageCurrent" :hide-on-single-page="false" >
            </el-pagination>
        </div> -->

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editVisible" width="70%" @close="close('form')">
            <el-form ref="form" :model="form" :rules="rules" label-width="100px">
                <el-form-item label="项目名称:" prop="name">
                    <el-input v-model="form.name" maxlength="20" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="编程语言:" prop="language">
                    <el-radio v-model="form.language" label="java" @change="handleLanguageChange">Java</el-radio>
                    <el-radio v-model="form.language" label="csharp" @change="handleLanguageChange">C#</el-radio>
                </el-form-item>
                <el-form-item label="数据库:" prop="data_base">
                    <el-radio v-model="form.data_base" label="mysql" :disabled="form.id>0">Mysql</el-radio>
                    <el-radio v-model="form.data_base" label="mssql" :disabled="form.id>0">Sqlserver</el-radio>
                </el-form-item>
                <el-form-item label="ORM:" prop="orm">
                    <el-select v-model="form.orm" placeholder="请选择">
                        <el-option v-for="item in orms" :key="item.value" :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="描述信息:">
                    <el-input type="textarea" placeholder="请输入内容" v-model="form.description"
                        maxlength="200" show-word-limit :autosize="{ minRows: 4, maxRows: 8}">
                    </el-input>
                </el-form-item>
                <el-form-item label="连接:">
                    <el-table :data="form.connection_list">
                        <el-table-column property="name" label="名称">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.name" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="value" label="Host" width="160">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.host" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="des" label="Port" width="80">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.port" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="des" label="User" width="120">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.user" ></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column property="des" label="Password" width="120">
                            <template slot-scope="scope">
                                <el-input v-model="scope.row.password" show-password></el-input>
                            </template>
                        </el-table-column>
                        <el-table-column align="right" width="80">
                            <template slot="header">
                                <el-button @click="form.connection_list.push({name:'', host:'' ,port:'', user:'', password:'' })" size="mini" type="primary">添加</el-button>
                            </template>
                            <template slot-scope="scope">
                                <el-button @click="form.connection_list.splice(scope.$index, 1);" size="mini" type="text" icon="el-icon-delete" class="red">删除</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="close('form')">取 消</el-button>
                <el-button @click="save('form')" type="primary">确 定</el-button>
            </span>
        </el-dialog>

    </div>
</template>

<script>
import { debuglog } from 'util';
    export default {
        data() {
            return {
                saveUrl: '/api/project/save',
                listUrl: '/api/project/list',
                deleteUrl: '/api/project/delete',
                search_word: '',
                editVisible: false,
                loading: false,
                data:[],
                pageTotal:0,
                pageSize:999999, //暂时不做分页，默认加载所有项目
                pageCurrent:1,
                form: {
                    id: 0,
                    name: '',
                    language: '',
                    description: '',
                    data_base: '',
                    orm:'',
                    connection_list:[]
                },
                rules: {
                    name: [
                        { required: true, message: '请输入项目名称', trigger: 'blur' }
                    ],
                    language: [
                        { required: true, message: '请选择编程语言', trigger: 'blur' }
                    ],
                    data_base: [
                        { required: true, message: '请选择项目所使用的数据库', trigger: 'blur' }
                    ],
                    orm: [
                        { required: true, message: '请选择ORM框架', trigger: 'blur' }
                    ]
                }
            }
        },
        computed: {
            orms() {
                if(this.form.language=='java') {
                    return [{label: "Mybatis",value: "mybatis"}]
                }
                if(this.form.language=='csharp') {
                    return [{label: "SmartSql",value: "smartSql"}]
                }
                return []
            }
        },
        created() {
            this.search();
        },
        methods: {
            search() {
                this.loading=true
                this.$axios.post(this.listUrl, { name:this.search_word, index:this.pageCurrent-1, size:this.pageSize, order_by:"project_created DESC" }).then(result=>{
                    if (result.success) {
                        this.data = result.data.list
                        this.pageTotal = result.data.total
                    }
                    this.loading=false
                })
            },
            add() {
                this.editVisible = true;
            },
            close(formName) {
                this.editVisible = false;
                this.$refs[formName].resetFields();
                this.form={
                    id: 0,
                    name: '',
                    language: '',
                    description: '',
                    data_base: '',
                    orm:'',
                    connection_list:[]
                }
                this.search()
            },
            edit(item) {
                this.form = {
                    id: item.id,
                    name: item.name,
                    language: item.language,
                    description: item.description,
                    data_base: item.data_base,
                    orm:item.orm,
                    connection_list:[]
                }
                this.editVisible = true;
            },
            editDB(pid){
                this.$router.push('/dbs/'+pid);
            },
            del(item) {
                this.$confirm(`确定要删除项目 : ${item.name}`, '提示信息', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.$axios.post(this.deleteUrl, item.id).then(result=>{
                        if (result.success) {
                            this.search()
                        }
                    })
                })
            },
            save(formName) {
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.$axios.post(this.saveUrl, this.form).then(result=>{
                            if (result.success) {
                                this.form.id=result.data
                                this.close(formName)
                            }
                        })
                    }
                });
            },
            handleLanguageChange(lab) {
                if(lab=='java') {
                    this.form.orm = "mybatis"
                }
                if(lab=='csharp') {
                    this.form.orm = "smartSql"
                }
            },
            handlePageChange(){
                this.search();
            }
        }
    }
</script>

<style scoped>
    .project-box {
        display:block;
        width:220px;
        height:240px;
        border:solid 1px #e6e6e6;
        margin-bottom: 20px;
        margin-top: 20px;
        margin-left:20px;
    }

    .project-box .mysql{
        background-image: url(../../assets/img/mysql-logo.jpg);
        background-repeat:no-repeat;
        background-size:100% 100%;
        -moz-background-size:100% 100%;
    }

    .project-box .mssql{
        background-image: url(../../assets/img/mssql-logo.jpg);
        background-repeat:no-repeat;
        background-size:100% 100%;
        -moz-background-size:100% 100%;
    }

    .project-box.add{
        background-image: url(../../assets/img/box-add.jpg);
        background-repeat:no-repeat;
        background-size:100% 100%;
        -moz-background-size:100% 100%;
        cursor: pointer;
    }

    .project-box .body .cbtn {
        float: right;
        margin-right: 10px;
        margin-top: 10px;
    }

    .project-box .body {
        height:200px;
        border-bottom:solid 1px #e6e6e6;
    }

    .project-box .title {
        height:40px;
        text-align: center;
        padding: 5px;
    }
</style>
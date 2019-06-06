<template>
    <div class="table">
        
        <div class="container">
            <div class="handle-box">
                <el-input v-model="search_word" placeholder="项目名称" class="handle-input mr10"></el-input>
                <el-button icon="el-icon-search" circle @click="search" title="查询"></el-button>
                <el-button type="primary" icon="el-icon-edit" circle @click="add" title="新增"></el-button>
            </div>
            <el-table :data="data" class="table" ref="multipleTable" v-loading="loading">
                <el-table-column type="expand">
                    <template slot-scope="props">
                        <el-form label-position="left" inline class="table-expand">
                        <el-form-item label="项目名称">
                            <span>{{ props.row.name }}</span>
                        </el-form-item>
                        <el-form-item label="编程语言">
                            <span>{{ props.row.language }}</span>
                        </el-form-item>
                        <el-form-item label="数据库">
                            <span>{{ props.row.data_base }}</span>
                        </el-form-item>
                        <el-form-item label="ORM">
                            <span>{{ props.row.orm }}</span>
                        </el-form-item>
                        <el-form-item label="描述信息">
                            <span>{{ props.row.description }}</span>
                        </el-form-item>
                        </el-form>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="项目名称" sortable >
                </el-table-column>
                <el-table-column label="操作" width="200" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-edit" @click="edit(scope.$index, scope.row)">编辑</el-button>
                        <el-button type="text" icon="el-icon-delete" class="red" >删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination">
                <el-pagination layout="prev, pager, next" :total="1000">
                </el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editVisible" width="40%" @close="close('form')">
            <el-form ref="form" :model="form" :rules="rules" label-width="100px">
                <el-form-item label="项目名称:" prop="name">
                    <el-input v-model="form.name" maxlength="20" show-word-limit></el-input>
                </el-form-item>
                <el-form-item label="编程语言:" prop="language">
                    <el-radio v-model="form.language" label="java" >Java</el-radio>
                    <el-radio v-model="form.language" label="csharp" >C#</el-radio>
                </el-form-item>
                <el-form-item label="数据库:" prop="data_base">
                    <el-radio v-model="form.data_base" label="mysql" >Mysql</el-radio>
                    <el-radio v-model="form.data_base" label="sqlserver" >Sqlserver</el-radio>
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
                search_word: '',
                editVisible: false,
                loading: false,
                data:[],
                form: {
                    id: 0,
                    name: '',
                    language: '',
                    description: '',
                    data_base: '',
                    orm:''
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
                this.$axios.post(this.listUrl, {name:this.search_word}).then(result=>{
                    if (result.success) {
                        this.data=result.data
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
                    orm:''
                }
            },
            edit(index, row) {
                this.idx = index;
                const item = this.data[index];
                this.form = {
                    id: item.id,
                    name: item.name,
                    language: item.language,
                    description: item.description,
                    data_base: item.data_base,
                    orm:item.orm
                }
                this.editVisible = true;
            },
            save(formName) {
                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        this.$axios.post(this.saveUrl, this.form).then(result=>{
                            if (result.success) {
                                this.form.id=result.data
                            }
                        })
                    }
                });
            },
            formatter(row, column) {
                switch (row.language){
                    case 'java':
                        return "Java"
                    case 'csharp':
                        return "C#"
                    default:
                        return ''
                }
            }
        }
    }
</script>

<style scoped>
    .table{
        width: 100%;
        font-size: 14px;
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
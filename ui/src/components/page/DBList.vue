<template>
    <div>
        <div class="container">
            <div class="handle-box">
                <el-input v-model="search_word" placeholder="数据库名称" class="handle-input mr10"></el-input>
                <el-button icon="el-icon-search" circle @click="search" title="查询"></el-button>
                <el-button type="primary" icon="el-icon-edit" circle @click="addDB" title="新增"></el-button>
            </div>
            <el-table :data="dbData" class="table" ref="DBTable" v-loading="loading">
                <el-table-column prop="name" label="数据库名称" sortable >

                </el-table-column>
            </el-table>
            <div class="pagination">
                <el-pagination background @current-change="search" layout="prev, pager, next" 
                :total="pageTotal" :page-size="pageSize" :current-page.sync="pageCurrent" :hide-on-single-page="true" >
                </el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editDBVisible" width="40%">
            <el-form ref="form" :model="form" :rules="rules" label-width="100px">
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
                <el-button @click="editDBVisible = false">取 消</el-button>
                <el-button @click="saveDB('form')" type="primary">确 定</el-button>
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
                saveDBUrl: '/api/database/save',
                listDBUrl:'/api/database/list',
                loading: false,
                pageTotal:0,
                pageSize:10,
                pageCurrent:1,
                search_word: '',
                editDBVisible: false,
                name: '',
                dbData:[],
                form: {
                    id: 0,
                    pid: parseInt(this.pid),
                    name: '',
                    description: ''
                },
                rules: {
                    name: [
                        { required: true, message: '请输入项目名称', trigger: 'blur' }
                    ]
                }
            }
        },
        created() {
            this.search();
        },
        computed:{
            
        },
        methods: {
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
            addDB() {
                this.editDBVisible = true;
                this.form.id=0;
                this.form.name="";
                this.form.description="";
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
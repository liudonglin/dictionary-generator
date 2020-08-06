<template>
    <div>
        <el-row>
            <el-form :inline="true" :model="form" class="mb20">
                <el-form-item label="模版名称:">
                    <el-input v-model="form.name" clearable maxlength="100"></el-input>
                </el-form-item>
                <el-form-item label="编程语言:">
                        <el-select v-model="form.language" clearable placeholder="请选择">
                            <el-option label="Java" value="java"></el-option>
                            <el-option label="C#" value="csharp"></el-option>
                            <el-option label="Golang" value="go"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="数据库:">
                        <el-select v-model="form.data_base" clearable placeholder="请选择">
                            <el-option label="Mysql" value="mysql"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="ORM:">
                        <el-select v-model="form.orm" clearable placeholder="请选择">
                            <el-option key="mybatis" label="Mybatis" value="mybatis"></el-option>
                            <el-option key="smartSql" label="SmartSql" value="smartSql"></el-option>
                            <el-option key="gorm" label="Gorm" value="gorm"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="search">查询</el-button>
                    </el-form-item>
            </el-form>
        </el-row>

        <el-row>
            <el-table :data="data" border style="width: 100%" class="mb20">
                <el-table-column prop="name" label="模版名称"></el-table-column>
                <el-table-column prop="language" label="编程语言" width="120"></el-table-column>
                <el-table-column prop="data_base" label="数据库" width="120"></el-table-column>
                <el-table-column prop="orm" label="映射框架" width="120"></el-table-column>
                <el-table-column prop="type" label="模版类型" width="120" 
                    :filters="[{ text: '初始', value: 'init' }, { text: '自定义', value: 'custom' }]"
                    :filter-method="filterType"
                    filter-placement="bottom-end">
                    <template slot-scope="scope">
                        <el-tag :type="scope.row.type === 'init' ? 'info' : 'success'"
                            disable-transitions>{{scope.row.type=== 'init' ? '初始' : '自定义'}}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="updated" label="编辑时间" width="180"></el-table-column>
                <el-table-column label="操作" width="180" align="center">
                        <template slot-scope="scope">
                            <el-button type="text" icon="el-icon-edit" @click="openColumnForm(scope.row)">查看</el-button>
                            <el-button v-if="scope.row.type === 'custom'" type="text" icon="el-icon-delete" @click="deleteColumn(scope.row)" class="red" >删除</el-button>
                        </template>
                    </el-table-column>
            </el-table>

            <el-pagination background layout="prev, pager, next" :total="pageTotal" 
                :page-size="pageSize" :current-page.sync="pageCurrent" @current-change="search">
            </el-pagination>
        </el-row>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                listUrl: '/api/templete/list',
                loading: false,
                data:[],
                search_word: '',
                pageTotal:0,
                pageSize:20, //暂时不做分页，默认加载所有项目
                pageCurrent:1,
                form: {
                    name:'',
                    language: '',
                    orm:'',
                    data_base: '',
                }
            }
        },
        created() {
            this.search();
        },
        methods: {
            search() {
                this.loading=true
                this.$axios.post(this.listUrl, { 
                    name:this.form.name, 
                    language:this.form.language,
                    data_base:this.form.data_base,
                    orm:this.form.orm,
                    index:this.pageCurrent-1, 
                    size:this.pageSize, 
                    order_by:"templete_updated DESC" })
                    .then(result=>{
                    if (result.success) {
                        this.data = result.data.list
                        this.pageTotal = result.data.total
                    }
                    this.loading=false
                })
            },
            filterType(value, row) {
                return row.type === value;
            },
            openColumnForm(col) {
                this.$router.push('/tplmgt/'+col.id);
            },
            deleteColumn(column){
                this.$confirm(`确定要删除该模版 : ${column.name}`, '提示信息', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    
                })
            }
        }
    }
</script>

<style scoped>

    .mb20{
        margin-bottom: 20px;
    }

    .red{
        color: #F56C6C;
    }

</style>
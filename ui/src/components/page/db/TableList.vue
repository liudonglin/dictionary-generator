<template>
<div v-loading="dbLoading" style="height: 100%;">

    <div class="mb60" style="border-bottom:solid 1px #e6e6e6;" v-if="dbid>0">
        <el-form :inline="true" :model="tableForm" :rules="rules">
            <el-form-item label="表 名:" prop="name">
                <el-input v-model="tableForm.name" style="width:280px;" maxlength="40" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="描 述:">
                <el-input v-model="tableForm.description" style="width:280px;" maxlength="200" show-word-limit></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" icon="el-icon-check" @click="saveTable(tableForm)" circle title="保存表"></el-button>
            </el-form-item>
        </el-form>
    </div>

    <div v-for="table in tables" :key="table.id" class="mb60">
        <el-form :inline="true" :model="table" :rules="rules" ref="tableform">
            <el-form-item label="表 名:" prop="name">
                <el-input v-model="table.name" style="width:280px;" maxlength="40" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="描 述:">
                <el-input v-model="table.description" style="width:280px;" maxlength="200" show-word-limit></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" icon="el-icon-check" @click="saveTable(table)" circle title="保存表"></el-button>
                <el-button type="danger" icon="el-icon-delete" @click="deleteTable(table)" circle title="删除表"></el-button>
                <el-button type="primary" icon="el-icon-plus" @click="openColumnForm({ id:0, tid:table.id })" circle title="添加列"></el-button>
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
            <el-table-column prop="length" label="数据长度" width="100" >
            </el-table-column>
            <el-table-column prop="null" label="可空" width="80" :formatter="nullFormatter">
            </el-table-column>
            <el-table-column prop="index" label="索引列" width="80" :formatter="indexFormatter">
            </el-table-column>
            <el-table-column label="枚举" width="120">
                <template slot-scope="scope">
                    <div>
                        <el-popover placement="top-start" width="360"
                        trigger="manual" v-if="scope.row.enum_list!=null&&scope.row.enum_list.length>0"
                        v-model="scope.row.enum_visible">
                        <p>{{scope.row.title}}</p>
                        <el-table :data="scope.row.enum_list">
                            <el-table-column property="key" label="字段"></el-table-column>
                            <el-table-column width="80" property="value" label="值"></el-table-column>
                            <el-table-column property="des" label="描述"></el-table-column>
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

    <!-- 编辑Column弹出框 -->
    <el-dialog title="编辑" :visible.sync="editColumnVisible" width="45%" @close="closeColumnForm('columnform')">
        <el-form ref="columnform" :model="columnForm" :rules="rules" label-width="100px">
            <el-form-item label="列名:" prop="name">
                <el-input v-model="columnForm.name" maxlength="40" show-word-limit></el-input>
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
                <el-switch v-model="columnForm.pk" @change="handlePKChange" active-text="YES" inactive-text="NO">
                </el-switch>
            </el-form-item>
            <el-form-item label="自增:">
                <el-switch v-model="columnForm.ai" active-text="YES" inactive-text="NO" :disabled="!columnForm.pk">
                </el-switch>
            </el-form-item>
            <el-form-item label="可空:">
                <el-switch v-model="columnForm.null" active-text="YES" inactive-text="NO" :disabled="columnForm.pk">
                </el-switch>
            </el-form-item>
            <el-form-item label="索引列:">
                <el-switch v-model="columnForm.index" active-text="YES" inactive-text="NO">
                </el-switch>
            </el-form-item>
            <el-form-item label="长度:">
                <el-input v-model="columnForm.length" maxlength="20" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="标题:" prop="title">
                <el-input v-model="columnForm.title" maxlength="40" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="枚举:" v-if="columnForm.data_type=='int'||columnForm.data_type=='bit'||columnForm.data_type=='tinyint'">
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
            
            <el-form-item label="描述:">
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

</div>
</template>
    
<script>
    import bus from '../../common/bus';
    export default {
        data() {
            return {
                listTableUrl:'/api/table/list',
                saveTableUrl: '/api/table/save',
                deleteTableUrl:'/api/table/delete',
                saveColumnUrl: '/api/column/save',
                deleteColumnUrl:'/api/column/delete',
                dbid:0,
                dbLoading: false,
                editColumnVisible:false,
                tables:[],
                search_word:'',
                tableForm:{
                    name:'',
                    description:'',
                    id:0,
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
            bus.$on('dbMgtSelectDBChange', val => {
                this.dbid = val;
                this.search();
            })

            bus.$on('headerSearchWordChange', val => {
                this.search_word = val;
                if (this.dbid==0) {
                    return;
                }
                this.search();
            })
        },
        methods:{
            search() {
                this.dbLoading = true;
                let postData = { name:this.search_word, did:this.dbid, index:0, size:999999, order_by:"table_created DESC" };
                this.$axios.post(this.listTableUrl,postData).then(result=>{
                    if (result.success) {
                        this.tables = result.data.list.map(table => {
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
                    this.dbLoading = false
                })
            },
            saveTable(table) {
                if (table.name==null||table.name==''){
                    return
                }
                table.did = this.dbid;
                this.$axios.post(this.saveTableUrl, table).then(result=>{
                    if (result.success) {
                        this.tableForm = {
                            name:'',
                            description:'',
                            id:0,
                        }
                        this.search()
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
                            let index =0;
                            this.tables.forEach(function(tab,i) {
                                if (table.id==tab.id) {
                                    index = i;
                                }
                            })
                            this.tables.splice(index, 1);
                        }
                    })
                })
            },
            openColumnForm(col) {
                if (col.id==0){
                    col.name = ''
                    col.data_type = ''
                    col.pk = false
                    col.index = false
                    col.ai = false
                    col.length = ''
                    col.null = false
                    col.enum = ''
                    col.enum_list = []
                    col.title = ''
                    col.description = ''
                }

                this.columnForm = JSON.parse(JSON.stringify(col))
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

                    // 非枚举字段去除枚举
                    let dtype =this.columnForm.data_type
                    if(dtype!="int"&&dtype!="bit"&&dtype!="tinyint"){
                        this.columnForm.enum_list=[];
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
                            // 遍历修改页面上对应列的数据 start
                            this.tables.forEach(function(table){
                                if (table.id==columnForm.tid) {
                                    if (columnForm.id>0) {
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
                                    } else {
                                        columnForm.id=result.data
                                        table.columns.push(JSON.parse(JSON.stringify(columnForm)))
                                    }
                                }
                            })
                            // 遍历修改页面上对应列的数据 end
                            this.closeColumnForm(formName);
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
                            this.tables.forEach(function(table){
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
                })
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
            handleEnumVisibleChange(row) {
                if (row.enum_visible){
                    row.enum_visible = false
                } else {
                    row.enum_visible = true
                }
            },
            handlePKChange(val){
                this.columnForm.ai=val
                if (val){
                    this.columnForm.null=false
                }
            },
        }
    }
</script>

<style scoped>
    
    .mb60{
        margin-bottom: 60px;
    }

    .red{
        color: #F56C6C;
    }

</style>
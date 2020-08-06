<template>
    <div>
        <el-row :gutter="20">
            <el-col :span="12" :offset="6">
                <el-form ref="form" :model="form" label-width="80px">
                    <el-form-item label="模版名称">
                        <el-input v-model="form.name"></el-input>
                    </el-form-item>
                    <el-form-item label="编程语言:" prop="language">
                        <el-radio v-model="form.language" label="java" >Java</el-radio>
                        <el-radio v-model="form.language" label="csharp" >C#</el-radio>
                        <el-radio v-model="form.language" label="go" >Golang</el-radio>
                    </el-form-item>
                    <el-form-item label="映射框架:" prop="orm">
                        <el-select v-model="form.orm" placeholder="请选择">
                            <el-option v-for="item in orms" :key="item.value" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="数据库:" prop="data_base">
                        <el-checkbox v-model="form.data_base" label="mysql">Mysql</el-checkbox>
                    </el-form-item>
                </el-form>
            </el-col>
            <el-col :span="6"></el-col>
        </el-row>

        <el-row>
             <el-col :span="16" :offset="4">
                <mavon-editor v-model="form.content" :toolbarsFlag="false" :shortCut="false" :subfield="false" defaultOpen="edit" placeholder="请编辑模版内容"/>
             </el-col>
             <el-col :span="4"></el-col>
        </el-row>
    </div>
</template>

<script>
import { mavonEditor } from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'
    export default {
        data() {
            return {
                loadUrl: '/api/templete/loadById',
                loading: false,
                templete_id:0,
                form: {
                    name:'',
                    content:'',
                    language: '',
                    orm:'',
                    data_base: '',
                    type:'custom'
                }
            }
        },
        components: {
            mavonEditor
        },
        mounted(){
            this.templete_id = parseInt(this.$route.params.templeteId)
            if(this.templete_id!=0) {
                this.load();
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
                if(this.form.language=='go') {
                    return [{label: "Gorm",value: "gorm"}]
                }
                return []
            }
        },
        created() {
        },
        methods: {
            load() {
                this.loading=true
                this.$axios.post(this.loadUrl, { 
                        templete_id:this.templete_id
                    })
                    .then(result=>{
                    if (result.success) {
                        this.form.name = result.data.name;
                        this.form.language = result.data.language;
                        this.form.orm = result.data.orm;
                        this.form.data_base = result.data.data_base;
                        this.form.content = result.data.content;
                    }
                    this.loading=false
                })
            }
        }
    }
</script>
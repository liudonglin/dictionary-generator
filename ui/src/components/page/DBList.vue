<template>
    <div>
        <h1> {{pid}} </h1>
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
                name: ''
            }
        },
        computed:{
            
        }
    }
</script>
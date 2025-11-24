<script lang="ts" setup>
import { GetConfig, SetConfig } from '../../bindings/tieEpub/tieepubservice'
import { onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
const cookieForm = reactive({
    cookie: '',
})
onMounted(() => {
    GetConfig("cookie").then(value => {
        cookieForm.cookie = value
    })
})
const updateCookie = () => {
    SetConfig("cookie", cookieForm.cookie).then(result => {
        if (result) {
            ElMessage.success({ message: 'Cookie设置成功', grouping: true, plain: true, placement: 'bottom' })
        } else {
            ElMessage.error({ message: 'Cookie设置失败', grouping: true, plain: true, placement: 'bottom' })
        }
    })
}

</script>
<template>
    <h1>Cookie设置</h1>
    <el-form label-position="top" :model="cookieForm">
        <el-form-item label="百度Cookie" prop="desc">
            <el-input placeholder="请填写已登录百度账号的有效Cookie" resize="none" :rows="10" v-model="cookieForm.cookie"
                type="textarea" />
        </el-form-item>
        <el-form-item>
            <div class="right">

                <el-button type="primary" @click="updateCookie">
                    保存Cookie
                </el-button>
            </div>
        </el-form-item>
    </el-form>
</template>
<style scoped lang="scss">
h1 {
    text-align: left;
    color: rgb(64, 130, 245);
    border-bottom: 2px solid rgb(64, 130, 245);
    font-size: 30px;
}

.right {
    display: flex;
    width: 100%;
    justify-content: right;
}
</style>
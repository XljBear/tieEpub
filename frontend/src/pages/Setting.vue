<script lang="ts" setup>
import { GetConfig, SetConfig } from '../../bindings/tieEpub/tieepubservice'
import { onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Browser } from "@wailsio/runtime"
const settingForm = reactive({
    aiApiKey: '',
})
onMounted(() => {
    GetConfig("ai-api-key").then(key => {
        settingForm.aiApiKey = key
    })
})
const updateSetting = () => {
    SetConfig("ai-api-key", settingForm.aiApiKey).then(result => {
        if (result) {
            ElMessage.success({ message: '设置保存成功', grouping: true, plain: true, placement: 'bottom' })
        } else {
            ElMessage.error({ message: '设置保存失败', grouping: true, plain: true, placement: 'bottom' })
        }
    })
}

const openLink = () => {
    Browser.OpenURL('https://cloud.siliconflow.cn/me/account/ak')
}
</script>
<template>
    <h1>设置</h1>
    <el-form label-position="top" :model="settingForm">
        <el-form-item label="Siliconflow API Key" prop="desc">
            <el-input type="password" show-password placeholder="请填写您的硅基流动API Key" v-model="settingForm.aiApiKey" />
            <div @click="openLink" class="link">从哪获取我的API Key?</div>
        </el-form-item>
        <el-form-item>
            <div class="right">
                <el-button type="primary" @click="updateSetting">
                    保存设置
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
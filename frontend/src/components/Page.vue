<script lang="ts" setup>
import MenuItem from './MenuItem.vue'
import Cookie from '../pages/Cookie.vue'
import Setting from '../pages/Setting.vue'
import Download from '../pages/Download.vue'
import { ref } from 'vue'
const nowPage = ref("download")
const toPage = (page: string) => {
    nowPage.value = page
}
</script>
<template>
    <div class="page">
        <div class="menus">
            <div class="logo">
                <img alt="" src="/images/logo.png" />
                <p>百度贴吧小说下载器</p>
                <p>v0.0.1</p>
            </div>
            <div class="menu-box">
                <MenuItem @click="toPage('download')" :active="nowPage == 'download'" icon="&#xe614;" label="下载" />
                <MenuItem @click="toPage('cookie')" :active="nowPage == 'cookie'" icon="&#xe654;" label="Cookie 设置" />
                <MenuItem @click="toPage('setting')" :active="nowPage == 'setting'" icon="&#xe60f;" label="设置" />
            </div>
        </div>
        <div class="page-content">
            <el-scrollbar height="100%">
                <div class="page-box">
                    <div v-show="nowPage == 'download'">
                        <Download />
                    </div>
                    <Cookie v-if="nowPage == 'cookie'" />
                    <Setting v-if="nowPage == 'setting'" />
                </div>
            </el-scrollbar>
        </div>
    </div>
</template>
<style scoped lang="scss">
.page {
    height: calc(100vh - 30px);
    // height: 100%;
    display: flex;
    border-bottom: 1px solid #999;
    border-left: 1px solid #999;
    border-right: 1px solid #999;
    border-radius: 0 0 8px 8px;
    box-sizing: border-box;

    .menus {
        height: 100%;
        width: 220px;
        box-shadow: 2px 0 4px rgba(0, 0, 0, 0.2);

        .logo {
            margin-top: 20px;
            width: 100%;
            padding: 20px 0;
            box-sizing: border-box;
            border-bottom: 1px solid #eee;

            img {
                width: 70px;
                height: 70px;
            }

            p {
                font-size: 12px;
                margin: 0;
                color: #666;
            }
        }

        .menu-box {
            padding: 0 20px;
        }
    }

    .page-content {
        flex: 1;
        padding: 0 5px;
        overflow-y: hidden;
        box-sizing: border-box;

        .page-box {
            padding: 5px 20px;
            box-sizing: border-box;
        }
    }
}
</style>
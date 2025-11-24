<script lang="ts" setup>
import { StartDownload, GetTieData, SaveEPUB, OnReady, CreateAiCover, DeleteChapter } from '../../bindings/tieEpub/tieepubservice'
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Events } from "@wailsio/runtime"
import { TieContent } from '../../bindings/tieEpub/tieba'
const downloadLoading = ref(false)
const downloadProgress = ref(0)
const downloadForm = reactive({
    url: '',
    minimumWord: 100,
    onlyLZ: true,
    filterLink: true,
    filterImg: true,
})
const aiPrompt = ref("")
const canCreateEpub = ref(false)
const tieData = ref<TieContent>()
onMounted(() => {
    Events.On("downloadError", (msg: { data: string }) => {
        downloadLoading.value = false
        ElMessage.error({ message: msg.data, grouping: true, plain: true, placement: 'bottom' })
    })
    Events.On("downloadProcess", (process: { data: number }) => {
        downloadProgress.value = process.data
    })
    Events.On("downloadSuccess", () => {
        downloadLoading.value = false
        refreshTieData()
        ElMessage.success({ message: "小说章节数据下载完成", grouping: true, plain: true, placement: 'bottom' })
    })
    Events.On("coverError", (msg: { data: string }) => {
        ElMessage.error({ message: "AI封面生成失败:" + msg.data, grouping: true, plain: true, placement: 'bottom' })
    })
    setTimeout(() => OnReady(), 500)
})
const refreshTieData = () => {
    GetTieData().then((result: [success: boolean, content: TieContent]) => {
        if (!result[0]) {
            return
        }
        canCreateEpub.value = true
        tieData.value = result[1]
        nowSelectChapterIndex.value = 0
        aiCoverSelectChapter.value = 0
    })
}
const enableAiCover = ref(false)
const coverLoading = ref(false)
const coverImgData = ref("")
const aiCoverSelectChapter = ref(0)
const getCover = () => {
    coverLoading.value = true
    CreateAiCover(aiPrompt.value, aiCoverSelectChapter.value).then((result: [imgBase64: string, error: string]) => {
        if (result[1] != "") {
            ElMessage.error({ message: "AI封面生成失败:" + result[1], grouping: true, plain: true, placement: 'bottom' })
        } else {
            coverImgData.value = result[0]
        }
    }).finally(() => {
        coverLoading.value = false
    })
}
onUnmounted(() => {
    Events.OffAll()
})
const startDownload = () => {
    downloadLoading.value = true
    downloadProgress.value = 0
    StartDownload(downloadForm.url, downloadForm.minimumWord, downloadForm.onlyLZ, downloadForm.filterLink, downloadForm.filterImg).catch(() => {
        downloadLoading.value = false
        ElMessage.error({ message: "启动任务发生错误，请检查参数", grouping: true, plain: true, placement: 'bottom' })
    })
}
const saveEPUB = () => {
    SaveEPUB(enableAiCover.value).then((result: boolean) => {
        if (result) {
            ElMessage.success({ message: "小说文件保存成功", grouping: true, plain: true, placement: 'bottom' })
        }
    })
}
const back = () => {
    aiCoverSelectChapter.value = 0
    enableAiCover.value = false
    coverImgData.value = ""
    aiPrompt.value = ""
    nowSelectChapterIndex.value = 0
    canCreateEpub.value = false
    tieData.value = undefined
}
const nowSelectChapterIndex = ref(0)
const selectChapter = (index: number) => {
    nowSelectChapterIndex.value = index
}
const deleteChapter = (index: number) => {
    DeleteChapter(index).then(() => {
        refreshTieData()
        ElMessage.success({ message: "章节删除成功", grouping: true, plain: true, placement: 'bottom' })
    })
}
</script>
<template>
    <h1>{{ canCreateEpub ? "创建EPUB文件" : "准备下载" }}</h1>
    <div v-if="!canCreateEpub">
        <el-form v-loading="downloadLoading" element-loading-text="数据获取中..." label-position="top" :model="downloadForm">
            <el-form-item label="贴子链接">
                <el-input size="large" placeholder="请输入小说网文类贴链接" v-model="downloadForm.url" />
            </el-form-item>
            <el-collapse expand-icon-position="left">
                <el-collapse-item title="高级设置">
                    <div>
                        <el-form-item label="楼层内容最少字数">
                            <el-input-number min="1" size="large"
                                @blur="() => { if (!downloadForm.minimumWord) { downloadForm.minimumWord = 100 } }"
                                style="width: 300px;" placeholder="请输入最少字数限制" v-model="downloadForm.minimumWord" />
                        </el-form-item>
                        <el-form-item label="">
                            <el-checkbox v-model="downloadForm.onlyLZ" label="只看楼主" />
                        </el-form-item>
                        <el-form-item label="">
                            <el-checkbox v-model="downloadForm.filterLink" label="过滤包含链接的楼层" />
                        </el-form-item>
                    </div>
                </el-collapse-item>
            </el-collapse>
            <el-form-item>
                <div class="right">
                    <el-button type="primary" @click="startDownload">
                        开始获取数据
                    </el-button>
                </div>
            </el-form-item>
        </el-form>
        <el-progress v-if="downloadLoading" striped striped-flow :duration="20" :text-inside="true" :stroke-width="20"
            :percentage="downloadProgress" />
    </div>
    <div class="epub" v-else>
        <div class="book-info">
            <div class="title">{{ tieData?.Title || "空白占位标题" }}</div>
            <div class="author"> 作者：{{ tieData?.Author || "佚名" }}<br>共{{ tieData?.TotalContent.length || 0 }}章节</div>
        </div>
        <div class="book-content">
            <div class="cover-block">
                <div v-loading="coverLoading" class="cover">
                    <div class="ai-cover" v-if="enableAiCover">
                        <img v-if="coverImgData" :src="coverImgData" />
                        <div class="ai-cover-empty" v-else>等待生成AI封面</div>
                    </div>
                    <div class="title-cover" v-else>{{ tieData?.Title || "空白占位标题" }}</div>
                </div>
                <div class="ai-cover-prompt">
                    <div class="option">
                        <el-checkbox v-model="enableAiCover" label="创作AI封面" />
                        <el-select :disabled="!enableAiCover" v-model="aiCoverSelectChapter" placeholder="选择参考章节">
                            <el-option v-for="(_, index) in tieData?.TotalContent" :label="`第${index + 1}章`"
                                :value="index" />
                        </el-select>
                    </div>
                    <div v-if="enableAiCover" class="prompt">
                        <el-input placeholder="风格设定提示词" v-model="aiPrompt" />
                        <el-button :loading="coverLoading" type="primary" @click="getCover">生成</el-button>
                    </div>
                </div>
            </div>
            <div class="chapter-block">
                <div class="chapter-list">
                    <el-scrollbar height="100%">
                        <div @click="selectChapter(index)" v-for="(_, index) in tieData?.TotalContent" :key="index"
                            :class="['chapter-item', nowSelectChapterIndex === index ? 'active' : '']">
                            <el-popconfirm width="160" @confirm="deleteChapter(index)" placement="bottom"
                                :hide-icon="true" cancel-button-text="取消" confirm-button-text="删除"
                                :title="`是否删除第${index + 1}章节？`">
                                <template #reference>
                                    <div class="remove">&#xeafb;</div>
                                </template>
                            </el-popconfirm>
                            第{{ index + 1 }}章
                        </div>
                    </el-scrollbar>
                </div>
                <div class="chapter-content">
                    <el-scrollbar height="100%">
                        <div class="content" v-html="tieData?.TotalContent[nowSelectChapterIndex]">
                        </div>
                    </el-scrollbar>
                </div>
            </div>
        </div>
        <div class="right">
            <el-button type="primary" @click="back">
                返回
            </el-button>
            <el-button type="primary" @click="saveEPUB">
                保存
            </el-button>
        </div>
    </div>
</template>
<style scoped lang="scss">
h1 {
    text-align: left;
    color: rgb(64, 130, 245);
    border-bottom: 2px solid rgb(64, 130, 245);
    font-size: 30px;
}

.right {
    margin-top: 20px;
    display: flex;
    width: 100%;
    justify-content: right;
}

.epub {
    .book-info {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .title {
            text-align: left;
            line-height: 32px;
            font-size: 20px;
            font-weight: bold;
            color: #333;
        }

        .author {
            text-align: right;
            color: #666;
            font-size: 14px;
            line-height: 20px;
        }
    }

    .book-content {
        display: flex;
        padding: 10px;
        box-sizing: border-box;
        margin-top: 10px;

        .cover-block {
            width: 250px;
            display: flex;
            flex-direction: column;
            align-items: center;

            .cover {
                width: 192px;
                height: 256px;
                border-radius: 8px;
                overflow: hidden;

                .ai-cover {
                    height: 100%;
                    width: 100%;
                    border-radius: 8px;
                    overflow: hidden;

                    .ai-cover-empty {
                        display: flex;
                        height: 100%;
                        width: 100%;
                        font-size: 20px;
                        align-items: center;
                        color: #333;
                        justify-content: center;
                        border: 1px solid #ccc;
                        border-radius: 8px;
                        box-sizing: border-box;
                    }

                    img {
                        border-radius: 8px;
                        overflow: hidden;
                        width: 100%;
                        height: 100%;
                    }
                }


                .title-cover {
                    padding: 10px 10px 10px 30px;
                    box-sizing: border-box;
                    display: flex;
                    height: 100%;
                    width: 100%;
                    font-size: 20px;
                    align-items: center;
                    color: #333;
                    justify-content: center;
                    background-image: url("/images/book.png");
                    background-position: center center;
                    background-repeat: no-repeat;
                }
            }

            .ai-cover-prompt {
                margin-top: 25px;
                width: 100%;

                .option {
                    padding: 5px;
                    box-sizing: border-box;
                    display: flex;
                    gap: 10px;
                    width: 100%;
                }

                .prompt {
                    margin-top: 5px;
                    display: flex;
                    gap: 2px;
                }
            }
        }

        .chapter-block {
            height: 280px;
            width: auto;
            flex: 1;
            border-radius: 10px;
            overflow: hidden;
            border: 1px solid #ccc;
            display: flex;

            .chapter-list {
                height: 280px;
                width: 120px;
                overflow: hidden;
                background-color: #f5f5f5;
                display: flex;
                flex-direction: column;
                gap: 5px;

                .chapter-item {
                    padding: 5px;
                    cursor: pointer;
                    height: 30px;
                    font-size: 12px;
                    border-bottom: 1px solid #ccc;
                    box-sizing: border-box;
                    color: #333;
                    line-height: 20px;
                    position: relative;

                    .remove {
                        position: absolute;
                        right: 5px;
                        color: #aaa;
                        font-size: 18px;
                        cursor: pointer;
                        opacity: 0;
                        transition: all 0.3s ease;
                    }

                    &:hover {
                        .remove {
                            opacity: 1;

                            &:hover {
                                color: rgb(199, 32, 32);
                            }
                        }
                    }

                    &.active {
                        background-color: #e6f7ff;

                        .remove {
                            opacity: 1;
                        }
                    }
                }
            }

            .chapter-content {
                flex: 1;
                border-left: 1px solid #ccc;

                .content {
                    width: 100%;
                    padding: 10px;
                    box-sizing: border-box;
                    font-size: 14px;
                    text-align: left;
                    color: #333;
                    word-break: break-all;
                }
            }
        }
    }

}
</style>
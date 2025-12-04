<script lang="ts" setup>
import { ElMessageBox } from 'element-plus'
import { Window } from "@wailsio/runtime"
const handleClose = () => {
    ElMessageBox.confirm('确认要退出了吗？', { title: '退出程序', confirmButtonText: '退出', cancelButtonText: '取消', closeOnClickModal: false }).then(() => {
        Window.Close()
    }).catch(() => {
    })
}
const handleMinimise = () => {
    Window.Minimise()
}
</script>
<template>
    <div class="title-bar">
        <div class="title">百度贴吧小说下载器</div>
        <div class="actions">
            <div @click="handleClose" class="button close">&#xe639;</div>
            <div @click="handleMinimise" class="button minimise">&#xe67a;</div>
            <div @click="" class="button nop"></div>
        </div>
    </div>
</template>
<style lang="scss" scoped>
.title-bar {
    width: 100%;
    height: 30px;
    --wails-draggable: drag;
    cursor: default;
    display: flex;
    align-items: center;
    background-color: #333;

    .title {
        flex: 1;
        color: white;
        font-size: 14px;
    }

    .actions {
        display: flex;
        gap: 6px;
        position: fixed;
        left: 10px;

        &:hover {
            .button {
                &.close {
                    color: rgb(143 28 19);

                }

                &.minimise {
                    color: rgb(124, 107, 31);
                }
            }
        }

        .button {
            --wails-draggable: none;
            border-radius: 50%;
            font-size: 9px;
            line-height: 10px;
            font-weight: bold;
            cursor: pointer;
            display: flex;
            width: 14px;
            height: 14px;
            justify-content: center;
            align-items: center;
            padding-top: 1px;
            box-sizing: border-box;
            transition: all 0.1s ease-in-out;

            &.close {
                background-color: rgb(237, 106, 94);
                color: rgb(237, 106, 94);

                &:active {
                    background-color: rgb(245, 140, 130);
                }
            }

            &.minimise {
                background-color: rgb(244, 190, 79);
                color: rgb(244, 190, 79);

                &:active {
                    background-color: rgb(255, 205, 131);
                }
            }

            &.nop {
                cursor: default;
                background-color: rgb(78, 78, 78);
            }
        }
    }
}
</style>
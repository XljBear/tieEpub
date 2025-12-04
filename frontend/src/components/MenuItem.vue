<script lang="ts" setup>
const props = defineProps({
    label: {
        type: String,
        default: ''
    },
    icon: {
        type: String,
        default: ''
    },
    active: {
        type: Boolean,
        default: false
    }
})
</script>
<template>
    <div :class="['menu-item', props.active ? 'active' : '']">
        <div class="icon">{{ props.icon }}</div>
        <div class="label">{{ props.label }}</div>
    </div>
</template>
<style lang="scss" scoped>
.menu-item {
    width: 100%;
    height: 40px;
    margin: 10px 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #666;
    cursor: pointer;
    transition: background-color 0.3s ease;
    border-radius: 14px;
    font-size: 17px;
    display: flex;
    justify-content: start;
    padding: 10px 20px;
    box-sizing: border-box;

    div {
        cursor: pointer;
    }

    .icon {
        margin-right: 10px;
    }

    &:hover {
        background-color: #ddd;
    }

    &.active {
        background-color: rgb(240, 249, 253);
        color: rgb(73, 165, 248);
    }
}
</style>
<script setup>
import { state } from '@/state';

const hours = [
    { name: '近 3 小时', value: '3' },
    { name: '近 6 小时', value: '6' },
    { name: '近 12 小时', value: '12' },
    { name: '近 1 天', value: '24' },
    { name: '近 3 天', value: '72' },
    { name: '近 1 周', value: '168' },
    { name: '时间不限', value: 'x' },
];

function switchTime(event) {
    window.location.href = '/?hour=' + event.target.value;
}
</script>

<template>
    <select class="form-select bg-dark text-white" @change="switchTime">
        <option v-for="opt in hours" :selected="opt.value == state.hour" :value="opt.value">{{ opt.name }}</option>
    </select>
</template>

<style scoped>
select {
    margin-right: 16px;
    background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23ffffff' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='m2 5 6 6 6-6'/%3e%3c/svg%3e");
}
</style>
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

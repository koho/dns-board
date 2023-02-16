<script setup>
import { ref, watch } from 'vue';
import { state } from '@/state';

const total = ref(0);

watch(() => state.ipCount, val => {
    val.forEach(e => {
        total.value += e.count;
    });
})
</script>

<template>
    <div style="overflow-y: auto;">
        <table class="table">
            <thead>
                <tr>
                    <th scope="col">#</th>
                    <th scope="col">IP</th>
                    <th scope="col">请求数</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(row, i) in state.ipCount" :key="'client-row-' + i">
                    <th scope="row">{{ i + 1 }}</th>
                    <td>{{ row.ip }}</td>
                    <td>{{ row.count }} / {{ Math.round((row.count / total) * 100) }}%</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<style scoped>

</style>
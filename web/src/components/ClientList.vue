<script setup>
import { onMounted, ref } from 'vue';
import axios from 'axios';
import { state, getTokenHeader } from '@/state';

const clientList = ref([]);
const total = ref(0);

onMounted(function () {
    axios.get('/api/client?hour=' + state.hour, getTokenHeader())
        .then(function (response) {
            clientList.value = response.data.sort((a, b) => b.count - a.count).map(e => {
                total.value += e.count;
                return e;
            });
        })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .finally(function () {
            // always executed
        });
});
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
                <tr v-for="(row, i) in clientList" :key="'client-row-' + i">
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
<script setup>
import { watch, ref } from 'vue';
import axios from 'axios';
import { state } from '@/state'

const total = ref(0);

watch(() => state.domainList, val => {
    val.forEach(e => {
        total.value += e.count;
    });
    getGeoLocation();
})

function getGeoLocation() {
    let c = parseInt(state.domainList.length / 100);
    if (state.domainList.length % 100 > 0) {
        c++;
    }
    for (let i = 0; i < state.domainList.length; i += 100) {
        axios.post('http://ip-api.com/batch?lang=zh-CN', state.domainList.slice(i, i + 100).map(x => x.ip)).then(resp => {
            resp.data.forEach((e, j) => {
                if (e.status != 'success')
                    return;
                state.domainList[i + j].lon = e.lon;
                state.domainList[i + j].lat = e.lat;
                state.domainList[i + j].loc = [e.country, e.regionName, e.city].join(' ');
            });
        }).catch(error => {
            // handle error
            console.log(error);
        }).finally(() => {
            if (--c == 0) {
                state.geoLoaded = true;
            }
        });
    }
}
</script>

<template>
    <div style="overflow-y: auto;">
        <table class="table">
            <thead>
                <tr>
                    <th scope="col">#</th>
                    <th scope="col">域名</th>
                    <th scope="col">IP</th>
                    <th scope="col">位置</th>
                    <th scope="col">请求数</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(row, i) in state.domainList" :key="'domain-row-' + i">
                    <th scope="row">{{ i + 1 }}</th>
                    <td>{{ row.name }}</td>
                    <td>{{ row.ip }}</td>
                    <td>{{ row.loc }}</td>
                    <td>{{ row.count }} / {{ Math.round((row.count / total) * 100) }}%</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<style scoped>

</style>
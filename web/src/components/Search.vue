<script setup>
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { getTokenHeader } from '@/state';
import csvDownload from 'json-to-csv-export';

const props = defineProps({
    id: {
        type: String,
        required: true
    }
});

const typeMap = {
    6: "CLIENT",
    8: "FORWARDER"
};

const colNames = [
    "#", "时间", "类型", "IP", "端口", "协议", "域名", "请求类型", "应答", "响应码", "耗时"
]

var resultList = ref([]);
var kwInput = ref("");

function onEnter(e) {
    axios.get("/api/search?kw=" + e.target.value, getTokenHeader()).then(resp => {
        resultList.value = resp.data;
    }).catch(error => {
        console.log(error);
    });
}

onMounted(function () {
    document.getElementById(props.id).addEventListener('hidden.bs.modal', function (e) {
        resultList.value = [];
        kwInput.value = "";
    });
});

function downloadCSV() {
    csvDownload({
        data: resultList.value,
        filename: kwInput.value,
        delimiter: ',',
        headers: colNames
    });
}
</script>

<template>
    <div class="modal" tabindex="-1" :id="id">
        <div class="modal-dialog modal-xl">
            <div class="modal-content bg-dark text-white">
                <div class="modal-header">
                    <form>
                        <label class="search-icon" for="search-input" id="search-label">
                            <svg width="20" height="20" viewBox="0 0 20 20">
                                <path
                                    d="M14.386 14.386l4.0877 4.0877-4.0877-4.0877c-2.9418 2.9419-7.7115 2.9419-10.6533 0-2.9419-2.9418-2.9419-7.7115 0-10.6533 2.9418-2.9419 7.7115-2.9419 10.6533 0 2.9419 2.9418 2.9419 7.7115 0 10.6533z"
                                    stroke="currentColor" fill="none" fill-rule="evenodd" stroke-linecap="round"
                                    stroke-linejoin="round"></path>
                            </svg>
                        </label>
                        <input aria-labelledby="search-label" id="search-input" autocomplete="off" autocorrect="off"
                            autocapitalize="off" enterkeyhint="search" spellcheck="false" autofocus="true"
                            placeholder="搜索域名或 IP" v-on:keyup.enter="onEnter" v-model="kwInput">
                    </form>
                    <button type="button" class="btn-close btn-close-white" :disabled="kwInput == ''"
                        @click="kwInput = ''"></button>
                </div>
                <div class="modal-body">
                    <div style="overflow-y: auto; max-height: 1000px;">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th scope="col" v-for="col in colNames">{{ col }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(row, i) in resultList">
                                    <th scope="row">{{ i + 1 }}</th>
                                    <td>{{ row.time }}</td>
                                    <td>{{ typeMap[row.type] }}</td>
                                    <td>{{ row.ip }}</td>
                                    <td>{{ row.port }}</td>
                                    <td>{{ row.protocol }}</td>
                                    <td>{{ row.domain }}</td>
                                    <td>{{ row.qtype }}</td>
                                    <td>{{ row.answer }}</td>
                                    <td>{{ row.rcode }}</td>
                                    <td>{{ row.duration }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-primary" @click="downloadCSV">下载</button>
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">关闭</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
form {
    align-items: center;
    box-shadow: inset 0 0 0 2px var(--docsearch-primary-color);
    display: flex;
    height: 56px;
    margin: 0;
    padding: 0 12px;
    position: relative;
    width: 100%;
}

input {
    background: transparent;
    border: 0;
    color: rgb(204, 204, 220);
    flex: 1;
    font: inherit;
    font-size: 1.2em;
    height: 100%;
    outline: none;
    padding: 0 0 0 8px;
    width: 80%;
}

.search-icon {
    align-items: center;
    color: #712cf9;
    display: flex;
    justify-content: center;
    margin: 0;
    padding: 0;
}

.search-icon svg {
    height: 24px;
    width: 24px;
    box-sizing: border-box;
    stroke-width: 1.6;
    vertical-align: middle;
}

.modal-header {
    padding: 0;
    padding-right: 16px;
    border-bottom: var(--bs-modal-header-border-width) solid #712cf9;
}

.modal-footer {
    border-top: none;
}
</style>
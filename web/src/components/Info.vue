<script setup>
import { onMounted, ref } from 'vue';
import axios from 'axios';
import humanizeDuration from 'humanize-duration';
import { getTokenHeader } from '@/state';
import { Modal } from 'bootstrap/dist/js/bootstrap';

var identity = ref("");
var version = ref("");
var uptime = ref("");
var searchDialog = null;

onMounted(function () {
    axios.get("/api/meta", getTokenHeader()).then(resp => {
        identity.value = resp.data.identity;
        version.value = resp.data.version;
        uptime.value = humanizeDuration(Date.now() - resp.data.startup, { language: "zh_CN", largest: 2 });
    }).catch(error => {
        console.log(error);
        if (error.response && error.response.status == 401) {
            logout();
        }
    });
    searchDialog = new Modal(document.getElementById('search'));
});

function logout() {
    localStorage.clear();
    window.location.href = '/login';
}

</script>

<template>
    <div id="navbarNavDarkDropdown">
        <ul class="navbar-nav">
            <li class="nav-item dropdown">
                <button class="btn btn-dark dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
                    {{ identity }}
                </button>
                <ul class="dropdown-menu dropdown-menu-end dropdown-menu-dark">
                    <li><span class="dropdown-item read-text">{{ version }}</span></li>
                    <li><span class="dropdown-item read-text">{{ uptime }}</span></li>
                    <li>
                        <hr class="dropdown-divider">
                    </li>
                    <li><a class="dropdown-item" href="#" v-on:click.prevent="searchDialog.show()">搜索</a></li>
                    <li><a class="dropdown-item" href="#" v-on:click.prevent="logout">退出登录</a></li>
                </ul>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.read-text {
    pointer-events: none;
    cursor: default;
}
</style>

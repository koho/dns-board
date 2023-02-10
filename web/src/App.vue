<script setup>
import { state } from '@/state';
import Card from './components/Card.vue';
import ClientList from './components/ClientList.vue';
import DomainList from './components/DomainList.vue';
import GeoMap from './components/GeoMap.vue'
import Info from './components/Info.vue';
import Latency from './components/Latency.vue';
import QType from './components/QType.vue';
import LineChart from './components/LineChart.vue';
import SwitchTime from './components/SwitchTime.vue';
import Search from './components/Search.vue';
import RCode from './components/RCode.vue';

var hour = new URLSearchParams(window.location.search).get('hour');
if (!hour) {
  hour = '3';
}
state.hour = hour;

if (!localStorage.getItem('token')) {
  window.location.href = "/login";
}

function reload() {
  window.location.reload();
}
</script>

<template>
  <nav class="navbar navbar-expand-lg fixed-top navbar-dark bg-dark">
    <div class="container-fluid">
      <a class="navbar-brand" href="#" @click="reload">
        <img src="/logo.svg" alt="" width="30" height="24" class="d-inline-block align-text-top">
        DNS Board
      </a>
      <form class="d-flex">
        <SwitchTime></SwitchTime>
        <Info></Info>
      </form>
    </div>
  </nav>
  <div class="container-fluid main">
    <div class="row">
      <div class="col">
        <Card title="每分钟请求数">
          <LineChart id="query-count" url="/api/count" text="每分钟请求数"></LineChart>
        </Card>
      </div>
      <div class="col">
        <Card title="请求耗时趋势（毫秒）">
          <Latency id="query-duration"></Latency>
        </Card>
      </div>
      <div class="col">
        <Card title="响应状态">
          <RCode id="rcode"></RCode>
        </Card>
      </div>
      <div class="col">
        <Card title="客户端排行">
          <ClientList></ClientList>
        </Card>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <Card title="域名位置分布">
          <GeoMap id="map"></GeoMap>
        </Card>
      </div>
      <div class="col">
        <Card title="请求域名排行">
          <DomainList></DomainList>
        </Card>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <Card title="请求类型占比">
          <QType id="qtype"></QType>
        </Card>
      </div>
      <div class="col">
        <Card title="响应大小趋势（字节）">
          <LineChart id="resp-size" url="/api/size" text="响应大小"></LineChart>
        </Card>
      </div>
      <div class="col">
        <Card title="缓存命中率">
          <LineChart id="cache-hit" url="/api/cache" text="命中率"></LineChart>
        </Card>
      </div>
    </div>
  </div>
  <Search id="search"></Search>
</template>

<style scoped>
.main {
  margin-top: 72px;
}

.container-fluid .row:first-child .card-body div:first-child {
  height: 300px;
}

.container-fluid .row:nth-child(2) .card-body div:first-child {
  height: 760px;
}

.col {
  flex: 1 0 600px;
}

nav {
  box-shadow: 0 1px 6px 0 rgb(0 0 0 / 20%);
}
</style>

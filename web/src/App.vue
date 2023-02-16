<script setup>
import { state, getTokenHeader } from '@/state';
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
import { onMounted } from 'vue';
import axios from 'axios';

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

function getIP(domain) {
  if (domain.ipv4.length > 0) {
    return domain.ipv4[0];
  } else if (domain.ipv6.length > 0) {
    return domain.ipv6[0];
  }
  return "";
}

const rcode2String = {
  0: "NOERROR",
  1: "FORMERR",
  2: "SERVFAIL",
  3: "NXDOMAIN",
  4: "NOTIMP",
  5: "REFUSED",
  6: "YXDOMAIN", // See RFC 2136
  7: "YXRRSET",
  8: "NXRRSET",
  9: "NOTAUTH",
  10: "NOTZONE",
  16: "BADSIG", // Also known as RcodeBadVers, see RFC 6891
  //	RcodeBadVers:        "BADVERS",
  17: "BADKEY",
  18: "BADTIME",
  19: "BADMODE",
  20: "BADNAME",
  21: "BADALG",
  22: "BADTRUNC",
  23: "BADCOOKIE",
}

function stat(data) {
  const rcodeMap = {};
  const timeGroup = data.reduce((group, row) => {
    const time = row.time.substr(0, 16).replace('T', ' ');
    group[time] = group[time] ?? { clientCount: 0, forwarderCount: 0, cacheHit: 0, cacheRate: 0, avgSize: 0, forwarders: {}, rcode: {} };
    if (row.type == 6) {
      group[time].clientCount++;
      if (row.duration == 0) {
        group[time].cacheHit++;
      }
      group[time].cacheRate = group[time].cacheHit / group[time].clientCount;
      group[time].avgSize -= (group[time].avgSize - row.size) / group[time].clientCount;
      group[time].rcode[row.rcode] = group[time].rcode[row.rcode] ?? 0;
      group[time].rcode[row.rcode]++;
      rcodeMap[row.rcode] = 0;
    } else if (row.type == 8) {
      group[time].forwarderCount++;
      group[time].forwarders[row.ip] = group[time].forwarders[row.ip] ?? { count: 0, avg: 0 };
      group[time].forwarders[row.ip].count++;
      group[time].forwarders[row.ip].avg -= (group[time].forwarders[row.ip].avg - row.duration) / group[time].forwarders[row.ip].count;
    }
    return group;
  }, {});
  state.queryCount = Object.entries(timeGroup).map(([time, obj]) => ({ time, value: obj.clientCount }));
  const durationData = [];
  Object.entries(timeGroup).forEach(([time, obj]) => {
    Object.entries(obj.forwarders).forEach(([forwarder, f]) => {
      durationData.push({ time, type: forwarder, value: Math.round(f.avg) });
    })
  });
  state.durationData = durationData;
  const rcodeData = [];
  Object.entries(timeGroup).forEach(([time, obj]) => {
    Object.entries(rcodeMap).forEach(([rcode, count]) => {
      rcodeData.push({ time, type: rcode2String[rcode], value: obj.rcode[rcode] ?? count });
    })
  });
  state.rcodeData = rcodeData;
  state.sizeData = Object.entries(timeGroup).map(([time, obj]) => ({ time, value: Math.round(obj.avgSize) }));
  state.cacheData = Object.entries(timeGroup).map(([time, obj]) => ({ time, value: Math.floor(obj.cacheRate * 100) / 100 }));

  const ipGroup = data.reduce((group, row) => {
    if (row.type == 6 && row.answer !== "") {
      group[row.ip] = group[row.ip] ?? 0;
      group[row.ip]++;
    }
    return group;
  }, {});
  state.ipCount = Object.entries(ipGroup).map(([ip, obj]) => ({ ip, count: obj })).sort((a, b) => b.count - a.count);

  const qtypeGroup = data.reduce((group, row) => {
    if (row.type == 6 && row.answer !== "") {
      group[row.qtype] = group[row.qtype] ?? 0;
      group[row.qtype]++;
    }
    return group;
  }, {});
  state.qtypeCount = Object.entries(qtypeGroup).map(([type, obj]) => ({ type, value: obj })).sort((a, b) => b.value - a.value);

  const domainGroup = data.reduce((group, row) => {
    if (row.type == 6 && row.answer !== "" && (row.qtype === "A" || row.qtype === "AAAA")) {
      group[row.domain] = group[row.domain] ?? { ipv4: new Set(), ipv6: new Set(), count: 0 };
      const newIP = row.answer.split(',');
      if (newIP.length > 0)
        group[row.domain].count++;
      if (row.qtype === "A") {
        for (let elem of newIP) {
          group[row.domain].ipv4.add(elem);
        }
      } else if (row.qtype === "AAAA") {
        for (let elem of newIP) {
          group[row.domain].ipv6.add(elem);
        }
      }
    }
    return group;
  }, {});
  state.domainList = Object.entries(domainGroup).map(([name, obj]) => ({ name, ipv4: Array.from(obj.ipv4), ipv6: Array.from(obj.ipv6), count: obj.count })).sort((a, b) => b.count - a.count).map(e => {
    return { ...e, ip: getIP(e) }
  });
}

onMounted(function () {
  axios.get('/api/raw?hour=' + state.hour, getTokenHeader()).then(resp => {
    stat(resp.data);
  }).catch(function (error) {
    // handle error
    console.log(error);
  })
})
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
          <LineChart id="query-count" :data="state.queryCount" text="每分钟请求数"></LineChart>
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
          <LineChart id="resp-size" :data="state.sizeData" text="响应大小"></LineChart>
        </Card>
      </div>
      <div class="col">
        <Card title="缓存命中率">
          <LineChart id="cache-hit" :data="state.cacheData" text="命中率"></LineChart>
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

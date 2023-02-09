<script setup>
import { onMounted } from 'vue';
import { Heatmap, G2 } from '@antv/g2plot';
import axios from 'axios';
import { state, getTokenHeader } from '@/state';

const props = defineProps({
    id: {
        type: String,
        required: true
    }
});

const theme = G2.getTheme('dark');
G2.registerTheme('transparent-dark', {
    ...theme,
    background: 'transparent'
});

onMounted(function () {
    axios.get('/api/rcode?hour=' + state.hour, getTokenHeader())
        .then(function (response) {
            const cfg = {
                data: response.data,
                theme: 'transparent-dark',
                xField: 'time',
                yField: 'type',
                colorField: 'value',
                legend: {
                    rail: {
                        defaultLength: 200
                    }
                },
                color: ['#047331', '#388C04', '#CACE17', '#E16519', '#CA0300'],
                meta: {
                    'value': {
                        alias: '计数'
                    }
                },
                xAxis: {
                    tickCount: 5,
                },
            }
            const plot = new Heatmap(props.id, cfg);
            plot.render();
        })
        .catch(function (error) {
            // handle error
            console.log(error);
        })
        .finally(function () {
            // always executed
        });
})
</script>

<template>
    <div :id="id"></div>
</template>

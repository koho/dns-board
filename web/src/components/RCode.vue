<script setup>
import { watch } from 'vue';
import { Heatmap, G2 } from '@antv/g2plot';
import { state } from '@/state';

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

watch(() => state.rcodeData, val => {
    const cfg = {
        data: val,
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
</script>

<template>
    <div :id="id"></div>
</template>

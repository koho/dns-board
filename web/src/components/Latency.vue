<script setup>
import { watch } from 'vue';
import { Line, G2 } from '@antv/g2plot';
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

watch(() => state.durationData, val => {
    const line = new Line(props.id, {
        padding: 'auto',
        forceFit: true,
        data: val,
        theme: 'transparent-dark',
        xField: 'time',
        yField: 'value',
        seriesField: 'type',
        xAxis: {
            tickCount: 5,
        },
        legend: {
            position: 'bottom',
            itemName: {
                style: {
                    fontSize: 20
                }
            }
        },
        smooth: true,
        area: {
            style: {
                fillOpacity: 0.15,
            },
        },
    });
    line.render();
})
</script>

<template>
    <div :id="id"></div>
</template>

<style scoped>

</style>
<script setup>
import { watch } from 'vue';
import { Line, G2 } from '@antv/g2plot';

const props = defineProps({
    id: {
        type: String,
        required: true
    },
    data: {
        required: true
    },
    text: {
        type: String,
        required: true
    }
});

const theme = G2.getTheme('dark');
G2.registerTheme('transparent-dark', {
    ...theme,
    background: 'transparent'
});

watch(() => props.data, val => {
    const line = new Line(props.id, {
        padding: 'auto',
        forceFit: true,
        data: val,
        theme: 'transparent-dark',
        xField: 'time',
        yField: 'value',
        xAxis: {
            tickCount: 5,
        },
        smooth: true,
        area: {
            style: {
                fillOpacity: 0.15,
            },
        },
        meta: {
            'value': {
                alias: props.text
            }
        }
    });
    line.render();
})
</script>

<template>
    <div :id="id"></div>
</template>

<style scoped>

</style>
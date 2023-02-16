<script setup>
import { watch } from 'vue';
import { Pie, G2 } from '@antv/g2plot';
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

watch(() => state.qtypeCount, val => {
    const cfg = {
        appendPadding: 10,
        data: val,
        theme: 'transparent-dark',
        angleField: 'value',
        colorField: 'type',
        radius: 0.9,
        label: {
            type: 'inner',
            offset: '-30%',
            content: ({ percent }) => `${(percent * 100).toFixed(0)}%`,
            style: {
                fontSize: 14,
                textAlign: 'center',
            },
        },
        legend: {
            itemName: {
                style: {
                    fontSize: 16
                }
            }
        },
        interactions: [{ type: 'element-active' }],
    }
    const piePlot = new Pie(props.id, cfg);
    piePlot.render();
})
</script>

<template>
    <div :id="id"></div>
</template>
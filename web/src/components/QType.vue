<script setup>
import { onMounted } from 'vue';
import { Pie, G2 } from '@antv/g2plot';
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
    axios.get('/api/qtype?hour=' + state.hour, getTokenHeader())
        .then(function (response) {
            const G = G2.getEngine('canvas');
            const cfg = {
                appendPadding: 10,
                data: response.data,
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
                interactions: [{ type: 'element-active' }],
            }
            const piePlot = new Pie(props.id, cfg);
            piePlot.render();
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
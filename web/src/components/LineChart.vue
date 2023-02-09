<script setup>
import { onMounted } from 'vue';
import { Line, G2 } from '@antv/g2plot';
import axios from 'axios';
import { state, getTokenHeader } from '@/state';

const props = defineProps({
    id: {
        type: String,
        required: true
    },
    url: {
        type: String,
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

onMounted(function () {
    axios.get(props.url + '?hour=' + state.hour, getTokenHeader())
        .then(function (response) {
            const line = new Line(props.id, {
                padding: 'auto',
                forceFit: true,
                data: response.data,
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

<style scoped>

</style>
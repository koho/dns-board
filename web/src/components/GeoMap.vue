<script setup>
import { state, getTokenHeader } from '@/state';
import { Map, NavigationControl, Popup } from 'maplibre-gl'
import { watch } from 'vue';
import axios from 'axios';

const props = defineProps({
    id: {
        type: String,
        required: true
    }
});

watch(() => state.geoLoaded, val => {
    const geojson = { type: 'FeatureCollection' }
    geojson.features = state.domainList.map(e => {
        return {
            type: 'Feature',
            properties: {
                domain: e.name,
                ip: e.ip
            },
            geometry: {
                type: 'Point',
                coordinates: [e.lon, e.lat]
            }
        }
    });
    axios.get("/api/meta", getTokenHeader()).then(resp => {
        loadMap(resp.data.mapUrl, geojson);
    }).catch(error => {
        console.log(error);
    })
});

function loadMap(url, data) {
    const map = new Map({
        container: props.id, // container's id or the HTML element in which MapLibre GL JS will render the map
        style: url, // style URL
        center: [-174.89, 38.15], // starting position [lng, lat]
        zoom: 2, // starting zoom
    });
    map.addControl(new NavigationControl(), 'top-left');
    map.on('load', function () {
        map.addSource('dns', {
            type: 'geojson',
            data,
            cluster: true,
            clusterMaxZoom: 14, // Max zoom to cluster points on
            clusterRadius: 50
        });
        map.addLayer({
            id: 'domain',
            type: 'circle',
            source: 'dns',
            paint: {
                // Use step expressions (https://maplibre.org/maplibre-gl-js-docs/style-spec/#expressions-step)
                // with three steps to implement three types of circles:
                //   * Blue, 20px circles when point count is less than 100
                //   * Yellow, 30px circles when point count is between 100 and 750
                //   * Pink, 40px circles when point count is greater than or equal to 750
                'circle-color': ['step', ['get', 'point_count'], '#5794f2', 100, '#b877d9', 750, '#f2495c'],
                'circle-radius': ['step', ['get', 'point_count'], 20, 100, 30, 750, 40],
                'circle-opacity': 0.3,
                'circle-stroke-color': ['step', ['get', 'point_count'], '#5794f2', 100, '#b877d9', 750, '#f2495c'],
                'circle-stroke-width': 1
            }
        });
        map.addLayer({
            id: 'domain-count',
            type: 'symbol',
            source: 'dns',
            filter: ['has', 'point_count'],
            layout: {
                'text-field': '{point_count_abbreviated}',
                'text-font': ['DIN Offc Pro Medium', 'Arial Unicode MS Bold'],
                'text-size': 12
            },
            paint: {
                'text-color': "#ffffff"
            }
        });
        map.addLayer({
            id: 'domain-point',
            type: 'circle',
            source: 'dns',
            filter: ['!', ['has', 'point_count']],
            paint: {
                'circle-color': '#5794f2',
                'circle-radius': 4,
                'circle-stroke-width': 1,
                'circle-opacity': 0.3,
                'circle-stroke-color': '#5794f2'
            }
        });
        // inspect a cluster on click
        map.on('click', 'domain', function (e) {
            var features = map.queryRenderedFeatures(e.point, {
                layers: ['domain']
            });
            var clusterId = features[0].properties.cluster_id;
            map.getSource('dns').getClusterExpansionZoom(clusterId, function (err, zoom) {
                if (err) return;
                if (!zoom) return;
                map.easeTo({
                    center: features[0].geometry.coordinates,
                    zoom: zoom
                });
            });
        });
        map.on('click', 'domain-point', function (e) {
            var coordinates = e.features[0].geometry.coordinates.slice();
            var domain = e.features[0].properties.domain;
            var ip = e.features[0].properties.ip;
            // Ensure that if the map is zoomed out such that
            // multiple copies of the feature are visible, the
            // popup appears over the copy being pointed to.
            while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
                coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
            }

            new Popup().setLngLat(coordinates).setHTML(
                domain + '<br>' + ip
            ).addTo(map);
        });
        map.on('mouseenter', 'domain', function () {
            map.getCanvas().style.cursor = 'pointer';
        });
        map.on('mouseleave', 'domain', function () {
            map.getCanvas().style.cursor = '';
        });
    });
}
</script>

<template>
    <div :id="id"></div>
</template>

<style scoped>

</style>

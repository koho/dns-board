import { reactive } from 'vue'

export const state = reactive({
    domainList: [],
    geoLoaded: false,
    hour: '',
    queryCount: [],
    durationData: [],
    rcodeData: [],
    sizeData: [],
    cacheData: [],
    ipCount: [],
    qtypeCount: [],
})

export function getTokenHeader() {
    return {
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    }
}

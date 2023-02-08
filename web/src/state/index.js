import { reactive } from 'vue'

export const state = reactive({
    domainList: [],
    geoLoaded: false,
    hour: '',
})

export function getTokenHeader() {
    return {
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
    }
}

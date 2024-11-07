import {ref} from 'vue'
import {apply} from 'json-merge-patch'

const state = ref({})

export function useState() {
    function patch(message) {
        state.value = apply(state.value, message)
    }

    function set(message) {
        state.value = message
    }

    return {
        state,
        patch,
        set
    }
}
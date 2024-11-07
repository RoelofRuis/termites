import {ref} from 'vue'

const messages = ref([])

export function useMessageStream(key) {
    function prepend(message) {
        messages.value.unshift(message)
    }

    function clear() {
        messages.value = []
    }

    return {
        messages,
        prepend,
        clear,
    }
}
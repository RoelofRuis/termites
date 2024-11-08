import {ref} from 'vue'

const messages = ref([])

export function useMessageStream() {
    function prepend(message) {
        message.index = messages.value.length
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
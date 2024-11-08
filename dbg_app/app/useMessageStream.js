import {ref} from 'vue'

const messages = ref([])
let index = 0

export function useMessageStream() {
    let lastPrependTime = Date.now()

    function prepend(message) {
        message.index = index++
        const newPrependTime = Date.now()
        if (newPrependTime - lastPrependTime > 500) {
            message.group_start = true
        }
        messages.value.unshift(message)
        lastPrependTime = newPrependTime
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
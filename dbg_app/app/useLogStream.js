import {ref} from 'vue'

const logs = ref([])

export function useLogStream() {
    function prepend(log) {
        log.index = logs.value.length
        logs.value.unshift(log)
    }

    function clear() {
        logs.value = []
    }

    return {
        logs,
        prepend,
        clear,
    }
}
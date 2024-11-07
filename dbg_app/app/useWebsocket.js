const topicSubscribers = []
let ws

export function useWebsocket() {
    const subscribe = (topic, callback) => {
        const regex = new RegExp('^' + topic.replace('*', '.*'))
        topicSubscribers.push({regex, callback})
    }

    const unsubscribe = (topic, callback) => {
        const regex = new RegExp('^' + topic.replace('*', '.*'))

        for (let i = 0; i < topicSubscribers.length; i++) {
            if (topicSubscribers[i].regex.toString() === regex.toString() && topicSubscribers[i].callback === callback) {
                topicSubscribers.splice(i, 1)
                break;
            }
        }
    }

    const open = (wsUrl) => {
        ws = new WebSocket(wsUrl);

        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data)
                if (message.topic) {
                    notifySubscribers(message.topic, message.payload)
                }
            } catch (err) {
                console.warn("failed to parse message: ", event.data, err)
            }
        }
    }

    const notifySubscribers = (topic, message) => {
        topicSubscribers.forEach(({regex, callback}) => {
            if (regex.test(topic)) {
                callback(message)
            }
        })
    }

    return {
        subscribe,
        unsubscribe,
        open,
    }
}
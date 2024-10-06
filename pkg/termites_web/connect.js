let storage = (function () {

    let methods = {};

    methods.hasStorage = function () {
        return typeof (Storage) !== "undefined"
    }

    methods.put = function (key, value) {
        if (!methods.hasStorage()) {
            return
        }
        window.localStorage.setItem(key, value);
    }

    methods.get = function (key) {
        if (!methods.hasStorage()) {
            return
        }
        return window.localStorage.getItem(key)
    }

    return methods;
})();

let connector = (function (storage) {
    let conn;

    let methods = {};

    methods.connect = function () {
        if (!window["WebSocket"]) {
            return false
        }
        let id = storage.get("id")
        let url = "ws://" + document.location.host + "/ws"
        if (id !== null) {
            url = url + "?id=" + id
        }

        conn = new WebSocket(url);

        conn.onopen = function (evt) {
            publish("onopen", {})
        }
        conn.onclose = function (evt) {
            publish("onclose", {})
        }
        conn.onmessage = function (evt) {
            const msg = JSON.parse(evt.data);
            onmessage(msg);
        }

        return true
    }

    function onmessage(msg) {
        if (msg.msg_type === "update") {
            publish("onupdate", msg);
            return;
        }

        if (msg.msg_type === "_connected") { // tells which id is linked to this client
            let id = msg.payload.id;
            storage.put("id", id);
            return;
        }

        if (msg.msg_type === "_close") {
            window.close();
        }
    }

    const subscriptions = [];

    function publish(event, data) {
        subscriptions.forEach((callback) => {
            callback(event, data)
        })
    }

    methods.subscribe = function (callback) {
        subscriptions.push(callback)
    }

    methods.send = function (type, object) {
        if (!conn || conn.readyState !== WebSocket.OPEN) {
            return
        }

        let msg = {
            timestamp: Date.now(),
            type: type,
            data: object,
        }

        conn.send(JSON.stringify(msg))
    }

    return methods
})(storage);

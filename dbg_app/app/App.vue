<script setup>
import NavBar from "./NavBar.vue";
import {onMounted, ref} from "vue";
import {useWebsocket} from "./useWebsocket";
import {useState} from "./useState";
import GraphContainer from "./GraphContainer.vue";
import MessageContainer from "./MessageContainer.vue";
import {useMessageStream} from "./useMessageStream";
import LogsContainer from "./LogsContainer.vue";
import {useLogStream} from "./useLogStream";

const { open, subscribe } = useWebsocket()
const { patch, set } = useState()
const { prepend: prependMessage } = useMessageStream()
const { prepend: prependLog } = useLogStream()

const tab = ref('graph')

onMounted(() => {
  subscribe("state/patch", patch)
  subscribe("state/full", set)
  subscribe("message", prependMessage)
  subscribe("log", prependLog)
  open("ws://" + document.location.host + "/ws")
})
</script>

<template>
  <div class="app-container">
    <NavBar @tab-selected="selected => tab = selected"/>
    <GraphContainer v-show="tab === 'graph'"/>
    <MessageContainer v-show="tab === 'messages'"/>
    <LogsContainer v-show="tab === 'logs'"/>
  </div>
</template>

<style>
html,body {
  margin: 0;
  padding: 0;
  background: #fff;
  font-family: sans,serif;
}

* {
  box-sizing: border-box;
}
</style>
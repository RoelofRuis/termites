<script setup>
import NavBar from "./NavBar.vue";
import {onMounted, ref} from "vue";
import {useWebsocket} from "./useWebsocket";
import {useState} from "./useState";
import GraphContainer from "./GraphContainer.vue";
import MessageContainer from "./MessageContainer.vue";
import {useMessageStream} from "./useMessageStream";

const { open, subscribe } = useWebsocket()
const { patch, set } = useState()
const { prepend } = useMessageStream()

const tab = ref('graph')

onMounted(() => {
  subscribe("state/patch", patch)
  subscribe("state/full", set)
  subscribe("message", prepend)
  open("ws://" + document.location.host + "/ws")
})

</script>

<template>
  <div class="app-container">
    <NavBar @tab-selected="selected => tab = selected"/>
    <GraphContainer v-if="tab === 'graph'"/>
    <MessageContainer v-if="tab === 'messages'"/>
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
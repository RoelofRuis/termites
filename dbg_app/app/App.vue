<script setup>
import NavBar from "./NavBar.vue";
import PageContainer from "./PageContainer.vue";
import {onMounted} from "vue";
import {useWebsocket} from "./useWebsocket";
import {useState} from "./useState";

const { open, subscribe } = useWebsocket()
const { patch, set } = useState()

onMounted(() => {
  subscribe("state/patch", patch)
  subscribe("state/full", set)
  open("ws://" + document.location.host + "/ws")
})

</script>

<template>
  <div class="app-container">
    <NavBar/>
    <PageContainer/>
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
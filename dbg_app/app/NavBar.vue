<script setup>

import {useState} from "./useState";
import {computed} from "vue";

defineEmits(['tabSelected'])

const { state } = useState()

const graphEnabled = computed(() => {
  return state.value.debugger?.graph_enabled ?? false;
})

const logsEnabled = computed(() => {
  return state.value.debugger?.logs_enabled ?? false;
})

const messagesEnabled = computed(() => {
  return state.value.debugger?.messages_enabled ?? false;
})

</script>

<template>
<nav class="navbar">
  <div class="title header-item">Termites Debugger</div>
  <div class="nav-list">
    <div
        v-if="graphEnabled"
        class="nav-item header-item"
        @click="$emit('tabSelected', 'graph')"
    >Graph</div>
    <div
        v-if="logsEnabled"
        class="nav-item header-item"
        @click="$emit('tabSelected', 'logs')"
    >Logs</div>
    <div
        v-if="messagesEnabled"
        class="nav-item header-item"
        @click="$emit('tabSelected', 'messages')"
    >Messages</div>
  </div>
</nav>
</template>

<style scoped>
.navbar {
  display: flex;
  flex-direction: row;
  border-bottom: 1px solid black;
  background: #e3e0bf;
  user-select: none;
  padding: 0 20px;
  height: 40px;
}

.header-item {
  padding: 10px;
}

.nav-list {
  display: flex;
  margin-left: auto;
}

.nav-item {
  flex: 1;
  cursor: pointer;
}
</style>
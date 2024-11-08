<script setup>

import {useState} from "./useState";
import {computed} from "vue";

const { state } = useState()

const graphPath = computed(() => {
  return state.value.graph?.path ?? null
})

</script>

<template>
<div class="page">
  <div class="graph-container">
    <img class="graph" alt="graph" :src="graphPath"/>
  </div>
</div>
</template>

<style scoped>
.page {
  overflow: scroll;
}

.graph-container {
  display: flex;
  transform: rotateX(180deg);
  flex-direction: row;
  max-height: calc(100vh - 40px);
  overflow: scroll;
}

.graph-container * {
  transform: rotateX(180deg);
}

.graph {
  flex: 1;
}
</style>
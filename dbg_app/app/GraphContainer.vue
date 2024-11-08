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
  <img class="graph" alt="graph" :src="graphPath"/>
</div>
</template>

<style scoped>
.page {
  width: 100vw;
  height: calc(100vh - 40px);
  overflow: scroll;
}

</style>
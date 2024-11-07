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
  <img :src="graphPath"/>
</div>
</template>

<style scoped>
.page {
  overflow: scroll;
}
</style>
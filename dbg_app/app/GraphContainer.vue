<script setup>

import {useState} from "./useState";
import {computed, ref} from "vue";

const { state } = useState()

const graphPath = computed(() => {
  return state.value.graph?.path ?? null
})

const zoom = ref(1.0)

</script>

<template>
<div class="page">
  <div class="zoom zoom-in" @click="zoom += 0.1">+</div>
  <div class="zoom zoom-out" @click="zoom -= 0.1">−</div>
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
  scale: v-bind(zoom);
}

.zoom {
  z-index: 1;
  position: fixed;
  width: 40px;
  height:40px;
  line-height: 38px;
  font-size: 24px;
  text-align: center;
  background: rgba(0, 0, 0, 0.1);
  user-select: none;
}

.zoom:hover {
  background: rgba(0, 0, 0, 0.3);
}

.zoom-in {
  top: 50px;
  left: 10px;
  border-bottom-left-radius: 50px;
  border-top-left-radius: 50px;
}

.zoom-out {
  top: 50px;
  left: 50px;
  border-bottom-right-radius: 50px;
  border-top-right-radius: 50px;
}
</style>
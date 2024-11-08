<script setup>

import {useLogStream} from "./useLogStream";

const { logs } = useLogStream()

</script>

<template>
<div class="logs-page">
  <pre>{{logs}}</pre>
</div>
</template>

<style scoped>
.logs-page {
  font-family: monospace;
  display: flex;
}
</style>
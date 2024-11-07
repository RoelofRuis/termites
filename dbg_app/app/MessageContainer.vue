<script setup>

import {useMessageStream} from "./useMessageStream";

const { messages } = useMessageStream()

</script>

<template>
<div class="message-page">
  <div v-for="message in messages" class="message">
    {{message.from_name}} {{message.from_port_name}} -> {{message.to_name}} {{message.to_port_name}}
  </div>
</div>
</template>

<style scoped>
.message-page {
  font-family: monospace;
}

.message {
  border: 1px solid black;
  border-radius: 5px;
  margin: 2px;
  padding: 2px;
}
</style>
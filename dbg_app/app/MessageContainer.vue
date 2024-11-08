<script setup>

import {useMessageStream} from "./useMessageStream";
import {msToTime} from "./util";

const { messages, clear } = useMessageStream()

</script>

<template>
<div class="message-page">
  <div class="message-table">
    <table>
      <thead>
        <tr class="header">
          <th></th>
          <th>Time</th>
          <th>From</th>
          <th>To</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="message in messages">
          <th>{{message.index}}</th>
          <td>{{msToTime(message.time_since_opened_ms)}}</td>
          <td>{{message.from_name}} ({{message.from_port_name}})</td>
          <td>{{message.to_name}} ({{message.to_port_name}})</td>
        </tr>
      </tbody>
    </table>
  </div>
  <div class="message-controls">
    <span @click="clear">Clear</span>
  </div>
</div>
</template>

<style scoped>
.message-page {
  font-family: monospace;
  display: flex;
}

.message-table {
  flex: 1;
}

.message-controls {
  flex: 1;
}

</style>
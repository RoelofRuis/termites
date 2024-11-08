<script setup>

import {useLogStream} from "./useLogStream";
import {msToTime} from "./util";

const { logs, clear } = useLogStream()

</script>

<template>
<div class="logs-page">
  <div class="logs-table">
    <table>
      <thead>
        <tr class="header">
          <th></th>
          <th>Time</th>
          <th>Message</th>
        </tr>
      </thead>
      <tbody>
        <tr
            v-for="log in logs"
            :class="{info: log.log_level === 'info', error: log.log_level === 'error', panic: log.log_level === 'panic'}"
        >
          <th>{{log.index}}</th>
          <td>{{msToTime(log.time_since_opened_ms)}}</td>
          <td>{{log.message}}<br/>{{log.error}}{{log.stack}}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <div class="logs-controls">
    <span @click="clear">Clear</span>
  </div>
</div>
</template>

<style scoped>
.logs-page {
  font-family: monospace;
  display: flex;
}

.logs-table {
  flex: 1;
}

.logs-controls {
  flex: 1;
}

.info {
  background: #a2c4d6;
}

.error {
  background: #dfbb8f;
}

.panic {
  background: #dc9999;
}
</style>
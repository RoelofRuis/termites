<script setup>
import { msToTime } from './util';
import { useLogStream } from './useLogStream';

const { logs, clear } = useLogStream();

const logClass = (log_level) => {
  return {
    'log-entry': true,
    'log-info': log_level === 'info',
    'log-error': log_level === 'error',
    'log-panic': log_level === 'panic',
  };
};
</script>

<template>
  <div class="log-container">
    <button class="clear-button" @click="clear">Clear Logs</button>
    <div v-for="log in logs" :key="log.id" :class="logClass(log.log_level)">
      <div class="log-header">
        <span class="log-time">{{ msToTime(log.time_since_opened_ms) }}</span>
        <span class="log-type">{{ log.log_level.toUpperCase() }}</span>
      </div>
      <div class="log-message">{{ log.message }}</div>
      <div v-if="log.stack_track" class="log-stacktrace">
        <pre>{{ log.stack_track }}</pre>
      </div>
    </div>
  </div>
</template>

<style scoped>
.log-container {
  font-family: Arial, sans-serif;
  max-width: 800px;
  margin: auto;
  padding-top: 10px;
}

.clear-button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  font-size: 0.9em;
  cursor: pointer;
  margin-bottom: 10px;
  border-radius: 4px;
  width: 100%;
}

.clear-button:hover {
  background-color: #0056b3;
}

.log-entry {
  border: 1px solid #ddd;
  border-radius: 4px;
  margin: 10px 0;
  padding: 10px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  font-weight: bold;
}

.log-time {
  font-size: 0.9em;
  color: #666;
}

.log-type {
  padding: 2px 5px;
  border-radius: 3px;
}

.log-info .log-type {
  background-color: #d9edf7;
  color: #31708f;
}

.log-error .log-type {
  background-color: #f2dede;
  color: #a94442;
}

.log-panic .log-type {
  background-color: #fcf8e3;
  color: #8a6d3b;
}

.log-message {
  margin: 5px 0;
}

.log-stacktrace {
  background-color: #f7f7f9;
  padding: 10px;
  border-radius: 4px;
  font-size: 0.85em;
  color: #b94a48;
  overflow-x: auto;
}
</style>
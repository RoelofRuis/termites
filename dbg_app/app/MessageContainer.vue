<script setup>
import { msToTime } from './util';
import {useMessageStream} from "./useMessageStream";

const { messages, clear } = useMessageStream();

</script>

<template>
  <div class="message-stream-container">
    <button class="clear-button" @click="clear">Clear messages</button>
    <div v-for="message in messages" :key="message.index">
      <div class="message-row">
        <span class="message-index">{{ message.index }}</span>
        <span class="message-time">{{ msToTime(message.time_since_opened_ms) }}</span>
        <span class="message-from">{{ message.from_name }} ({{ message.from_port_name }})</span>
        <span class="message-arrow">→</span>
        <span class="message-to">{{ message.to_name }} ({{ message.to_port_name }})</span>
        <span class="message-error" v-if="message.error">{{ message.error }}</span>
      </div>
      <div v-if="message.group_start" class="spacer">---</div>
    </div>
  </div>
</template>

<style scoped>
.message-stream-container {
  padding-top: 10px;
  font-family: Arial, sans-serif;
  max-width: 800px;
  margin: auto;
}

.message-error {
  flex: 1;
  color: #dd3322;
  text-align: right;
  padding-right: 20px;
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

.message-row {
  display: flex;
  align-items: center;
  padding: 4px 0;
  font-size: 0.85em;
  border-bottom: 1px solid #ddd;
}

.message-index {
  width: 40px;
  text-align: right;
  margin-right: 20px;
  color: #666;
}

.message-time {
  width: 80px;
  margin-right: 20px;
  color: #888;
}

.message-from,
.message-to {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 250px;
}

.message-arrow {
  margin: 0 6px;
  color: #666;
}

.spacer {
  text-align: center;
  color: #aaa;
  font-size: 0.9em;
  padding: 4px 0;
  font-weight: bold;
}
</style>
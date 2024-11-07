<script setup>

defineEmits(['tabSelected'])

</script>

<template>
<nav class="navbar">
  <div class="title header-item">Termites Debugger</div>
  <div class="nav-list">
    <div
        class="nav-item header-item"
        @click="$emit('tabSelected', 'graph')"
    >Graph</div>
    <div
        class="nav-item header-item"
        @click="$emit('tabSelected', 'messages')"
    >Messages</div>
  </div>
</nav>
</template>

<style scoped>
.navbar {
  display: flex;
  flex-direction: row;
  border-bottom: 1px solid black;
  background: #e3e0bf;
  user-select: none;
  padding: 0 20px;
}

.header-item {
  padding: 10px;
}

.nav-list {
  display: flex;
  margin-left: auto;
}

.nav-item {
  flex: 1;
  cursor: pointer;
}
</style>
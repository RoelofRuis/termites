<script setup>

connector.connect();

const img = document.getElementById('routing-graph')

connector.subscribe(function(tpe, data) {
  if (tpe === 'opnupdate') {
    const path = data.fields.routing_graph.path
    if (typeof path !== "undefined") {
      img.src = path
    }
  }
})

</script>

<template>
  <nav class="navbar">
    <div class="container">
      <div class="navbar-title">Termites Debugger</div>
      <ul class="navbar-nav">
        <li class="navbar-item">
          <a class="nav-link">Graph</a>
        </li>
      </ul>
    </div>
  </nav>
  <div class="graph">
    <object id="routing-graph" type="image/svg+xml"></object>
  </div>
</template>
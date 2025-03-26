<template>
  <div id="app">
    <TelemetryDashboard />
  </div>
</template>

<script>
import TelemetryDashboard from "./components/TelemetryDashboard.vue";
import websocket from "./services/websocket";

export default {
  name: "App",
  components: {
    TelemetryDashboard
  },
  created() {
    // Use a relative URL to use the same host and port as the served page.
    const wsUrl = "ws://" + window.location.host + "/ws";
    const vm = this;
    websocket.connect(wsUrl, (data) => {
      // Dispatch the telemetry data to Vuex
      vm.$store.dispatch("telemetry/updateTelemetry", data);
    });
  },
  beforeDestroy() {
    websocket.close();
  }
};
</script>
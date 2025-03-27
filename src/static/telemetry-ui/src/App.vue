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

      if (data && data.type) {
        if(data.type === "telemetry"){
          vm.$store.dispatch("telemetry/updateTelemetry", data.cars);
        }
        else if(data.type === "port_list") {
          vm.$store.dispatch("telemetry/updatePorts", data.ports);
        }
        else {
        console.warn("Received unknown message from WebSocket:", data);
      }
      }
      
    });
  },
  beforeUnmount() {
    websocket.close();
  }
};
</script>
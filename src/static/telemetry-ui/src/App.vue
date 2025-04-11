<template>
  <div id="app">
    <!-- When on Overview, render the overview and menu overlay -->
    <div v-if="$route.name === 'Overview'">
      <MobileOverview
        :hiddenCars="hiddenCars"
        :carOrder="carOrder"
        @hamburger-clicked="menuOpen = true"
        @reorder-cars="onReorderCars"
      />
      <MobileMenu
        v-if="menuOpen"
        :hiddenCars="hiddenCars"
        :carOrder="carOrder"
        @close="menuOpen = false"
        @toggle-car="toggleCarVisibility"
      />
    </div>
    <!-- For all other routes, simply render the router view -->
    <router-view v-else />
  </div>
</template>

<script>
import MobileOverview from "./components/MobileOverview.vue";
import MobileMenu from "./components/MobileMenu.vue";
import websocket from "./services/websocket";

export default {
  name: "App",
  components: { MobileOverview, MobileMenu },
  data() {
    return {
      menuOpen: false,
      hiddenCars: [],
      carOrder: []
    };
  },
  methods: {
    toggleCarVisibility(carNum, visible) {
      if (!visible) {
        if (!this.hiddenCars.includes(carNum)) {
          this.hiddenCars.push(carNum);
        }
      } else {
        this.hiddenCars = this.hiddenCars.filter(num => num !== carNum);
      }
    },
    onReorderCars(newOrder) {
      this.carOrder = newOrder.map(car => car.CN);
      console.log("New car order:", this.carOrder);
    }
  },
  created() {
    const wsUrl = "ws://" + window.location.host + "/ws";
    websocket.connect(wsUrl, (data) => {
      console.log("App.vue received:", data);
      if (data && data.type) {
        if (data.type === "telemetry") {
          if (data.cars && Array.isArray(data.cars)) {
            this.$store.dispatch("telemetry/updateTelemetry", data.cars);
            if (!this.carOrder.length && data.cars.length) {
              this.carOrder = data.cars.map(car => car.CN);
            }
          } else {
            console.warn("Telemetry message missing 'cars' array", data);
          }
        } else if (data.type === "port_list") {
          if (data.ports && Array.isArray(data.ports)) {
            this.$store.dispatch("telemetry/updatePorts", data.ports);
          } else {
            console.warn("Port list message missing 'ports' array", data);
          }
        } else if (data.type === "lap_times") {
        if (data.lap_times && Array.isArray(data.lap_times)) {
          this.$store.dispatch("lapTimer/updateLapTimes", data.lap_times);
          }
        } else {
          console.warn("Unknown WebSocket message type:", data);
        }
      
      } else if (Array.isArray(data)) {
        this.$store.dispatch("telemetry/updateTelemetry", data);
      }
    });
  },
  beforeUnmount() {
    websocket.close();
  }
};
</script>

<style>
:root {
  --color-primary: #fceb09;
  --color-secondary: #424242;
  --color-accent: #82B1FF;
  --color-background: #FFFFFF;
  --color-text: #550101;
  --color-overlay: rgba(0, 0, 0, 0.5);
}
body {
  margin: 0;
  font-family: sans-serif;
  overflow: hidden;
}
</style>

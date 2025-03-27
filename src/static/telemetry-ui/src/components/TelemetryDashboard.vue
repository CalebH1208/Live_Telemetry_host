<template>
  <div>
    <h1>Telemetry Dashboard</h1>

    <!-- COM Port Selection -->
    <div>
      <label for="com-port">Select COM Port:</label>
      <select v-model="selectedPort" :disabled="availablePorts.length === 0 || availablePorts[0] === 'No active ports'">
        <option v-for="port in availablePorts" :key="port" :value="port">{{ port }}</option>
      </select>
      <button @click="connectToPort">Connect</button>
    </div>
    <hr>
    <!-- Car Telemetry Display -->
    <div v-for="car in telemetryData" :key="car.CN" class="car-section">
      <h2>Car {{ car.CN }}</h2>
      <table border="1">
        <thead>
          <tr>
            <th>Name</th>
            <th>Value</th>
            <th>Unit</th>
            <th>Filter</th>
            <th>Precision</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loop through telemetry values for this car -->
          <tr v-for="tv in (car.TV || [])" :key="tv.N">
            <td>{{ tv.N }}</td>
            <td>{{ formatValue(tv.V, getPrecision(car, tv)) }}</td>
            <td>{{ tv.U }}</td>
            <td>
              <input
                type="checkbox"
                :checked="getFilter(car, tv)"
                @change="toggleFilter(car, tv, $event)"
              />
            </td>
            <td>
              <select :value="getPrecision(car, tv)" @change="changePrecision(car, tv, $event)">
                <!-- Options 0 to 7 -->
                <option v-for="n in 8" :key="n-1" :value="n-1">{{ n-1 }}</option>
              </select>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import Cookies from "js-cookie";
import { mapGetters } from "vuex";
import websocket from "@/services/websocket";

export default {
  name: "TelemetryDashboard",
  data() {
    return {
      selectedPort: ""
    };
  },
  computed: {
    ...mapGetters("telemetry", ["telemetryData", "availablePorts"])
  },
  methods: {
    // Cookie key functions.
    filterCookieKey(car, tv) {
      return `${car.CN}_${tv.N}_F`;
    },
    precisionCookieKey(car, tv) {
      return `${car.CN}_${tv.N}_P`;
    },
    // Get filter state from cookie (default true).
    getFilter(car, tv) {
      const key = this.filterCookieKey(car, tv);
      const cookieVal = Cookies.get(key);
      return cookieVal === undefined ? true : cookieVal !== "false";
    },
    // Toggle filter and update cookie.
    toggleFilter(car, tv, event) {
      const newValue = event.target.checked;
      const key = this.filterCookieKey(car, tv);
      if (newValue) {
        Cookies.remove(key);
      } else {
        Cookies.set(key, "false");
      }
      this.$forceUpdate();
    },
    // Get precision value from cookie (default 3).
    getPrecision(car, tv) {
      const key = this.precisionCookieKey(car, tv);
      const cookieVal = Cookies.get(key);
      return cookieVal === undefined ? 3 : parseInt(cookieVal, 10);
    },
    // Change precision and update cookie accordingly.
    changePrecision(car, tv, event) {
      const newValue = parseInt(event.target.value, 10);
      const key = this.precisionCookieKey(car, tv);
      if (newValue === 3) {
        Cookies.remove(key);
      } else {
        Cookies.set(key, newValue.toString());
      }
      this.$forceUpdate();
    },
    // Format the telemetry value based on precision.
    formatValue(value, precision) {
      // Use toFixed then parseFloat to remove trailing zeros.
      const fixed = value.toFixed(precision);
      return precision > 0 ? parseFloat(fixed).toString() : fixed;
    },
    // Send selected COM port to backend.
    connectToPort() {
      if (this.selectedPort && this.selectedPort !== "No active ports") {
        const msg = { type: "select_port", port: this.selectedPort };
        // Assuming your WebSocket service is imported globally or via a plugin.
        websocket.sendMessage(JSON.stringify(msg));
      }
    }
  }
};
</script>

<style scoped>
.car-section {
  margin-bottom: 2rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 0.5rem;
}

th,
td {
  padding: 0.5rem;
  text-align: left;
}
</style>

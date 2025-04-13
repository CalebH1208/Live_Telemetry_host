<template>
  <div class="host-settings-container">
    <!-- Fixed Top Bar with Back Button -->
    <div class="top-bar">
      <button class="back-button" @click="goBack">&larr;</button>
      <h1 class="title">Host Settings</h1>
    </div>
    <!-- Main Content -->
    <div class="content">
      <div class="port-selection">
        <label for="port-dropdown">Select Port:</label>
        <select v-model="selectedPort" id="port-dropdown">
          <option v-for="port in availablePorts" :key="port" :value="port">
            {{ port }}
          </option>
        </select>
        <button class="connect-button" @click="connectPort">Connect</button>
        <button class="reset-button" @click="reset">Reset Application</button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import websocket from "@/services/websocket";

export default {
  name: "MobileHostSettings",
  data() {
    return {
      selectedPort: ""
    };
  },
  computed: {
    ...mapGetters("telemetry", ["availablePorts"])
  },
  methods: {
    connectPort() {
      if (this.selectedPort && this.selectedPort !== "No active ports") {
        const msg = { type: "select_port", port: this.selectedPort };
        websocket.sendMessage(JSON.stringify(msg));
      }
    },
    reset() {
      const msg = { type: "reset", kill: "true" };
      websocket.sendMessage(JSON.stringify(msg));
    },
    goBack() {
      // Navigate back to the Overview page.
      this.$router.push({ name: "Overview" });
    }
  }
};
</script>

<style scoped>
.host-settings-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--color-background);
  color: var(--color-text);
}
.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-primary);
  color: var(--color-text);
  padding: 20px;
  margin-bottom: 10px;
  z-index: 10;
}
.back-button {
  position: absolute;
  left: 0px;
  background: none;
  border: none;
  font-size: 8vw;
  color: var(--color-text);
  cursor: pointer;
}
.title {
  margin: 0;
  font-size: 8vw;
}
.content {
  margin-top: 80px;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.port-selection {
  display: flex;
  flex-direction: column;
  align-items: center;
}
#port-dropdown {
  margin: 10px 0;
  font-size: 5vw;
  padding: 10px;
}
.connect-button,
.reset-button {
  font-size: 5vw;
  padding: 10px 20px;
  background-color: var(--color-secondary);
  color: var(--color-text);
  border: none;
  border-radius: 5px;
  margin: 10px 0;
  margin-top: 50px;
}
</style>

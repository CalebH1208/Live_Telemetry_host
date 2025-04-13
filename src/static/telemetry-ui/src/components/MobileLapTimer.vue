<template>
    <div class="lap-timer-container">
      <!-- Fixed Top Bar -->
      <div class="top-bar">
        <button class="back-button" @click="goBack">&larr;</button>
        <h1 class="title">Lap Timer</h1>
      </div>
      <!-- Lap Times List -->
      <div class="lap-list">
        <div v-if="lapTimes && lapTimes.length > 0">
          <div v-for="(lap, index) in lapTimes" :key="index" class="lap-item">
            <span class="lap-number">Lap {{ index }}:</span>
            <span class="lap-time">{{ lap }}</span>
          </div>
        </div>
        <div v-else class="no-laps">
          No laps recorded yet.
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { mapGetters } from "vuex";
  export default {
    name: "MobileLapTimer",
    computed: {
      ...mapGetters("lapTimer", ["lapTimes"])
    },
    methods: {
      goBack() {
        this.$router.push({ name: "Overview" });
      }
    }
  };
  </script>
  
  <style scoped>
  .lap-timer-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background-color: var(--color-background);
    color: var(--color-text);
    font-size: 1.5em;
  }
  
  /* Fixed Top Bar */
  .top-bar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    background-color: var(--color-primary);
    color: var(--color-text);
    padding: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
  }
  .back-button {
    position: absolute;
    left: 20px;
    background: none;
    border: none;
    font-size: 2em;
    color: var(--color-text);
    cursor: pointer;
  }
  .title {
    margin: 0;
    font-size: 2em;
  }
  
  /* Lap List */
  .lap-list {
    margin-top: 80px;
    flex: 1;
    overflow-y: auto;
    padding: 20px;
  }
  .lap-item {
    display: flex;
    justify-content: space-between;
    padding: 15px;
    margin-bottom: 10px;
    background-color: var(--color-secondary);
    color: var(--color-text);
    border-radius: 5px;
  }
  .lap-number {
    font-weight: bold;
  }
  .no-laps {
    text-align: center;
    font-size: 1.2em;
    margin-top: 20px;
  }
  </style>
  
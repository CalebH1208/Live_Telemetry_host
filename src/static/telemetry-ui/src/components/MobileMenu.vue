<template>
    <div class="menu-overlay" @click.self="$emit('close')">
      <div class="menu-container">
        <div class="menu-buttons">
          <button class="menu-item" @click="onHostSettings">Host Settings</button>
          <button class="menu-item" @click="onLapTimer">Lap Timer</button>
        </div>
        <h3 class="menu-header">Cars:</h3>
        <div class="cars-list">
          <div
            v-for="(car, index) in sortedCars"
            :key="car.CN"
            class="car-item"
            :class="{ alt: index % 2 === 1 }"
          >
            <span>Car {{ car.CN }}</span>
            <button class="visibility-button" @click="toggleCar(car.CN)">
              {{ isHidden(car.CN) ? "Show" : "Hide" }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { mapGetters } from "vuex";
  export default {
    name: "MobileMenu",
    props: {
      carOrder: {
        type: Array,
        default: () => []
      },
      hiddenCars: {
        type: Array,
        default: () => []
      }
    },
    computed: {
      ...mapGetters("telemetry", ["telemetryData"]),
      sortedCars() {
        if (!this.carOrder || !this.carOrder.length) {
          return this.telemetryData;
        }
        return this.telemetryData.slice().sort((a, b) => {
          const aIdx = this.carOrder.indexOf(a.CN);
          const bIdx = this.carOrder.indexOf(b.CN);
          return aIdx - bIdx;
        });
      }
    },
    methods: {
      onHostSettings() {
        // Navigate to the Host Settings page using vue-router.
        this.$router.push({ name: "HostSettings" });
      },
      onLapTimer() {
        this.$router.push({ name: "LapTimer" });
      },
      // Toggle the hidden status of a car:
      toggleCar(carNum) {
        // If the car is currently hidden, emit event to show it; otherwise, hide it.
        if (this.isHidden(carNum)) {
          this.$emit("toggle-car", carNum, true);
        } else {
          this.$emit("toggle-car", carNum, false);
        }
      },
      isHidden(carNum) {
        return this.hiddenCars.includes(carNum);
      }
    }
  };
  </script>
  
  <style scoped>
  .menu-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: var(--color-overlay);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }
  .menu-container {
    background-color: var(--color-background);
    color: var(--color-text);
    width: 80%;
    max-width: 400px;
    padding: 30px;
    border-radius: 10px;
    text-align: center;
    font-size: 1.5em;
  }
  .menu-buttons {
    margin-bottom: 20px;
  }
  .menu-item {
    width: 100%;
    padding: 15px;
    margin-bottom: 15px;
    background-color: var(--color-primary);
    color: var(--color-text);
    border: none;
    font-size: 1.5em;
    border-radius: 5px;
  }
  .menu-header {
    margin-top: 20px;
    font-size: 1.8em;
  }
  .cars-list {
    text-align: left;
    margin-top: 10px;
  }
  .car-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    font-size: 1.5em;
  }
  .car-item.alt {
    background-color: var(--color-accent);
  }
  .visibility-button {
    padding: 5px 10px;
    font-size: 1.2em;
    border: none;
    border-radius: 5px;
    background-color: var(--color-primary);
    color: var(--color-background);
    cursor: pointer;
  }
  </style>
  
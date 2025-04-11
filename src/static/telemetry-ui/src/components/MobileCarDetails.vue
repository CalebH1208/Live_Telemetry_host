<template>
  <div class="details-container" v-if="currentCar">
    <!-- Fixed Top Bar -->
    <div class="top-bar">
      <button class="back-button" @click="goBack">&larr;</button>
      <h1 class="title">Car: {{ currentCar.CN }}</h1>
    </div>
    <!-- Draggable Telemetry List -->
    <div class="telemetry-list">
      <draggable 
        :list="localTelemetry" 
        item-key="N" 
        :animation="300"
        @start="onDragStart" 
        @end="onDragEnd">
        <template #item="{ element, index }">
          <div class="telemetry-item" :class="{ alt: index % 2 === 1 }">
            <!-- Row Main is always visible; clicking toggles expansion -->
            <div class="row-main" @click.stop="toggleExpanded(element)">
              <span class="name">{{ element.N }}</span>
              <span class="value">{{ formatValue(element.V, getPrecision(element)) }}</span>
              <span class="unit">{{ element.U }}</span>
            </div>
            <!-- Collapsible Controls -->
            <transition name="fade">
              <div class="row-controls" v-show="isExpanded(element)" @click.stop>
                <button class="filter-button" @click.stop="toggleFilterButton(element)">
                  Filter: {{ getFilterText(element) }}
                </button>
                <label class="precision-label" @click.stop>
                  Precision:
                  <select :value="getPrecision(element)" @change.stop="changePrecision(element, $event)">
                    <option v-for="n in 8" :key="n-1" :value="n-1">{{ n-1 }}</option>
                  </select>
                </label>
              </div>
            </transition>
          </div>
        </template>
      </draggable>
    </div>
  </div>
  <div v-else class="loading">
    <div class="top-bar">
      <button class="back-button" @click="goBack">&larr;</button>
      <h1 class="title">Waiting ...</h1>
    </div>
    Loading car details...
  </div>
</template>

<script>
import draggable from "vuedraggable";
import Cookies from "js-cookie";

export default {
  name: "MobileCarDetails",
  components: { draggable },
  data() {
    return {
      localTelemetry: [], // Local copy for draggable reordering.
      expanded: {},       // Track expanded state keyed by telemetry name.
      dragging: false     // Flag to ignore clicks while dragging.
    };
  },
  computed: {
    // Look up the current car based on the route parameter.
    currentCar() {
      const carId = parseInt(this.$route.params.id, 10);
      return this.$store.getters["telemetry/telemetryData"].find(
        (car) => car.CN === carId
      );
    }
  },
  watch: {
    // When currentCar changes, load telemetry order and reset expanded.
    currentCar: {
      handler(newCar, oldCar) {
        if (!oldCar || newCar.CN !== oldCar.CN) {
          if (newCar && newCar.TV) {
            this.loadTelemetryWithSavedOrder(newCar);
            this.expanded = {};
          }
        }
      },
      immediate: true
    },
    // Watch telemetryData (deeply) and update localTelemetry if there are new values.
    telemetryData: {
      handler(newData) {
        if (this.currentCar && this.currentCar.TV) {
          this.loadTelemetryWithSavedOrder(this.currentCar);
        }
      },
      deep: true,
      immediate: true
    }
  },
  methods: {
    goBack() {
      this.$router.push({ name: "Overview" });
    },
    telemetryOrderCookieKey() {
      return `car_${this.currentCar.CN}_telemetry_order`;
    },
    saveTelemetryOrder() {
      const order = this.localTelemetry.map((item) => item.N);
      Cookies.set(this.telemetryOrderCookieKey(), JSON.stringify(order), { expires: 365 });
    },
    loadTelemetryWithSavedOrder(car) {
      const savedOrder = Cookies.get(this.telemetryOrderCookieKey());
      if (savedOrder) {
        try {
          const orderArray = JSON.parse(savedOrder);
          this.localTelemetry = [...car.TV].sort(
            (a, b) => orderArray.indexOf(a.N) - orderArray.indexOf(b.N)
          );
        } catch (error) {
          console.error("Error parsing telemetry order cookie", error);
          this.localTelemetry = [...car.TV];
        }
      } else {
        this.localTelemetry = [...car.TV];
      }
    },
    onDragStart() {
      this.dragging = true;
    },
    onDragEnd(evt) {
      this.dragging = false;
      this.saveTelemetryOrder();
    },
    // Toggle expanded state for a telemetry item.
    toggleExpanded(tv) {
      if (this.dragging) return; // Prevent toggle during drag.
      // Toggle the expanded state keyed by the telemetry name.
      this.expanded[tv.N] = !this.expanded[tv.N];
      // Force reactivity update.
      this.expanded = { ...this.expanded };
    },
    isExpanded(tv) {
      return this.expanded[tv.N] === true;
    },
    // Filter button as toggle (instead of checkbox).
    filterCookieKey(tv) {
      return `${this.currentCar.CN}_${tv.N}_F`;
    },
    getFilterText(tv) {
      const key = this.filterCookieKey(tv);
      const cookieVal = Cookies.get(key);
      return cookieVal === undefined || cookieVal !== "false" ? "ON" : "OFF";
    },
    toggleFilterButton(tv) {
      const key = this.filterCookieKey(tv);
      const current = Cookies.get(key);
      if (current === undefined || current !== "false") {
        Cookies.set(key, "false");
      } else {
        Cookies.remove(key);
      }
      this.$forceUpdate();
    },
    // Precision methods.
    precisionCookieKey(tv) {
      return `${this.currentCar.CN}_${tv.N}_P`;
    },
    getPrecision(tv) {
      const key = this.precisionCookieKey(tv);
      const cookieVal = Cookies.get(key);
      return cookieVal === undefined ? 3 : parseInt(cookieVal, 10);
    },
    changePrecision(tv, event) {
      const newValue = parseInt(event.target.value, 10);
      const key = this.precisionCookieKey(tv);
      if (newValue === 3) {
        Cookies.remove(key);
      } else {
        Cookies.set(key, newValue.toString());
      }
      this.$forceUpdate();
    },
    formatValue(value, precision) {
      const factor = Math.pow(10, precision);
      return parseFloat((Math.floor(value * factor) / factor).toFixed(precision));
    },
    onDragUpdate(evt) {
      this.$emit("update-telemetry-order", this.localTelemetry);
    }
  }
};
</script>

<style scoped>
.details-container {
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
  color: var(--color-background);
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
  color: var(--color-background);
  cursor: pointer;
}
.title {
  margin: 0;
  font-size: 2em;
}

/* Draggable Telemetry List */
.telemetry-list {
  margin-top: 80px;
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.telemetry-item {
  background-color: var(--color-secondary);
  color: var(--color-background);
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 8px;
  transition: background-color 0.3s, max-height 0.3s, opacity 0.3s;
  overflow: hidden;
}
.telemetry-item.alt {
  background-color: var(--color-accent);
}
.row-main {
  display: flex;
  justify-content: space-around;
  margin-bottom: 10px;
  cursor: pointer;
}
.row-controls {
  display: flex;
  justify-content: space-around;
  font-size: 0.9em;
}
.row-controls label {
  display: flex;
  align-items: center;
}
.row-controls select {
  margin-left: 5px;
}
/* Fade transition for expansion */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.loading {
  display: flex;
  height: 100vh;
  justify-content: center;
  align-items: center;
  font-size: 2em;
}
</style>

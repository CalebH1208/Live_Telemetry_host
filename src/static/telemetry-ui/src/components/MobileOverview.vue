<template>
  <div class="overview-container">
    <!-- Fixed Top Bar -->
    <div class="top-bar">
      <!-- Modified hamburger button -->
      <button class="hamburger" @click="handleHamburger">&#9776;</button>
      <h1 class="title">MR Telemetry</h1>
    </div>

    <!-- Scrollable Car Placards with Draggable -->
    <div class="cars-list">
      <draggable :list="localCars" item-key="CN" @update="onDragUpdate">
        <template #item="{ element, index }">
          <!-- Tapping a placard navigates to Car Details -->
          <div class="car-placard" :class="{ alt: index % 2 === 1 }" @click="goToCarDetails(element)">
            <h2>Car {{ element.CN }}</h2>
            <div v-for="(tv, idx) in firstN(element,20)" :key="tv.N" class="telemetry-block">
              <div class="telemetry-row">
                <span class="telemetry-name">{{ tv.N }}</span>
                <span class="telemetry-value">{{ formatValue(tv.V, getPrecision(element, tv)) }}</span>
                <span class="telemetry-unit">{{ tv.U }}</span>
              </div>
            </div>
          </div>
        </template>
      </draggable>
    </div>

    <!-- Fixed Footer -->
    <div class="footer">
      <span>Last lap time:</span>
      <div class="lap-time-placeholder">{{ lastLap }}</div>
    </div>
  </div>
</template>

<script>
import draggable from "vuedraggable";
import Cookies from "js-cookie";
import { mapGetters } from "vuex";

export default {
  name: "MobileOverview",
  components: { draggable },
  props: {
    hiddenCars: {
      type: Array,
      default: () => []
    },
    carOrder: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      localCars: []
    };
  },
  computed: {
    ...mapGetters("telemetry", ["telemetryData"]),
    // Get lap times from the lapTimer module.
    ...mapGetters("lapTimer", ["lapTimes"]),
    // Calculate the last lap time (if any); fallback to "--:--"
    lastLap() {
      return this.lapTimes && this.lapTimes.length
        ? this.lapTimes[this.lapTimes.length - 1]
        : "--:--";
    }
  },
  watch: {
    telemetryData: {
      handler(newData) {
        this.updateLocalCars(newData);
      },
      immediate: true,
      deep: true
    },
    hiddenCars: {
      handler() {
        this.updateLocalCars(this.telemetryData);
      },
      immediate: true
    },
    carOrder: {
      handler() {
        this.updateLocalCars(this.telemetryData);
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    updateLocalCars(data) {
      if (!data) return;
      let filtered = data.filter(car => !this.hiddenCars.includes(car.CN));
      if (this.carOrder && this.carOrder.length) {
        filtered.sort((a, b) => {
          const aIdx = this.carOrder.indexOf(a.CN);
          const bIdx = this.carOrder.indexOf(b.CN);
          return aIdx - bIdx;
        });
      }
      this.localCars = filtered.slice();
    },
    // New firstN method: returns the first n telemetry values based on stored order.
    firstN(car, n) {
      if (!car || !car.TV) return [];
      
      const stored = Cookies.get(`car_${car.CN}_telemetry_order`);
      if (stored) {
        let order;
        try {
          order = JSON.parse(stored);
        } catch (err) {
          console.error("Error parsing telemetry order cookie for car", car.CN, err);
          return car.TV.slice(0, n);
        }
        
        if (Array.isArray(order) && order.length > 0) {
          const orderedTelemetry = order
            .map(name => car.TV.find(tv => tv.N === name))
            .filter(item => item !== undefined);
          if (orderedTelemetry.length < n) {
            const remaining = car.TV.filter(tv => !orderedTelemetry.some(item => item.N === tv.N));
            return orderedTelemetry.concat(remaining).slice(0, n);
          }
          return orderedTelemetry.slice(0, n);
        }
      }
      return car.TV.slice(0, n);
    },
    formatValue(value, precision) {
      const factor = Math.pow(10, precision);
      return parseFloat((Math.floor(value * factor) / factor).toFixed(precision));
    },
    precisionCookieKey(car, tv) {
      return `${car.CN}_${tv.N}_P`;
    },
    getPrecision(car, tv) {
      const key = this.precisionCookieKey(car, tv);
      const cookieVal = Cookies.get(key);
      return cookieVal === undefined ? 3 : parseInt(cookieVal, 10);
    },
    onDragUpdate(evt) {
      this.$emit("reorder-cars", this.localCars);
    },
    handleHamburger() {
      this.$emit("hamburger-clicked");
    },
    goToCarDetails(car) {
      this.$router.push({ name: "CarDetails", params: { id: car.CN } });
    }
  }
};
</script>

<style scoped>
.overview-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--color-background);
  color: var(--color-text);
  font-size: 1.5em;
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
.hamburger {
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
.cars-list {
  margin-top: 80px;
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.car-placard {
  background-color: var(--color-secondary);
  color: var(--color-text);
  padding: 20px;
  margin-top: 35px;
  border-radius: 8px;
  font-size: 4vw;
  transition: background-color 0.3s;
  cursor: pointer;
}
.car-placard.alt {
  background-color: var(--color-accent);
}
.telemetry-block {
  margin-top: 10px;
}
.telemetry-row {
  display: flex;
  justify-content: space-around;
  padding: 5px 0;
}
.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: var(--color-secondary);
  color: var(--color-text);
  padding: 20px;
  display: flex;
  align-items: center;
  font-size: 1.5em;
  z-index: 10;
}
.lap-time-placeholder {
  flex: 1;
  text-align: right;
  padding-right: 50px;
  margin-right: 25px;
}
</style>

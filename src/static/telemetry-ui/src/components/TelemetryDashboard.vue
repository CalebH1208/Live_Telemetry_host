<template>
  <div>
    <h1>Telemetry Dashboard</h1>
    <!-- Loop through each car -->
    <div v-for="car in cars" :key="car.CN" class="car-section">
      <h2>Car {{ car.CN }}</h2>
      <table border="1">
        <thead>
          <tr>
            <th>Name</th>
            <th>Value</th>
            <th>Unit</th>
          </tr>
        </thead>
        <tbody>
          <!-- Loop through telemetry values for this car (guard against undefined) -->
          <tr v-for="tv in (car.TV || [])" :key="tv.N">
            <td>{{ tv.N }}</td>
            <td>{{ tv.V }}</td>
            <td>{{ tv.U }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  name: "TelemetryDashboard",
  computed: {
    // Retrieve the telemetry data (an array of cars) from Vuex
    cars() {
      const data = this.$store.getters["telemetry/telemetryData"];
      console.log("Telemetry data from store:", data);
      return data;
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

th, td {
  padding: 0.5rem;
  text-align: left;
}
</style>

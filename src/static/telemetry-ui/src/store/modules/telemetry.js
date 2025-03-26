const state = {
    // Holds an array of Car objects (each with properties CN, AF, TV, etc.)
    data: []
  };
  
  const mutations = {
    // Set the telemetry data
    setTelemetryData(state, payload) {
      state.data = payload;
    }
  };
  
  const actions = {
    updateTelemetry({ commit }, data) {
      // data is expected to be an array of Car objects from the Go backend
      commit("setTelemetryData", data);
    }
  };
  
  const getters = {
    telemetryData: (state) => state.data
  };
  
  export default {
    namespaced: true,
    state,
    mutations,
    actions,
    getters
  };
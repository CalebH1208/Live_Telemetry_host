const state = {
    // Holds an array of Car objects (each with properties CN, AF, TV, etc.)
    data: [],
    ports: []
  };
  
  const mutations = {
    // Set the telemetry data
    setTelemetryData(state, payload) {
      state.data = payload;
    },
    setPorts(state, ports) {
      state.ports = ports;
    }
  };
  
  const actions = {
    updateTelemetry({ commit }, data) {
      // data is expected to be an array of Car objects from the Go backend
      commit("setTelemetryData", data);
    },
    updatePorts({ commit }, ports) {
      commit("setPorts", ports);
    }
  };
  
  const getters = {
    telemetryData: (state) => state.data,
    availablePorts: (state) => state.ports,
  };
  
  export default {
    namespaced: true,
    state,
    mutations,
    actions,
    getters
  };
// src/store/modules/lapTimer.js
const state = {
    lapTimes: []  // Array of lap time strings.
  };
  
  const mutations = {
    setLapTimes(state, lapTimes) {
      state.lapTimes = lapTimes;
    }
  };
  
  const actions = {
    updateLapTimes({ commit }, lapTimes) {
      commit("setLapTimes", lapTimes);
    }
  };
  
  const getters = {
    lapTimes: (state) => state.lapTimes
  };
  
  export default {
    namespaced: true,
    state,
    mutations,
    actions,
    getters
  };
  
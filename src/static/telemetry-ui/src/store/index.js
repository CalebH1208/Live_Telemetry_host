import { createStore } from 'vuex';
import telemetry from './modules/telemetry';
import lapTimer from './modules/lapTimer';

export default createStore({
  modules: {
    telemetry,
    lapTimer
  }
});

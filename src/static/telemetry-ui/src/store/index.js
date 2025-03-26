import { createStore } from 'vuex';
import telemetry from './modules/telemetry';

export default createStore({
  modules: {
    telemetry
  }
});

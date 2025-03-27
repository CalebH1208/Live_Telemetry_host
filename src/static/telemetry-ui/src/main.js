import { createApp } from 'vue';
import App from './App.vue';
import store from './store';
import websocket from "./services/websocket";

const app = createApp(App);
app.config.globalProperties.$socket = websocket;
app.use(store);
app.mount('#app');
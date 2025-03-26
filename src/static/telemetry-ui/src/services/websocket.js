export default {
    socket: null,
    connect(url, onMessageCallback) {
      this.socket = new WebSocket(url);
  
      this.socket.onopen = () => {
        console.log("WebSocket connection opened.");
      };
  
      this.socket.onmessage = (event) => {
        console.log("Received data:", event.data);
        try {
          // Assuming the data is a JSON array of car objects
          const data = JSON.parse(event.data);
          if (onMessageCallback && typeof onMessageCallback === "function") {
            onMessageCallback(data);
          }
        } catch (err) {
          console.error("Error parsing JSON:", err);
        }
      };
  
      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
  
      this.socket.onclose = () => {
        console.log("WebSocket connection closed.");
      };
    },
    sendMessage(message) {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(message);
      }
    },
    close() {
      if (this.socket) {
        this.socket.close();
      }
    }
  };
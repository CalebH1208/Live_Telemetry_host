const WebSocketService = {
  socket: null,
  connect(url, onMessageCallback) {
    this.socket = new WebSocket(url);
  
    this.socket.onopen = () => {
      console.log("WebSocket connection opened.");
    };
  
    this.socket.onmessage = (event) => {
      try {
        // Try to parse the received data as JSON.
        let data = JSON.parse(event.data);
        console.log("WebSocket received:", data);
        if (onMessageCallback && typeof onMessageCallback === "function") {
          onMessageCallback(data);
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
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
    } else {
      console.error("Socket not open. Cannot send message.");
    }
  },
  close() {
    if (this.socket) {
      this.socket.close();
    }
  }
};

export default WebSocketService;
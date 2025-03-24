// Create a WebSocket connection.
// Use the correct protocol (ws:// for HTTP, wss:// for HTTPS).
const socket = new WebSocket("ws://MRTelemetry.com:8080/ws");

socket.onopen = function(event) {
  console.log("WebSocket connection opened.");
};

socket.onmessage = function(event) {
  try {
    // Parse the received JSON data (an array of car objects).
    const cars = JSON.parse(event.data);
    updateDashboard(cars);
  } catch (error) {
    console.error("Error parsing telemetry data:", error);
  }
};

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
};

socket.onerror = function(error) {
  console.error("WebSocket error:", error);
};

// Function to update the HTML dashboard.
function updateDashboard(cars) {
  const telemetryDiv = document.getElementById("telemetry");
  telemetryDiv.innerHTML = ""; // Clear previous content.
  cars.forEach(car => {
    const carElem = document.createElement("div");
    carElem.className = "car";
    carElem.innerHTML = `<h2>Car ${car.CN}</h2>
                         <pre>${JSON.stringify(car, null, 2)}</pre>`;
    telemetryDiv.appendChild(carElem);
  });
}
package main

import (
	"car"
	"log"
	"net"
	"net/http"
	"serial"
	"strconv"
	"strings"
	"sync"
	"time"
	"ws"
)

var (
	carsMap  = make(map[int]*car.Car)
	carMutex sync.RWMutex
)

func main() {
	log.SetPrefix("Main: ")
	serialChan := make(chan string, 100)
	portName := "COM8"

	port, err := serial.Open_serial(portName)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	go serial.Read_serial_message(port, serialChan)
	go process_Serial_Data(serialChan)
	go Update_active_flags(carsMap, time.Second*5)

	wsManager := ws.NewManager()
	go wsManager.StartBroadcast(getCarsSlice)

	http.HandleFunc("/ws", wsManager.HandleWS)
	http.Handle("/", http.FileServer(http.Dir("./static/telemetry-ui/dist")))

	log.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
	printServerAddresses(8080)
}

/*
gets the data out of the channel from the serial parse
expects format of :

CN:(the car number here)
(data name),(unit name),(data value) * N
|

if not in this format it will break don't @ me
*/
func process_Serial_Data(in <-chan string) {
	for msg := range in {
		lines := strings.Split(msg, "\n")
		if len(lines) < 3 {
			log.Printf("Received message too short, ignoring: %s", msg)
			continue
		}

		header := strings.TrimSpace(lines[0])
		if !strings.HasPrefix(header, "CN:") {
			log.Printf("Invalid message header, expected 'CN:': %s", header)
			continue
		}
		carNumStr := strings.TrimSpace(strings.TrimPrefix(header, "CN:"))
		carNum, err := strconv.Atoi(carNumStr)
		if err != nil {
			log.Printf("Invalid car number (%s): %v", carNumStr, err)
			continue
		}

		dataLines := lines[1 : len(lines)-1]
		dataStr := strings.Join(dataLines, "\n")

		carMutex.Lock()
		c, exists := carsMap[carNum]
		if !exists {
			// Create a new car if it does not exist.
			c = &car.Car{
				Car_num:      carNum,
				Active:       true,
				Telem_values: []car.Telem_value{},
			}
			carsMap[carNum] = c
		}
		if err := c.Update_Car(dataStr); err != nil {
			log.Printf("Error updating car %d: %v", carNum, err)
		}
		carMutex.Unlock()
	}
}

func Update_active_flags(cars_map map[int]*car.Car, timout_period time.Duration) {
	for {
		for _, car := range cars_map {
			carMutex.Lock()
			car.Update_active_flag(timout_period)
			carMutex.Unlock()
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func getCarsSlice() []*car.Car {
	carMutex.RLock()
	defer carMutex.RUnlock()
	carsSlice := make([]*car.Car, 0, len(carsMap))
	for _, c := range carsMap {
		carsSlice = append(carsSlice, c)
	}
	return carsSlice
}

func printServerAddresses(port int) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting IP addresses: %v", err)
		return
	}
	log.Println("Server may be reachable at:")
	for _, addr := range addrs {
		// Check if the address is an IP address and not a loopback.
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// Only consider IPv4 addresses.
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				log.Printf("http://%s:%d/", ip4.String(), port)
			}
		}
	}
}

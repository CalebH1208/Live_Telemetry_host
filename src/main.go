package main

import (
	"car"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"serial"
	"strconv"
	"strings"
	"sync"
	"time"
	"ws"
)

var (
	carsMap    = make(map[int]*car.Car)
	carMutex   sync.RWMutex
	serialChan = make(chan string, 100)
)

func main() {
	log.SetPrefix("Main: ")

	serial.Set_up_serial_channel(serialChan)

	go process_Serial_Data(serialChan)
	go Update_active_flags(carsMap, time.Second*5)

	wsManager := ws.NewManager()
	go wsManager.StartBroadcast(getCarsSlice)
	go wsManager.Send_available_ports()

	http.HandleFunc("/ws", wsManager.HandleWS)

	static_dir := "./static/telemetry-ui/dist"
	//http.Handle("/", http.FileServer(http.Dir(static_dir)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the absolute path of the requested file
		path := filepath.Join(static_dir, r.URL.Path)
		// Check if file exists and is not a directory
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() {
			// If the file doesn't exist or it's a directory,
			// serve the index.html so the SPA can handle the route.
			http.ServeFile(w, r, static_dir)
			return
		}
		// Otherwise, serve the file.
		http.ServeFile(w, r, path)
	})

	log.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
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
		if msg == "kill" {
			carsMap = make(map[int]*car.Car)
		}
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

func Empty_car_list() {

}

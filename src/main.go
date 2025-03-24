package main

import (
	"car"
	"log"
	"serial"
	"strconv"
	"strings"
	"sync"
	"time"
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

	for {
		printCarsMap()
		time.Sleep(time.Second * 5)
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

func printCarsMap() {
	carMutex.RLock()
	defer carMutex.RUnlock()

	log.Println("----- Current Car Data -----")
	for carNum, c := range carsMap {
		serialized, err := c.Serialize()
		if err != nil {
			log.Printf("Error serializing car %d: %v", carNum, err)
			continue
		}
		log.Printf("Car %d: %s", carNum, serialized)
	}
	log.Println("----------------------------")
}

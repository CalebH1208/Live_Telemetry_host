package serial

import (
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var (
	channel chan string
	restart bool
)

func Get_port_list() ([]string, error) {
	log.SetPrefix("Serial: ")
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if len(ports) == 0 {
		return []string{"no active ports"}, nil
	}
	var ret []string
	for _, port := range ports {
		ret = append(ret, port.Name)
	}
	return ret, nil
}

func Read_serial_message(port_name string) {
	log.SetPrefix("Serial: ")
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(port_name, mode)
	if err != nil {
		log.Print(err)
		return
	}

	var buffer string
	buf := make([]byte, 1024)

	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Print(err)
			break
		}
		if restart {
			log.Print("Restarting serial ports")
			channel <- "kill"
			break
		}
		if n > 0 {
			buffer += string(buf[:n])

			for {
				var startIdx int
				startIdxcar := strings.Index(buffer, "CN:")
				startIdxlap := strings.Index(buffer, "LT:")
				if startIdxcar < startIdxlap {
					startIdx = startIdxcar
				} else {
					startIdx = startIdxlap
				}
				if startIdx == -1 {
					break
				}

				endIdx := strings.Index(buffer[startIdx:], "|")
				if endIdx == -1 {
					break
				}
				endIdx = startIdx + endIdx

				fullMessage := buffer[startIdx : endIdx+1]
				channel <- fullMessage

				buffer = buffer[endIdx+1:]
			}
		}
	}
	log.Print("closing port: ", port_name)
	port.Close()
}

func Set_up_serial_channel(mainChan chan string) {
	channel = mainChan
}

func Restart_serial() {
	restart = true
	time.Sleep(time.Second)
	restart = false
}

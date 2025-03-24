package serial

import (
	"log"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func Get_port_list() (string, error) {
	log.SetPrefix("Serial: ")
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	if len(ports) == 0 {
		return "No active Ports", nil
	}
	ret := ""
	for _, port := range ports {
		ret += port.Name
		ret += "\n"
	}
	return ret, nil
}

func Open_serial(port_name string) (serial.Port, error) {
	log.SetPrefix("Serial: ")
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open(port_name, mode)
	if err != nil {
		log.Fatal(err)
	}
	return port, err

}

func Read_serial_message(port serial.Port, out chan<- string) {
	log.SetPrefix("Serial: ")
	var buffer string
	buf := make([]byte, 1024)

	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if n > 0 {
			buffer += string(buf[:n])

			for {
				startIdx := strings.Index(buffer, "CN:")
				if startIdx == -1 {
					break
				}

				endIdx := strings.Index(buffer[startIdx:], "|")
				if endIdx == -1 {
					break
				}
				endIdx = startIdx + endIdx

				fullMessage := buffer[startIdx : endIdx+1]
				out <- fullMessage

				buffer = buffer[endIdx+1:]
			}
		}
	}
}

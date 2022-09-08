package main

import (
	"github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet"
	"github.com/tarm/serial"
	"log"
	"serial_test/messages"
	"strings"
	"time"
)

func main() {

	port := ""
	if list, err := serialdet.List(); err == nil {
		for _, p := range list {
			log.Print(p.Description(), " ", p.Path())
			if strings.Contains(p.Description(), "Raspberry") {
				port = p.Path()
			}
		}
	}
	config := &serial.Config{
		Name:        port,
		Baud:        115200,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	dev := messages.SassiDev{
		StartTime: time.Now(),
	}

	dev.Listen(stream, stream)
}

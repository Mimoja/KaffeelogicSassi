package messages

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sigurn/crc16"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type SassiDev struct {
	Max_payload_size    uint32
	Manufacturer_code   string
	Crc_table           *crc16.Table
	Serial_number       string
	Max_filename_size   uint32
	Model_number        string
	Manufacturer_domain string
	Description         string
}

func (dev *SassiDev) Parse(in string) (FullSassiMessage, error) {
	println("\n-----\nin: ", in)
	message := ParseSassiMessage(in)
	full_message := message.ToWrappedMessage().ParsePipedFields()

	if message.Message_type == CONNECTION_REQUEST {
		con_req := full_message.(ConnectionRequest)
		if dev.Manufacturer_code != "" && dev.Manufacturer_code != con_req.Manufacturer_code {
			return full_message, fmt.Errorf("Wrong device for reconnect. ManufacturerCode does not match")
		}
		if dev.Serial_number != "" && dev.Serial_number != con_req.Serial_number {
			return full_message, fmt.Errorf("Wrong device for reconnect. SerialNumber does not match")
		}
		dev.Max_payload_size = con_req.Max_packet_size
		dev.Max_filename_size = con_req.Max_packet_size
		dev.Manufacturer_code = con_req.Manufacturer_code
		dev.Serial_number = con_req.Serial_number
		dev.Model_number = con_req.Model_number
		dev.Manufacturer_domain = con_req.Manufacturer_domain
		dev.Description = con_req.Description
		param := crc16.CRC16_XMODEM
		param.Init = con_req.Crc_initial
		dev.Crc_table = crc16.MakeTable(param)
	}

	if dev.Crc_table == nil {
		panic("No connection estabilshed yet")
	}
	// Get checksum of the raw byte stream
	last := strings.LastIndexAny(in, "|")
	checksum := crc16.Checksum([]byte(in[:last+1]), dev.Crc_table)

	if checksum != full_message.CRC() {
		return full_message, fmt.Errorf("Message does not match checksum!")
	} else {
		println("Checksum valid!")
	}
	// Get checksum of the parsed and emitted byte stream
	if checksum != dev.GenerateChecksum(full_message) {
		panic("Checksum matched but we lost it during parsing!")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	//enc.Encode(full_message)
	if full_message.String() != in {
		panic("Failure during formating")
	}
	return full_message, nil
}

func (dev *SassiDev) GenerateChecksum(message FullSassiMessage) uint16 {
	mString := message.String()
	last := strings.LastIndexAny(mString, "|")
	checksum := crc16.Checksum([]byte(mString[:last+1]), dev.Crc_table)
	return checksum
}

func (dev *SassiDev) Listen(r io.Reader) {
	startTime := time.Now()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		input := scanner.Text()
		msg, err := dev.Parse(input)
		if err != nil {
			log.Fatal(err)
		}
		//TODO statemachine updates here!
		if msg.Type() == CONNECTION_REQUEST || msg.Type() == TIME_SYNC_REQUEST {
			now := time.Now().UTC()
			ts := time.Since(startTime).Milliseconds()

			tsMsg := NewTimeSync(
				dev,
				ts,
				now)
			fmt.Printf("TimeSync answer: %v\n", tsMsg.String())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

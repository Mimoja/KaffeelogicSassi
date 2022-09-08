package messages

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sigurn/crc16"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type SassiDev struct {
	Connected           bool
	Max_payload_size    uint32
	Manufacturer_code   string
	Crc_inital          uint16
	Crc_table           *crc16.Table
	Serial_number       string
	Max_filename_size   uint32
	Model_number        string
	Manufacturer_domain string
	Description         string
	StartTime           time.Time
	waitingForAck       bool
	cmdQueue            []FullSassiMessage
}

func (dev *SassiDev) Parse(in string) (FullSassiMessage, error) {
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
		dev.Max_filename_size = con_req.Max_filename_size
		dev.Manufacturer_code = con_req.Manufacturer_code
		dev.Serial_number = con_req.Serial_number
		dev.Model_number = con_req.Model_number
		dev.Manufacturer_domain = con_req.Manufacturer_domain
		dev.Description = con_req.Description
		dev.Crc_inital = con_req.Crc_initial
		param := crc16.CRC16_XMODEM
		param.Init = dev.Crc_inital
		dev.Crc_table = crc16.MakeTable(param)
		dev.Connected = true
	}

	if dev.Crc_table == nil {
		// FIXME
		//panic("No connection estabilshed yet")
		return full_message, nil
	}
	// Get checksum of the raw byte stream
	last := strings.LastIndexAny(in, "|")
	checksum := crc16.Checksum([]byte(in[:last+1]), dev.Crc_table)

	if checksum != full_message.CRC() {
		return full_message, fmt.Errorf("Message does not match checksum!")
	}
	// Get checksum of the parsed and emitted byte stream
	if checksum != dev.GenerateChecksum(full_message) {
		//	panic("Checksum matched but we lost it during parsing!")
	}

	if full_message.String() != in {
		//	panic("Failure during formating")
	}
	return full_message, nil
}

func (dev *SassiDev) GenerateChecksum(message FullSassiMessage) uint16 {
	mString := message.String()
	last := strings.LastIndexAny(mString, "|")
	checksum := crc16.Checksum([]byte(mString[:last+1]), dev.Crc_table)
	return checksum
}

func (dev *SassiDev) sendMessage(writer io.Writer, message FullSassiMessage) {
	println("Sending message: ", message.String())
	writer.Write([]byte(message.String() + "\r"))
	dev.waitingForAck = true
}

func (dev *SassiDev) sendAck(writer io.Writer, message FullSassiMessage) {
	println("Sending ack message: ", message.String())
	writer.Write([]byte(message.String() + "\r"))
	dev.waitingForAck = false
}

func (dev *SassiDev) Listen(r io.Reader, writer io.Writer) {
	startTime := time.Now()

	liveData := []byte{}
	fileData := []byte{}
	//lastFile := ""
	reader := bufio.NewReader(r)
	for true {
		input, err := reader.ReadString('\r')
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			log.Fatal(err)
		}
		input = strings.ReplaceAll(input, "\r", "")
		msg, err := dev.Parse(input)
		if err != nil {
			log.Fatal(err)
		}
		switch msg.Type() {
		case CONNECTION_REQUEST:
			fallthrough
		case TIME_SYNC_REQUEST:
			now := time.Now().UTC()
			ts := time.Since(startTime).Milliseconds()

			tsMsg := NewTimeSync(
				dev,
				ts,
				now)
			fmt.Printf("TimeSync answer: %v\n", tsMsg.String())
			dev.sendMessage(writer, tsMsg)

		case TIME_SYNC_ACK:
			println("Recieved Time Sync ACK")

		case FILE_RESPONSE:
			println("Got file response")
			fr := msg.(FileResponse)
			fileData = append(fileData, fr.Data...)
			if fr.OutcomeCode == SASSI_OUTCOME_CODE_LAST_PACKET {
				err := ioutil.WriteFile("./offlinedata.txt", fileData, 0777)
				if err != nil {
					panic(err)
				}
			}
		case NOTIFY_STATUS:
			println("Status Notify", input)
			nstat := msg.(NotifyStatus)
			if nstat.InfoCode == SASSI_INFO_CODE_APPLIANCE_BUSY {
				liveData = []byte{}
			}

			//na := NewNotifyFileAck(dev, time.Since(startTime).Milliseconds(), stat.InfoData, 0)
			//dev.sendMessage(writer, na)

			if nstat.InfoCode == SASSI_INFO_CODE_APPLIANCE_NOT_BUSY {
				err := ioutil.WriteFile("./livedata.txt", liveData, 0777)
				if err != nil {
					panic(err)
				}
				fileData = []byte{}
				//fr := NewDirListRequest(dev, time.Since(startTime).Milliseconds(), "kaffelogic/roast-logs/")
				fr := NewFileRequest(dev, time.Since(startTime).Milliseconds(), nstat.InfoData)
				dev.sendMessage(writer, fr)
				//lastFile = nstat.InfoData
			}
		case DIR_LIST_RESPONSE:
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(msg)
			dlr := msg.(DirListResponse)
			if dlr.Outcome_code == 128 {
				fr := NewFileRequest(dev, time.Since(startTime).Milliseconds(), "kaffelogic/roast-logs/log0026.klog")
				dev.sendMessage(writer, fr)
			}

		case NOTIFY_FILE:
			nf := msg.(NotifyFile)
			fmt.Printf("Notify Outcome Code is %v %v %v\n", nf.OutcomeCode, nf.SequenceNumber, nf.DateModified)
			liveData = append(liveData, nf.Data...)
		default:
			println("Unhandeled Message: ", msg.Type())
		}

	}
}

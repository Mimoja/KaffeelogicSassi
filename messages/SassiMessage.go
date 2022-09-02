package messages

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type SassiMessageType uint8

const (
	PACKET_ACK SassiMessageType = iota + 1
	CONNECTION_REQUEST
	TIME_SYNC
	TIME_SYNC_ACK
	DIR_LIST_REQUEST
	DIR_LIST_RESPONSE
	FILE_REQUEST
	FILE_RESPONSE
	FILE_PUSH
	FILE_PUSH_ACK
	FILE_DELETE
	FILE_DELETE_ACK
	INFO_REQUEST
	INFO_RESPONSE
	ACTION_SEND
	ACTION_ACK
)
const (
	TIME_SYNC_REQUEST SassiMessageType = iota + 20
)
const (
	NOTIFY_STATUS = iota + 30
	NOTIFY_STATUS_ACK
	NOTIFY_FILE
	NOTIFY_FILE_ACK
)

type SassiMessage struct {
	Manufacturer_code string
	Message_type      SassiMessageType
	Crc               uint16
	Timestamp         string
	Piped_fields      []string
}

func (message SassiMessage) CRC() uint16 {
	return message.Crc
}
func (message SassiMessage) SetCRC(crc uint16) {
	message.Crc = crc
}
func (message SassiMessage) Type() SassiMessageType {
	return message.Message_type
}

func (message SassiMessage) ParsePipedFields() FullSassiMessage {
	message.Piped_fields = []string{}
	return message
}

func (message SassiMessage) String() string {
	if strings.Contains(message.Timestamp, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(message.Manufacturer_code, "|") {
		panic("AHHHH!")
	}
	if len(message.Manufacturer_code) > 2 {
		panic("AHHHH!")
	}

	piped_fields := ""
	if len(message.Piped_fields) > 0 {
		piped_fields = strings.Join(message.Piped_fields, "|")
		piped_fields += "|"
	}
	return fmt.Sprintf("%s*%d|%s|%s%04x",
		message.Manufacturer_code,
		message.Message_type,
		message.Timestamp,
		piped_fields,
		message.Crc,
	)
}

func (message *SassiMessage) ClearPipedFields() {
	message.Piped_fields = []string{}
}

type FullSassiMessage interface {
	ParsePipedFields() FullSassiMessage
	String() string
	CRC() uint16
	Type() SassiMessageType
}

func (message SassiMessage) ToWrappedMessage() FullSassiMessage {
	switch message.Message_type {
	case PACKET_ACK:
		return PacketAck{SassiMessage: message}
	case CONNECTION_REQUEST:
		return ConnectionRequest{SassiMessage: message}
	case TIME_SYNC:
		return TimeSync{SassiMessage: message}
	case TIME_SYNC_ACK:
		return TimeSyncAck{SassiMessage: message}
	case DIR_LIST_REQUEST:
		return DirListRequest{SassiMessage: message}
	case DIR_LIST_RESPONSE:
		return DirListResponse{SassiMessage: message}
	case FILE_REQUEST:
		return FileRequest{SassiMessage: message}
	case FILE_RESPONSE:
		return FileResponse{SassiMessage: message}
	case FILE_PUSH:
		return FilePush{SassiMessage: message}
	case FILE_PUSH_ACK:
		return FilePushAck{SassiMessage: message}
	case FILE_DELETE:
		return FileDelete{SassiMessage: message}
	case FILE_DELETE_ACK:
		return FileDeleteAck{SassiMessage: message}
	case INFO_REQUEST:
		return InfoRequest{SassiMessage: message}
	case INFO_RESPONSE:
		return InfoResponse{SassiMessage: message}
	case ACTION_SEND:
		return ActionSend{SassiMessage: message}
	case ACTION_ACK:
		return ActionAck{SassiMessage: message}
	case TIME_SYNC_REQUEST:
		return TimeSyncRequest{SassiMessage: message}
	case NOTIFY_STATUS:
		return NotifyStatus{SassiMessage: message}
	case NOTIFY_STATUS_ACK:
		return NotifyStatusAck{SassiMessage: message}
	case NOTIFY_FILE:
		return NotifyFile{SassiMessage: message}
	case NOTIFY_FILE_ACK:
		return NotifyFileAck{SassiMessage: message}
	default:
		log.Panicf("Unknown type: %v", message.Message_type)
	}
	log.Panicf("Unparsed type: %v", message.Message_type)
	return nil
}

func NewSassiMessage(dev *SassiDev, ts int64, mType SassiMessageType) SassiMessage {
	timestamp := strconv.FormatInt(ts, 16)
	return SassiMessage{
		Manufacturer_code: dev.Manufacturer_code,
		Message_type:      mType,
		Crc:               0,
		Timestamp:         timestamp,
		Piped_fields:      []string{},
	}
}

func ParseSassiMessage(in string) SassiMessage {
	pipes := strings.Split(in, "|")
	if len(pipes) < 3 {
		panic("Not enough data for header")
	}
	// Header
	header := pipes[0]
	hdrs := strings.Split(header, "*")
	if len(hdrs[0]) > 2 {
		panic("Manufactoring code too long")
	}
	code := hdrs[0][:2]
	mtype, err := strconv.ParseUint(hdrs[1], 10, 8)
	if err != nil {
		panic(err)
	}

	// Timestamp
	timestamp := pipes[1]

	// CRC from last piped filed
	mLen := len(pipes)
	crc, err := strconv.ParseUint(pipes[mLen-1], 16, 32)
	if err != nil {
		panic(err)
	}

	return SassiMessage{
		Manufacturer_code: code,
		Message_type:      SassiMessageType(mtype),
		Crc:               uint16(crc),
		Timestamp:         timestamp,
		Piped_fields:      pipes[2 : mLen-1],
	}
}

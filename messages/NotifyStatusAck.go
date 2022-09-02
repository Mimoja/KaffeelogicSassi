package messages

import (
	"strconv"
)

type NotifyStatusAck struct {
	SassiMessage
	_reserved string
	InfoCode  uint8
	Serial    string
}

func NewNotifyStatusAck(dev *SassiDev, timestamp int64,
	InfoCode uint8,
) NotifyStatusAck {
	cr := NotifyStatusAck{
		SassiMessage: NewSassiMessage(dev, timestamp, NOTIFY_STATUS_ACK),
		_reserved:    "",
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r NotifyStatusAck) ParsePipedFields() FullSassiMessage {
	r._reserved = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = uint8(oCode)
	r.Serial = r.Piped_fields[0]
	return r
}

func (r NotifyStatusAck) String() string {
	r.Piped_fields = []string{
		r._reserved,
		strconv.Itoa(int(r.InfoCode)),
		r.Serial,
	}
	return r.SassiMessage.String()
}

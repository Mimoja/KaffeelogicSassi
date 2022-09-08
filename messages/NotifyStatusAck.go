package messages

import (
	"strconv"
)

type NotifyStatusAck struct {
	SassiMessage
	_reserved string
	InfoCode  InfoType
}

func NewNotifyStatusAck(dev *SassiDev, timestamp int64,
	InfoCode InfoType,
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
	r.InfoCode = InfoType(oCode)
	return r
}

func (r NotifyStatusAck) String() string {
	r.Piped_fields = []string{
		r._reserved,
		strconv.Itoa(int(r.InfoCode)),
	}
	return r.SassiMessage.String()
}

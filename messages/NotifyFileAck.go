package messages

import (
	"strconv"
)

type NotifyFileAck struct {
	SassiMessage
	FilePath    string
	OutcomeCode uint8
}

func NewNotifyFileAck(dev *SassiDev, timestamp int64,
	FilePath string,
	OutcomeCode uint8,
) NotifyFileAck {
	cr := NotifyFileAck{
		SassiMessage: NewSassiMessage(dev, timestamp, NOTIFY_FILE_ACK),
		FilePath:     FilePath,
		OutcomeCode:  OutcomeCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r NotifyFileAck) ParsePipedFields() FullSassiMessage {
	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.OutcomeCode = uint8(oCode)
	return r
}

func (r NotifyFileAck) String() string {

	r.Piped_fields = []string{
		r.FilePath,
		strconv.Itoa(int(r.OutcomeCode)),
	}
	return r.SassiMessage.String()
}

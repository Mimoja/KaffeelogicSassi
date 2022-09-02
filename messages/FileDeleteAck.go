package messages

import (
	"strconv"
)

type FileDeleteAck struct {
	SassiMessage
	FilePath    string
	OutcomeCode uint8
}

func NewFileDeleteAck(dev *SassiDev, timestamp int64,
	FilePath string,
	OutcomeCode uint8,
) FileDeleteAck {
	cr := FileDeleteAck{
		SassiMessage: NewSassiMessage(dev, timestamp, FILE_DELETE_ACK),
		FilePath:     FilePath,
		OutcomeCode:  OutcomeCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r FileDeleteAck) ParsePipedFields() FullSassiMessage {
	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.OutcomeCode = uint8(oCode)
	return r
}

func (r FileDeleteAck) String() string {

	r.Piped_fields = []string{
		r.FilePath,
		strconv.Itoa(int(r.OutcomeCode)),
	}
	return r.SassiMessage.String()
}

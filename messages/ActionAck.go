package messages

import (
	"strconv"
)

type ActionAck struct {
	SassiMessage
	InfoData string
	InfoCode uint8
	Serial   string
}

func NewActionRespone(dev *SassiDev, timestamp int64,
	InfoData string,
	InfoCode uint8,
) ActionAck {
	cr := ActionAck{
		SassiMessage: NewSassiMessage(dev, timestamp, INFO_REQUEST),
		InfoData:     InfoData,
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r ActionAck) ParsePipedFields() FullSassiMessage {
	r.InfoData = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = uint8(oCode)
	r.Serial = r.Piped_fields[0]
	return r
}

func (r ActionAck) String() string {
	r.Piped_fields = []string{
		r.InfoData,
		strconv.Itoa(int(r.InfoCode)),
		r.Serial,
	}
	return r.SassiMessage.String()
}

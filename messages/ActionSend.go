package messages

import (
	"strconv"
)

type ActionSend struct {
	SassiMessage
	InfoData string
	InfoCode uint8
	Serial   string
}

func NewActionRequest(dev *SassiDev, timestamp int64,
	InfoData string,
	InfoCode uint8,
) ActionSend {
	cr := ActionSend{
		SassiMessage: NewSassiMessage(dev, timestamp, INFO_REQUEST),
		InfoData:     InfoData,
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r ActionSend) ParsePipedFields() FullSassiMessage {
	r.InfoData = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = uint8(oCode)
	r.Serial = r.Piped_fields[0]
	return r
}

func (r ActionSend) String() string {
	r.Piped_fields = []string{
		r.InfoData,
		strconv.Itoa(int(r.InfoCode)),
		r.Serial,
	}
	return r.SassiMessage.String()
}

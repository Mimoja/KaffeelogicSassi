package messages

import (
	"strconv"
)

type InfoRequest struct {
	SassiMessage
	InfoData string
	InfoCode InfoType
}

func NewInfoRequest(dev *SassiDev, timestamp int64,
	InfoData string,
	InfoCode InfoType,
) InfoRequest {
	cr := InfoRequest{
		SassiMessage: NewSassiMessage(dev, timestamp, INFO_REQUEST),
		InfoData:     InfoData,
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r InfoRequest) ParsePipedFields() FullSassiMessage {
	r.InfoData = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = InfoType(oCode)
	return r
}

func (r InfoRequest) String() string {
	r.Piped_fields = []string{
		r.InfoData,
		strconv.Itoa(int(r.InfoCode)),
	}
	return r.SassiMessage.String()
}

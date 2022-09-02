package messages

import (
	"strconv"
)

type InfoResponse struct {
	SassiMessage
	InfoData string
	InfoCode uint8
}

func NewInfoResponse(dev *SassiDev, timestamp int64,
	InfoData string,
	InfoCode uint8,
) InfoResponse {
	cr := InfoResponse{
		SassiMessage: NewSassiMessage(dev, timestamp, INFO_RESPONSE),
		InfoData:     InfoData,
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r InfoResponse) ParsePipedFields() FullSassiMessage {
	r.InfoData = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = uint8(oCode)
	return r
}

func (r InfoResponse) String() string {
	r.Piped_fields = []string{
		r.InfoData,
		strconv.Itoa(int(r.InfoCode)),
	}
	return r.SassiMessage.String()
}

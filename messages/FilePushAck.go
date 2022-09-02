package messages

import (
	"strconv"
)

type FilePushAck struct {
	SassiMessage
	FilePath    string
	OutcomeCode uint8
}

func NewFilePushAck(dev *SassiDev, timestamp int64,
	FilePath string,
	OutcomeCode uint8,
) FilePushAck {
	cr := FilePushAck{
		SassiMessage: NewSassiMessage(dev, timestamp, FILE_PUSH),
		FilePath:     FilePath,
		OutcomeCode:  OutcomeCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r FilePushAck) ParsePipedFields() FullSassiMessage {
	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.OutcomeCode = uint8(oCode)
	return r
}

func (r FilePushAck) String() string {

	r.Piped_fields = []string{
		r.FilePath,
		strconv.Itoa(int(r.OutcomeCode)),
	}
	return r.SassiMessage.String()
}

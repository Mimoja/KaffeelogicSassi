package messages

import (
	"strconv"
)

type FileDelete struct {
	SassiMessage
	FilePath   string
	ActionCode uint8
}

func NewFileDelete(dev *SassiDev, timestamp int64,
	FilePath string,
	ActionCode uint8,
) FileDelete {
	cr := FileDelete{
		SassiMessage: NewSassiMessage(dev, timestamp, FILE_DELETE),
		FilePath:     FilePath,
		ActionCode:   ActionCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r FileDelete) ParsePipedFields() FullSassiMessage {
	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.ActionCode = uint8(oCode)
	return r
}

func (r FileDelete) String() string {

	r.Piped_fields = []string{
		r.FilePath,
		strconv.Itoa(int(r.ActionCode)),
	}
	return r.SassiMessage.String()
}

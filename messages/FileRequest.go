package messages

import (
	"strings"
)

type FileRequest struct {
	SassiMessage
	FilePath string
}

func NewFileRequest(dev *SassiDev, timestamp int64, FilePath string) FileRequest {
	new := FileRequest{
		SassiMessage: NewSassiMessage(dev, timestamp, FILE_REQUEST),
		FilePath:     FilePath,
	}
	new.Crc = dev.GenerateChecksum(new)
	return new
}

func (r FileRequest) ParsePipedFields() FullSassiMessage {
	r.FilePath = r.Piped_fields[0]
	return r
}

func (r FileRequest) String() string {
	if strings.Contains(r.FilePath, "|") {
		panic("AHHHH!")
	}

	r.Piped_fields = []string{
		r.FilePath,
	}
	return r.SassiMessage.String()
}

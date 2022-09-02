package messages

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type FileResponse struct {
	SassiMessage
	FilePath       string
	OutcomeCode    uint8
	DateModified   time.Time
	SequenceNumber uint16
	Data           []byte
}

func NewFileResponse(dev *SassiDev, timestamp int64,
	FilePath string,
	OutcomeCode uint8,
	DateModified time.Time,
	SequenceNumber uint16,
	B64Data []byte,
) FileResponse {

	cr := FileResponse{
		SassiMessage:   NewSassiMessage(dev, timestamp, FILE_RESPONSE),
		FilePath:       FilePath,
		OutcomeCode:    OutcomeCode,
		DateModified:   DateModified,
		SequenceNumber: SequenceNumber,
		Data:           B64Data,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r FileResponse) ParsePipedFields() FullSassiMessage {

	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.OutcomeCode = uint8(oCode)

	dateModifiedString := r.Piped_fields[2]

	weekRemoved := dateModifiedString[:8] + dateModifiedString[9:]
	r.DateModified, err = time.Parse("20060102150405", weekRemoved)
	if err != nil {
		panic(err)
	}

	seq, err := strconv.ParseUint(r.Piped_fields[3], 10, 16)
	if err != nil {
		panic(err)
	}
	r.SequenceNumber = uint16(seq)
	r.Data, err = base64.StdEncoding.DecodeString(r.Piped_fields[4])
	if err != nil {
		panic(err)
	}
	return r
}

func (r FileResponse) String() string {
	b64Data := base64.StdEncoding.EncodeToString(r.Data)
	dateModified := fmt.Sprintf("%s%d%s",
		r.DateModified.Format("20060102"),
		int(r.DateModified.Weekday()),
		r.DateModified.Format("150405"))

	r.Piped_fields = []string{
		r.FilePath,
		strconv.Itoa(int(r.OutcomeCode)),
		dateModified,
		strconv.Itoa(int(r.SequenceNumber)),
		b64Data,
	}
	return r.SassiMessage.String()
}

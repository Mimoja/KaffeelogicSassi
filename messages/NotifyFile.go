package messages

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"
)

type NotifyFile struct {
	SassiMessage
	FilePath       string
	OutcomeCode    OutcomeCode
	DateModified   time.Time
	SequenceNumber uint16
	Data           []byte
}

func NewNotifyFile(dev *SassiDev, timestamp int64,
	FilePath string,
	OutcomeCode OutcomeCode,
	DateModified time.Time,
	SequenceNumber uint16,
	B64Data []byte,
) NotifyFile {
	cr := NotifyFile{
		SassiMessage:   NewSassiMessage(dev, timestamp, NOTIFY_FILE),
		FilePath:       FilePath,
		OutcomeCode:    OutcomeCode,
		DateModified:   DateModified,
		SequenceNumber: SequenceNumber,
		Data:           B64Data,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r NotifyFile) ParsePipedFields() FullSassiMessage {

	r.FilePath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.OutcomeCode = OutcomeCode(oCode)

	if r.Piped_fields[2] != "" {
		dateModifiedString := r.Piped_fields[2]

		weekRemoved := dateModifiedString[:8] + dateModifiedString[9:]
		r.DateModified, err = time.Parse("20060102150405", weekRemoved)
		if err != nil {
			panic(err)
		}
	}

	seq, err := strconv.ParseUint(r.Piped_fields[3], 10, 16)
	if err != nil {
		panic(err)
	}
	r.SequenceNumber = uint16(seq)

	r.Data, err = base64.StdEncoding.DecodeString(r.Piped_fields[4])
	if err != nil {
		log.Panicf("%v:\n%s", err, r.Piped_fields[4])
	}
	return r
}

func (r NotifyFile) String() string {
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

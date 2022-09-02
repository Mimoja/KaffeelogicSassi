package messages

import (
	"encoding/base64"
	"strconv"
)

type DirListResponse struct {
	SassiMessage
	DirPath        string
	Outcome_code   uint8
	formatCode     uint8
	SequenceNumber uint16
	Data           []byte
}

func NewDirListResponse(dev *SassiDev, timestamp int64,
	DirPath string,
	OutcomeCode uint8,
	SequenceNumber uint16,
	B64Data []byte) DirListResponse {

	new := DirListResponse{
		SassiMessage:   NewSassiMessage(dev, timestamp, DIR_LIST_RESPONSE),
		DirPath:        DirPath,
		Outcome_code:   OutcomeCode,
		formatCode:     1,
		SequenceNumber: SequenceNumber,
		Data:           B64Data,
	}

	new.Crc = dev.GenerateChecksum(new)
	return new
}
func (r DirListResponse) ParsePipedFields() FullSassiMessage {

	r.DirPath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.Outcome_code = uint8(oCode)

	fCode, err := strconv.ParseUint(r.Piped_fields[2], 10, 8)
	if err != nil {
		panic(err)
	}
	r.formatCode = uint8(fCode)

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

func (r DirListResponse) String() string {
	b64Data := base64.StdEncoding.EncodeToString(r.Data)

	r.Piped_fields = []string{
		r.DirPath,
		strconv.Itoa(int(r.Outcome_code)),
		strconv.Itoa(int(r.formatCode)),
		strconv.Itoa(int(r.SequenceNumber)),
		b64Data,
	}
	return r.SassiMessage.String()
}

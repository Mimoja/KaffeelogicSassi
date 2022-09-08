package messages

import (
	"strconv"
)

type ActionSend struct {
	SassiMessage
	ActionData string
	ActionCode ActionType
	Serial     string
}

func NewActionSend(dev *SassiDev, timestamp int64,
	ActionData string,
	ActionCode ActionType,
) ActionSend {
	cr := ActionSend{
		SassiMessage: NewSassiMessage(dev, timestamp, ACTION_SEND),
		ActionData:   ActionData,
		ActionCode:   ActionCode,
		Serial:       dev.Serial_number,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r ActionSend) ParsePipedFields() FullSassiMessage {
	r.ActionData = r.Piped_fields[0]

	aCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.ActionCode = ActionType(aCode)
	r.Serial = r.Piped_fields[2]
	return r
}

func (r ActionSend) String() string {
	r.Piped_fields = []string{
		r.ActionData,
		strconv.Itoa(int(r.ActionCode)),
		r.Serial,
	}
	return r.SassiMessage.String()
}

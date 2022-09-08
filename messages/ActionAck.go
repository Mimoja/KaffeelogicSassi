package messages

import (
	"strconv"
)

type ActionAck struct {
	SassiMessage
	ActionData string
	ActionCode ActionType
	Serial     string
}

func NewActionAck(dev *SassiDev, timestamp int64,
	ActionData string,
	ActionCode ActionType,
) ActionAck {
	cr := ActionAck{
		SassiMessage: NewSassiMessage(dev, timestamp, ACTION_ACK),
		ActionData:   ActionData,
		ActionCode:   ActionCode,
		Serial:       dev.Serial_number,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r ActionAck) ParsePipedFields() FullSassiMessage {
	r.ActionData = r.Piped_fields[0]

	aCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.ActionCode = ActionType(aCode)
	r.Serial = r.Piped_fields[2]
	return r
}

func (r ActionAck) String() string {
	r.Piped_fields = []string{
		r.ActionData,
		strconv.Itoa(int(r.ActionCode)),
		r.Serial,
	}
	return r.SassiMessage.String()
}

package messages

import (
	"strconv"
)

type NotifyStatus struct {
	SassiMessage
	InfoData string
	InfoCode InfoType
}

func NewNotifyStatus(dev *SassiDev, timestamp int64,
	InfoData string,
	InfoCode InfoType,
) NotifyStatus {
	cr := NotifyStatus{
		SassiMessage: NewSassiMessage(dev, timestamp, NOTIFY_STATUS),
		InfoData:     InfoData,
		InfoCode:     InfoCode,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	return cr
}

func (r NotifyStatus) ParsePipedFields() FullSassiMessage {
	r.InfoData = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.InfoCode = InfoType(oCode)
	return r
}

func (r NotifyStatus) String() string {
	r.Piped_fields = []string{
		r.InfoData,
		strconv.Itoa(int(r.InfoCode)),
	}
	return r.SassiMessage.String()
}

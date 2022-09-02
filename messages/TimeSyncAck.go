package messages

type TimeSyncAck struct {
	SassiMessage
}

func NewTimeSyncAck(dev *SassiDev, timestamp int64) TimeSyncAck {
	tsa := TimeSyncAck{
		SassiMessage: NewSassiMessage(dev, timestamp, TIME_SYNC_ACK),
	}
	tsa.Crc = dev.GenerateChecksum(tsa)
	return tsa
}

func (r TimeSyncAck) ParsePipedFields() FullSassiMessage {
	r.Piped_fields = []string{}
	return r
}

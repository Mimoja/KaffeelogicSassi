package messages

type TimeSyncRequest struct {
	SassiMessage
}

func NewTimeSyncRequest(dev *SassiDev, timestamp int64) TimeSyncRequest {
	tsr := TimeSyncRequest{
		SassiMessage: NewSassiMessage(dev, timestamp, TIME_SYNC_REQUEST),
	}
	tsr.Crc = dev.GenerateChecksum(tsr)
	return tsr
}

func (r TimeSyncRequest) ParsePipedFields() FullSassiMessage {
	r.Piped_fields = []string{}
	return r
}

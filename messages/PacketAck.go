package messages

type PacketAck struct {
	SassiMessage
}

func NewPacketAck(dev *SassiDev, timestamp int64) PacketAck {
	pa := PacketAck{
		SassiMessage: NewSassiMessage(dev, timestamp, PACKET_ACK),
	}
	pa.Crc = dev.GenerateChecksum(pa)
	return pa
}

func (r PacketAck) ParsePipedFields() FullSassiMessage {
	r.Piped_fields = []string{}
	return r
}

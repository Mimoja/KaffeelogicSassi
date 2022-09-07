package messages

import (
	"strconv"
	"strings"
)

type ConnectionRequest struct {
	SassiMessage
	Platform_code       uint8
	Capabilities        string
	Serial_number       string
	Sassi_version       uint8
	Model_number        string
	Manufacturer_domain string
	Description         string
	Max_packet_size     uint32
	Max_filename_size   uint32
	Crc_initial         uint16
}

func NewConnectionRequest(dev *SassiDev, timestamp int64,
	platform_code uint8,
	Capabilities string,
	Sassi_version uint8,
	Crc_initial uint16) ConnectionRequest {

	cr := ConnectionRequest{
		SassiMessage:        NewSassiMessage(dev, timestamp, CONNECTION_REQUEST),
		Platform_code:       platform_code,
		Capabilities:        Capabilities,
		Serial_number:       dev.Serial_number,
		Sassi_version:       Sassi_version,
		Model_number:        dev.Model_number,
		Manufacturer_domain: dev.Manufacturer_domain,
		Description:         dev.Description,
		Max_packet_size:     dev.Max_payload_size,
		Max_filename_size:   dev.Max_filename_size,
		Crc_initial:         Crc_initial,
	}
	cr.Crc = dev.GenerateChecksum(cr)
	cr.String()
	return cr
}

func (r ConnectionRequest) ParsePipedFields() FullSassiMessage {
	pc, err := strconv.ParseUint(r.Piped_fields[0], 10, 8)
	if err != nil {
		panic(err)
	}
	r.Platform_code = uint8(pc)
	r.Capabilities = r.Piped_fields[1]
	r.Serial_number = r.Piped_fields[2]
	sVersion, err := strconv.ParseUint(r.Piped_fields[3], 10, 8)
	if err != nil {
		panic(err)
	}
	r.Sassi_version = uint8(sVersion)
	r.Model_number = r.Piped_fields[4]
	r.Manufacturer_domain = r.Piped_fields[5]
	r.Description = r.Piped_fields[6]

	max_pack, err := strconv.ParseUint(r.Piped_fields[7], 10, 32)
	if err != nil {
		panic(err)
	}
	r.Max_packet_size = uint32(max_pack)

	max_name, err := strconv.ParseUint(r.Piped_fields[8], 10, 32)
	if err != nil {
		panic(err)
	}
	r.Max_filename_size = uint32(max_name)

	// CRC inital is hex!
	crci, err := strconv.ParseUint(r.Piped_fields[9], 16, 16)
	if err != nil {
		panic(err)
	}
	r.Crc_initial = uint16(crci)

	return r
}

func (r ConnectionRequest) String() string {
	if strings.Contains(r.Capabilities, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(r.Serial_number, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(r.Model_number, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(r.Manufacturer_domain, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(r.Description, "|") {
		panic("AHHHH!")
	}

	r.Piped_fields = []string{
		strconv.Itoa(int(r.Platform_code)),
		r.Capabilities,
		r.Serial_number,
		strconv.Itoa(int(r.Sassi_version)),
		r.Model_number,
		r.Manufacturer_domain,
		r.Description,
		strconv.Itoa(int(r.Max_packet_size)),
		strconv.Itoa(int(r.Max_filename_size)),
		strconv.FormatInt(int64(r.Crc_initial), 16),
	}
	return r.SassiMessage.String()
}

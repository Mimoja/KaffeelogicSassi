package messages

import (
	"github.com/sigurn/crc16"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func dev() *SassiDev {
	param := crc16.CRC16_XMODEM
	param.Init = 0x8241
	dev := SassiDev{
		Max_payload_size:    512,
		Manufacturer_code:   "KL",
		Crc_table:           crc16.MakeTable(param),
		Serial_number:       "C35315",
		Max_filename_size:   192,
		Model_number:        "KN1007B",
		Manufacturer_domain: "kaffelogic.com",
		Description:         "",
	}
	return &dev
}

func TestPacketAck(t *testing.T) {
	reference := PacketAck{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      PACKET_ACK,
			Crc:               0x1fde,
			Timestamp:         "aabbccdd",
			Piped_fields:      []string{},
		},
	}
	new_m := NewPacketAck(dev(), 0xAABBCCDD)
	should_s := "KL*1|aabbccdd|1fde"

	msg, _ := dev().Parse(should_s)
	assert.Equal(t, reference, msg.(PacketAck))
	assert.Equal(t, reference, new_m)
	assert.Equal(t, reference.String(), should_s)
	assert.Equal(t, new_m.String(), should_s, "Failure when emitting message")
}

func TestConnectionRequest(t *testing.T) {
	dev := dev()
	reference := ConnectionRequest{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      CONNECTION_REQUEST,
			Crc:               0xfec0,
			Timestamp:         "2a3e2f",
			Piped_fields:      []string{},
		},
		Platform_code:       1,
		Capabilities:        "0",
		Serial_number:       "C35315",
		Sassi_version:       1,
		Model_number:        "KN1007B",
		Manufacturer_domain: "kaffelogic.com",
		Description:         "",
		Max_packet_size:     512,
		Max_filename_size:   192,
		Crc_initial:         0x861a,
	}
	should_s := "KL*2|2a3e2f|1|0|C35315|1|KN1007B|kaffelogic.com||512|192|861a|fec0"

	msg, _ := dev.Parse(should_s)
	cMsg := msg.(ConnectionRequest)
	assert.Equal(t, cMsg.Crc_initial, uint16(0x861a))
	assert.Equal(t, dev.Crc_inital, uint16(0x861a))
	assert.Equal(t, dev.Max_payload_size, uint32(512))
	assert.Equal(t, dev.Max_filename_size, uint32(192))

	new_m := NewConnectionRequest(dev, 0x2a3e2f, 1, "0", 1, 0x861a)

	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestTimeSync(t *testing.T) {
	reference := TimeSync{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      TIME_SYNC,
			Crc:               0x0211,
			Timestamp:         "f0b",
			Piped_fields:      []string{},
		},
		platform_code: 10,
		capabilities:  "",
		CurrentDate:   "202206105000541",
		sassi_version: 1,
	}
	then, _ := time.Parse("20060102150405", "20220610000541")

	new_m := NewTimeSync(dev(), 0xf0b, then)
	should_s := "KL*3|f0b|10||202206105000541|1|0211"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(TimeSync)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestTimeSyncAck(t *testing.T) {
	reference := TimeSyncAck{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      TIME_SYNC_ACK,
			Crc:               0x3f56,
			Timestamp:         "f0b",
			Piped_fields:      []string{},
		},
	}
	new_m := NewTimeSyncAck(dev(), 0xf0b)
	should_s := "KL*4|f0b|3f56"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(TimeSyncAck)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestDirListRequest(t *testing.T) {
	reference := DirListRequest{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      DIR_LIST_REQUEST,
			Crc:               0x4e21,
			Timestamp:         "1218",
			Piped_fields:      []string{},
		},
		DirPath:    "kaffelogic/roast-logs",
		formatCode: 1,
	}
	new_m := NewDirListRequest(dev(), 0x1218, "kaffelogic/roast-logs")
	should_s := "KL*5|1218|kaffelogic/roast-logs||1|4e21"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(DirListRequest)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestDirListResponse(t *testing.T) {
	reference := DirListResponse{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      DIR_LIST_RESPONSE,
			Crc:               0x1ee7,
			Timestamp:         "1268",
			Piped_fields:      []string{},
		},
		DirPath:        "",
		Outcome_code:   128,
		formatCode:     1,
		SequenceNumber: 2,
		Data:           []byte("roast-profiles"),
	}

	new_m := NewDirListResponse(dev(), 0x1268, "", 128, 2, []byte("roast-profiles"))
	should_s := "KL*6|1268||128|1|2|cm9hc3QtcHJvZmlsZXM=|1ee7"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(DirListResponse)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestFileRequest(t *testing.T) {
	reference := FileRequest{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      FILE_REQUEST,
			Crc:               0x1585,
			Timestamp:         "1250",
			Piped_fields:      []string{},
		},
		FilePath: "roast-profiles/00032.kpro",
	}
	new_m := NewFileRequest(dev(), 0x1250, "roast-profiles/00032.kpro")
	should_s := "KL*7|1250|roast-profiles/00032.kpro|1585"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(FileRequest)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestFileResponse(t *testing.T) {
	then, _ := time.Parse("20060102150405", "20220415101056")

	reference := FileResponse{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      FILE_RESPONSE,
			Crc:               0xf668,
			Timestamp:         "1298",
			Piped_fields:      []string{},
		},
		FilePath:       "roast-profiles/00032.kpro",
		OutcomeCode:    128,
		DateModified:   then,
		SequenceNumber: 1,
		Data:           []byte("hello"),
	}

	new_m := NewFileResponse(dev(), 0x1298,
		"roast-profiles/00032.kpro",
		128,
		then,
		1,
		[]byte("hello"))
	should_s := "KL*8|1298|roast-profiles/00032.kpro|128|202204155101056|1|aGVsbG8=|f668"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(FileResponse)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestFilePush(t *testing.T) {
	then, _ := time.Parse("20060102150405", "20220413205331")

	reference := FilePush{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      FILE_PUSH,
			Crc:               0xcf8e,
			Timestamp:         "1272",
			Piped_fields:      []string{},
		},
		FilePath:       "new file",
		ActionCode:     129,
		DateModified:   then,
		SequenceNumber: 1,
		Data:           []byte("hello"),
	}

	new_m := NewFilePush(dev(), 0x1272,
		"new file",
		129,
		then,
		1,
		[]byte("hello"))
	should_s := "KL*9|1272|new file|129|202204133205331|1|aGVsbG8=|cf8e"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(FilePush)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

func TestFileDelete(t *testing.T) {
	reference := FileDelete{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      FILE_DELETE,
			Crc:               0x8b98,
			Timestamp:         "1319",
			Piped_fields:      []string{},
		},
		FilePath:   "roast-profiles/00032.kpro",
		ActionCode: 1,
	}

	new_m := NewFileDelete(dev(), 0x1319,
		"roast-profiles/00032.kpro",
		1,
	)
	should_s := "KL*11|1319|roast-profiles/00032.kpro|1|8b98"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(FileDelete)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}
func TestFileDeleteAck(t *testing.T) {
	reference := FileDeleteAck{
		SassiMessage: SassiMessage{
			Manufacturer_code: "KL",
			Message_type:      FILE_DELETE_ACK,
			Crc:               0xdcbc,
			Timestamp:         "1287",
			Piped_fields:      []string{},
		},
		FilePath:    "roast-profiles/00032.kpro",
		OutcomeCode: 3,
	}

	new_m := NewFileDeleteAck(dev(), 0x1287,
		"roast-profiles/00032.kpro",
		3,
	)
	should_s := "KL*12|1287|roast-profiles/00032.kpro|3|dcbc"

	msg, _ := dev().Parse(should_s)
	cMsg := msg.(FileDeleteAck)
	cMsg.ClearPipedFields()
	assert.Equal(t, reference, cMsg)
	assert.Equal(t, reference, new_m)
	assert.Equal(t, should_s, new_m.String(), "Failure when emitting message")
}

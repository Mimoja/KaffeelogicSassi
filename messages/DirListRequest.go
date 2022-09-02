package messages

import (
	"strconv"
	"strings"
	"unsafe"
)

type DirListRequest struct {
	SassiMessage
	DirPath    string
	_reserved  string
	formatCode uint8
}

func NewDirListRequest(dev *SassiDev, timestamp int64, DirPath string) DirListRequest {
	new := DirListRequest{
		SassiMessage: NewSassiMessage(dev, timestamp, DIR_LIST_REQUEST),
		DirPath:      DirPath,
		_reserved:    "",
		formatCode:   1,
	}
	new.Crc = dev.GenerateChecksum(new)
	return new
}

func (r DirListRequest) ParsePipedFields() FullSassiMessage {
	r.DirPath = r.Piped_fields[0]
	r._reserved = r.Piped_fields[1]
	fCode, err := strconv.ParseUint(r.Piped_fields[2], 10, 8)
	if err != nil {
		panic(err)
	}
	r.formatCode = uint8(fCode)
	return r
}

func (r DirListRequest) String() string {
	if strings.Contains(r.DirPath, "|") {
		panic("AHHHH!")
	}
	if len(r._reserved) != 0 {
		panic("AHHHH!")
	}

	r.Piped_fields = []string{
		r.DirPath,
		r._reserved,
		strconv.Itoa(int(r.formatCode)),
	}
	return (*((*SassiMessage)(unsafe.Pointer(&r)))).String()
}

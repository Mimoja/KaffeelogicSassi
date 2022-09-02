package messages

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type TimeSync struct {
	SassiMessage
	platform_code uint8
	capabilities  string
	CurrentDate   string
	sassi_version uint8
}

func (r TimeSync) ParsePipedFields() FullSassiMessage {
	pc, err := strconv.ParseUint(r.Piped_fields[0], 10, 8)
	if err != nil {
		panic(err)
	}
	r.platform_code = uint8(pc)
	r.capabilities = r.Piped_fields[1]
	r.CurrentDate = r.Piped_fields[2]
	sVersion, err := strconv.ParseUint(r.Piped_fields[3], 10, 8)
	if err != nil {
		panic(err)
	}
	r.sassi_version = uint8(sVersion)
	return r
}

func NewTimeSync(dev *SassiDev, timestamp int64, date time.Time) TimeSync {
	currentDate := fmt.Sprintf("%s%d%s", date.Format("20060102"), int(date.Weekday()), date.Format("150405"))

	tsm := TimeSync{
		SassiMessage:  NewSassiMessage(dev, timestamp, TIME_SYNC),
		platform_code: 10,
		capabilities:  "",
		CurrentDate:   currentDate,
		sassi_version: 1,
	}
	tsm.Crc = dev.GenerateChecksum(tsm)

	return tsm
}

func (r TimeSync) String() string {
	if strings.Contains(r.capabilities, "|") {
		panic("AHHHH!")
	}
	if strings.Contains(r.CurrentDate, "|") {
		panic("AHHHH!")
	}

	r.Piped_fields = []string{
		strconv.Itoa(int(r.platform_code)),
		r.capabilities,
		r.CurrentDate,
		strconv.Itoa(int(r.sassi_version)),
	}
	return (*((*SassiMessage)(unsafe.Pointer(&r)))).String()
}

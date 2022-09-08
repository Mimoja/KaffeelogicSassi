package messages

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DirListResponse struct {
	SassiMessage
	DirPath        string
	Outcome_code   uint8
	formatCode     uint8
	SequenceNumber uint16
	Data           []DirListEntry
}

type DirListEntry struct {
	Directory  bool
	Name       string
	ChangeDate time.Time
	Size       uint32
}

func (entry DirListEntry) String() string {
	if entry.Directory {
		return fmt.Sprintf(">\t%s\t\t\r", entry.Name)
	}
	changeDate := fmt.Sprintf("%s%d%s",
		entry.ChangeDate.Format("20060102"),
		int(entry.ChangeDate.Weekday()),
		entry.ChangeDate.Format("150405"))

	return fmt.Sprintf(" \t%s\t%s\t%d\r", entry.Name, changeDate, entry.Size)
}

func parseDirListEntry(in string) DirListEntry {
	var entry DirListEntry
	in = strings.ReplaceAll(in, "\r", "")
	comp := strings.Split(in, "\t")
	entry.Name = comp[1]
	if comp[0] == ">" {
		entry.Directory = true
		return entry
	} else {
		entry.Directory = false
	}

	changeString := comp[2]
	weekRemoved := changeString[:8] + changeString[9:]
	date, err := time.Parse("20060102150405", weekRemoved)
	if err != nil {
		panic(err)
	}
	entry.ChangeDate = date

	size, err := strconv.ParseUint(comp[3], 10, 32)
	if err != nil {
		panic(err)
	}
	entry.Size = uint32(size)

	return entry
}

func NewDirListResponse(dev *SassiDev, timestamp int64,
	DirPath string,
	OutcomeCode uint8,
	SequenceNumber uint16,
	data []DirListEntry) DirListResponse {

	new := DirListResponse{
		SassiMessage:   NewSassiMessage(dev, timestamp, DIR_LIST_RESPONSE),
		DirPath:        DirPath,
		Outcome_code:   OutcomeCode,
		formatCode:     1,
		SequenceNumber: SequenceNumber,
		Data:           data,
	}

	new.Crc = dev.GenerateChecksum(new)
	return new
}
func (r DirListResponse) ParsePipedFields() FullSassiMessage {

	r.DirPath = r.Piped_fields[0]
	oCode, err := strconv.ParseUint(r.Piped_fields[1], 10, 8)
	if err != nil {
		panic(err)
	}
	r.Outcome_code = uint8(oCode)

	fCode, err := strconv.ParseUint(r.Piped_fields[2], 10, 8)
	if err != nil {
		panic(err)
	}
	r.formatCode = uint8(fCode)

	seq, err := strconv.ParseUint(r.Piped_fields[3], 10, 16)
	if err != nil {
		panic(err)
	}
	r.SequenceNumber = uint16(seq)

	rawData, err := base64.StdEncoding.DecodeString(r.Piped_fields[4])
	if err != nil {
		panic(err)
	}
	entryStrings := strings.Split(string(rawData), "\r")
	entries := []DirListEntry{}

	for _, entryString := range entryStrings {
		if len(entryString) == 0 {
			continue
		}
		entries = append(entries, parseDirListEntry(entryString))
	}
	r.Data = entries
	return r
}

func (r DirListResponse) String() string {

	listing := []string{}
	for _, entry := range r.Data {
		listing = append(listing, entry.String())
	}
	payload := strings.Join(listing, "")
	b64Data := base64.StdEncoding.EncodeToString([]byte(payload))

	r.Piped_fields = []string{
		r.DirPath,
		strconv.Itoa(int(r.Outcome_code)),
		strconv.Itoa(int(r.formatCode)),
		strconv.Itoa(int(r.SequenceNumber)),
		b64Data,
	}
	return r.SassiMessage.String()
}

package main

import (
	"encoding/json"
	"os"
	"serial_test/messages"
	"strings"
)

func main() {
	connectionRequest := "KL*2|2d3c9|1|0|C50009|1|KN1007B|kaffelogic.com|512|192|8241|f045\n" +
		"KL*4|f0b|3f56\n" + // Time sync ack
		"KL*20|123|990d\n" + // Time sync ReQ
		"KL*4|f0b|3f56\n" // Time sync ack

	r := strings.NewReader(connectionRequest)

	dev := messages.SassiDev{}
	dev.Listen(r)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	dlr := messages.NewDirListRequest(&dev, 0x1218, "kaffelogic/roast-logs")
	println(dlr.String())
	enc.Encode(dlr)

	nfr := messages.NewFileRequest(&dev, 0x1250, "roast-profiles/00032.kpro")
	println(nfr.String())
	enc.Encode(nfr)

}

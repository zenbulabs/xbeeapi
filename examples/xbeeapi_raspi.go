package main

import (
	"fmt"
	"github.com/tarm/serial"
	"github.com/zenbulabs/xbeeapi"
	"time"
)

func readCb(frame *xbeeapi.Frame, status xbeeapi.XBeeReadStatus) {
	fmt.Println("FrameLength:", frame.Length, "FrameType:", frame.FrameData.FrameType(), "Data:", frame.FrameData.Data(), "Checksum:", frame.Checksum)
	if frame.FrameData.FrameType() == xbeeapi.FrameTypeATCommand {
		atcommand, _ := xbeeapi.ParseATCommand(frame.FrameData)
		fmt.Println("FrameID:", atcommand.FrameID, "CmdType:", atcommand.Command, "Params:", atcommand.Params)
	} else if frame.FrameData.FrameType() == xbeeapi.FrameTypeATCommandResponse {
		//TODO
	}
	fmt.Println("Status:", status.StatusCode, "Error:", status.Error)
}

func main() {
	port, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyAMA0", Baud: 9600})
	if err != nil {
		fmt.Println("Error with port:", err)
		return
	}
	api := xbeeapi.NewXBeeAPI(port, readCb)
	api.Start()
	time.Sleep(100 * time.Millisecond)

	atcommands := []string{
		"NI",
		"AP",
	}

	for i, at := range atcommands {
		api.SendFrames(&xbeeapi.ATCommand{FrameID: byte(i + 1), Command: at})
	}
	time.Sleep(5000 * time.Millisecond)

	api.Finish()

	time.Sleep(5000 * time.Millisecond)
}

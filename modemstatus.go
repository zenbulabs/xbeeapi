package xbeeapi

import "fmt"

const (
	ModemHardwareReset      = 0x00
	ModemWatchdogTimerReset = 0x01
	ModemJoined             = 0x02
	ModemDisassociated      = 0x03
	ModemCoordinatorStarted = 0x06
	ModemNetworkKeyUpdated  = 0x07
	ModemNetworkWokeUp      = 0x0b
	ModemNetworkSleeping    = 0x0c
	ModemNetworkOvervoltage = 0x0d
	ModemKeyEstablished     = 0x10
	ModemConfigChangeInJoin = 0x11
	ModemStackError         = 0x80
)

type ModemStatus struct {
	Status byte
}

func ParseModemStatus(rfd *RawFrameData) (*ModemStatus, error) {
	if !rfd.IsValid() || len(rfd.Data()) != 1 {
		return nil, &FrameParseError{msg: "Expecting frame type modem status"}
	}

	return &ModemStatus{Status: rfd.Data()[0]}, nil
}

func (ms *ModemStatus) Description() string {
	switch ms.Status {
	case ModemHardwareReset:
		return "Hardware Reset"
	case ModemWatchdogTimerReset:
		return "Watchdog Timer Reset"
	case ModemJoined:
		return "Joined Network"
	case ModemDisassociated:
		return "Disassociated with Network"
	case ModemCoordinatorStarted:
		return "Coordinator Started"
	case ModemNetworkKeyUpdated:
		return "Network Security Key Updated"
	case ModemNetworkWokeUp:
		return "Network Woke Up"
	case ModemNetworkSleeping:
		return "Network Went To Sleep"
	case ModemNetworkOvervoltage:
		return "Voltage Supply Limit Exceeded"
	case ModemKeyEstablished:
		return "Key Establishment Completed"
	case ModemConfigChangeInJoin:
		return "Modem Config Changed While Join in Progress"
	case ModemStackError:
		return "Network Stack Error"
	}

	return fmt.Sprintf("Unknown Modem Status: %x", ms.Status)
}

func (ms *ModemStatus) RawFrameData() *RawFrameData {
	return NewRawFrameData(ms.Status)
}

func (ms *ModemStatus) IsValid() bool {
	return ms.Status <= ModemStackError
}

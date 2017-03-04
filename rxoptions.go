package xbeeapi

type RxOptionFlag byte

const (
	RxOptionPacketAcked        RxOptionFlag = 0x01
	RxOptionBroadcastPacket    RxOptionFlag = 0x02
	RxOptionEnableAPSEncyption RxOptionFlag = 0x20
	RxOptionUseTimeout         RxOptionFlag = 0x40
)

func setRxOptionsFlags(options byte, rxOptionFlags ...RxOptionFlag) byte {
	var newOptions byte
	for _, flag := range rxOptionFlags {
		newOptions |= byte(flag)
	}
	return newOptions
}

func isRxOptionsFlagSet(options byte, rxOptionFlag RxOptionFlag) bool {
	return options&byte(rxOptionFlag) != 0
}

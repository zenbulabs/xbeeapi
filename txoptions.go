package xbeeapi

type TxOptionFlag byte

const (
	TxOptionDisableRetries      TxOptionFlag = 0x01
	TxOptionIndirectAddressing  TxOptionFlag = 0x04
	TxOptionMulticastAddressing TxOptionFlag = 0x08
	TxOptionEnableAPSEncyption  TxOptionFlag = 0x20
	TxOptionUseTimeout          TxOptionFlag = 0x40
)

func setTxOptionsFlags(options byte, txOptionFlags ...TxOptionFlag) byte {
	var newOptions byte
	for _, flag := range txOptionFlags {
		newOptions |= byte(flag)
	}
	return newOptions
}

func isTxOptionsFlagSet(options byte, txOptionFlag TxOptionFlag) bool {
	return options&byte(txOptionFlag) != 0
}

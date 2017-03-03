package xbeeapi

import (
	"errors"
	"fmt"
	"strconv"
)

func concat(s []byte, others ...[]byte) []byte {
	for _, o := range others {
		if o != nil {
			s = append(s, o...)
		}
	}

	return s
}

func hexToBytes(hex string) ([]byte, error) {
	result := []byte{}
	for _, h := range hex {
		b, err := strconv.ParseUint(string(h), 16, 8)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(b))
	}

	return result, nil
}

func bytesToHex(byteArray []byte) string {
	hex := ""
	for _, b := range byteArray {
		hex += strconv.FormatUint(uint64(b), 16)
	}
	return hex
}

func address16(address string) ([]byte, error) {
	a, err := hexToBytes(address)
	if err != nil {
		return nil, err
	}
	if len(a) != 4 {
		return nil, errors.New(fmt.Sprintln("Expected 4 hex characters for address16 field:", address))
	}

	return a, err
}

func parseAddress16(address []byte) (string, error) {
	b := bytesToHex(address)
	if len(b) != 4 {
		return "", errors.New(fmt.Sprintln("Expecting 16-bit address:", address))
	}

	return b, nil
}

func address64(address string) ([]byte, error) {
	a, err := hexToBytes(address)
	if err != nil {
		return nil, err
	}
	if len(a) != 16 {
		return nil, errors.New(fmt.Sprintln("Expected 16 hex characters for address64 field:", address))
	}

	return a, err
}

func parseAddress64(address []byte) (string, error) {
	b := bytesToHex(address)
	if len(b) != 16 {
		return "", errors.New(fmt.Sprintln("Expecting 64-bit address:", address))
	}

	return b, nil
}

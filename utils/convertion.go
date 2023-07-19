package utils

import "encoding/binary"

func Uint16ToBytes(number uint16) []byte {
	byteArray := make([]byte, 2)
	binary.BigEndian.PutUint16(byteArray, number)
	return byteArray
}

func Uint32ToBytes(number uint32) []byte {
	byteArray := make([]byte, 4)
	binary.BigEndian.PutUint32(byteArray, number)
	return byteArray
}

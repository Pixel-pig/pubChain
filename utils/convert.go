package utils

import (
	"bytes"
	"encoding/binary"
)

/**
 * 将int64的数据转换为[]byte
 */
func Int2byte(num int64) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

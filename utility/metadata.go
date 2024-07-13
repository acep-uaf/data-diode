package utility

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"
)

func MakeTimestamp() int64 {
	return time.Now().UnixMicro()
}

func Verification(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func EncapsulatePayload(message string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(message))
	return encoded
}

func UnencapsulatePayload(message string) string {
	decoded, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Println(">> [!] Error decoding the message: ", err)
	}
	return string(decoded)
}

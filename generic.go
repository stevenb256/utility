package utility

import (
	"crypto/md5"
	"encoding/base64"
)

// md5Bytes returns md5 hash of bytes
func md5Bytes(buf []byte) string {
	h := md5.Sum(buf)
	return base64.StdEncoding.EncodeToString(h[:])
}

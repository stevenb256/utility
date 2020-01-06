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

// MinUint32 min of two uint32
func MinUint32(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

// MaxUint32 Max of two uint32
func MaxUint32(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

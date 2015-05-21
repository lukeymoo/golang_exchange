package helper

import (
	"encoding/hex"
	"crypto/md5"
)

/**
	Encodes a given string into an MD5 sum and returns it in string format
*/
func MD5String(text string) (string) {
	hashed := md5.New()
	hashed.Write([]byte(text))
	return hex.EncodeToString(hashed.Sum(nil))
}
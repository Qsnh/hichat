package main

import (
	"crypto"
	"encoding/hex"
)

func Md5(s string) string {
	md5 := crypto.MD5.New()
	md5.Write([]byte(s))
	return hex.EncodeToString(md5.Sum(nil))
}

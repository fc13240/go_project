package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func main() {

	msg := "hello world"
	key := "test"

	md5Inst := md5.New()
	md5Inst.Write([]byte(msg))
	byteMd5 := md5Inst.Sum(nil)
	strMd5 := hex.EncodeToString(byteMd5)
	fmt.Printf("md5 result = %x, %d, %s, %d\n", byteMd5, len(byteMd5), strMd5, len(strMd5))

	sha1Inst := sha1.New()

	sha1Inst.Write([]byte(msg))
	byteSha1 := sha1Inst.Sum([]byte(key))
	strSha1 := hex.EncodeToString(byteSha1)
	fmt.Printf("sha1 result = %x, %d, %s, %d\n", byteSha1, len(byteSha1), strSha1, len(strSha1))

	var strFmtSha1 = fmt.Sprintf("%x", byteSha1)

	fmt.Println(strFmtSha1)
}

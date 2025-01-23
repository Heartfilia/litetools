package utils

import "crypto/md5"

func GetMd5(bt []byte) (md5Buf []byte) {
	hashed := md5.New()
	hashed.Write(bt)
	md5Buf = hashed.Sum(nil)
	return
}

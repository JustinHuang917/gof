package gofcore

import (
	"crypto/md5"
	"crypto/rand"
	"io"
	"math/big"
	"strconv"
	"time"
)

func getMd5Hex(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return string(h.Sum(nil))
}

func genUId() (string, error) {
	guid := ""
	i, err := rand.Int(rand.Reader, big.NewInt(10))
	if err == nil {
		x := *i
		s := strconv.Itoa(int(x.Int64()))
		guid = getMd5Hex(time.Now().UTC().Format(time.ANSIC) + s)
	}
	return guid, err
}

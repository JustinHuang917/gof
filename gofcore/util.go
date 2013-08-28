// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func getMd5Hex(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	r := fmt.Sprintf("%x", h.Sum(nil))
	return r
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

func firstCharToUpper(s string) string {
	index := 0
	s1 := strings.Map(func(c rune) rune {
		index++
		if index == 1 {
			return unicode.ToUpper(c)
		}
		return unicode.ToLower(c)
	}, s)
	return s1
}

func indexOfString(originalStr, str string, startIndex int) int {
	l := len(originalStr)
	if l == 0 {
		return -1
	}
	index := 0
	if startIndex == 0 {
		index = strings.Index(originalStr, str)
	} else {
		tempStr := originalStr[startIndex:]
		tempIndex := strings.Index(tempStr, str)
		if tempIndex == -1 {
			return -1
		} else {
			index = tempIndex + startIndex
		}
	}
	return index
}

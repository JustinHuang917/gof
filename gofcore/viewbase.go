// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"errors"
	"fmt"
)

type ViewBase struct {
}

func (d *ViewBase) Writeout(out *bytes.Buffer, content interface{}) {
	s := fmt.Sprint(content)
	out.WriteString(s)
}

func (v *ViewBase) ErrorHandle(msg string) error {
	fmt.Println(msg)
	return errors.New(msg)
}

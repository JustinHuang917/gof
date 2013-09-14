// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

import (
	"bytes"
	"io"
)

type IActionResult interface {
	Invoke(context *HttpContext)
}

type ViewResult struct {
	Content *bytes.Buffer
}

type FileContent struct {
	ContentPath string
}

type JsonResult struct {
	Content []byte
}

func setHeader(context *HttpContext, key, value string) {
	context.ResponseWriter.Header().Set(key, value)
}

func (v *ViewResult) Invoke(context *HttpContext) {
	w := context.ResponseWriter.(io.Writer)
	w.Write(v.Content.Bytes())
}

func (j *JsonResult) Invoke(context *HttpContext) {
	setHeader(context, "Content-Type", "application/json;charset=UTF-8")
	w := context.ResponseWriter.(io.Writer)
	w.Write(j.Content)
}

func (f *FileContent) Invoke(context *HttpContext) {
	
}

// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

type ParseResult struct {
	Imports       string
	LayoutPath    string
	ModelTypeName string
	OutPutContent string
}

func (p *ParseResult) IsNeedLayout() bool {
	return p.LayoutPath != ""
}

var (
	Writeout_begin = "d.Writeout(out,"
	Writeout_end   = ")\n"
	Line           = 0
)

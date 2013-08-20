// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/justinhuang917/gof/goftool/parser"
	"strings"
)

const (
	open_tag           = "<%"
	close_tag          = "%>"
	out_tag            = "="
	end_server_tag     = "}"
	import_tag         = "import"
	model_declare_tag  = "model"
	layout_declare_tag = "layout"
	else_tag           = "}else{"
	else_if_tag        = "}elseif{"
	helper_tag         = "helper"
)

func toChars(s string) []rune {
	var chars = make([]rune, 1, 10)
	for _, c := range s {
		chars = append(chars, c)
	}
	return chars
}

type HtmlParser struct {
	//ParserModel *ParserModel
	src       []rune
	srcString string
	lineCount int
	result    *parser.ParseResult
}

func NewHtmlParser(content string) *HtmlParser {
	p := &HtmlParser{}
	p.src = toChars(content)
	p.srcString = content
	p.result = &parser.ParseResult{}
	return p
}

func (h *HtmlParser) Parse() *parser.ParseResult {
	var result = ""
	src := h.srcString //string(c.src)
	for _, span := range strings.Split(src, open_tag) {
		innerSpans := strings.Split(span, close_tag)
		l := len(innerSpans)
		if l == 1 {
			result += h.html(innerSpans[0])
		} else {
			code0 := innerSpans[0]
			code1 := innerSpans[1]
			result += h.logic(code0)
			if len(code1) > 0 {
				result += h.html(code1)
			}
		}
	}
	h.result.OutPutContent = result
	return h.result
}

func (h *HtmlParser) html(code string) string {
	h.lineCount += len(strings.Split("\n", code))
	code = strings.Replace(code, "\r", "\\r", -1)
	code = strings.Replace(code, "\n", "\\n", -1)
	code = strings.Replace(code, `"`, `\"`, -1)
	code = parser.Writeout_begin + "\"" + code + "\"" + parser.Writeout_end
	code = code + "\n"
	return code
}

func (h *HtmlParser) logic(code string) string {
	h.lineCount += len(strings.Split("\n", code))
	keyword := h.getKeyword(code)
	var parseFunc = h.getParseFunc(keyword)
	code = parseFunc(code)
	code = code + "\n"
	return code
}

func (h *HtmlParser) getKeyword(code string) (keyword string) {
	ks := strings.Split(code, " ")
	if strings.Index(code, out_tag) == 0 {
		keyword = out_tag
	} else if strings.Index(code, end_server_tag) == 0 {
		k := readSkipBlank(code, 2)
		k = strings.Replace(k, " ", "", -1)
		switch {
		case k == end_server_tag:
			keyword = end_server_tag
		case k == else_tag:
			keyword = else_tag
		case k == else_if_tag:
			keyword = else_if_tag
		}
	} else {
		keyword = ks[0]
	}
	return
}

func (h *HtmlParser) getParseFunc(keyword string) func(code string) string {
	var parserFunc func(code string) string
	switch {
	case keyword == import_tag:
		parserFunc = func(code string) string {
			h.result.Imports = strings.Join(strings.Split(code, " ")[1:], " ")
			h.result.Imports = strings.Replace(h.result.Imports, ";", "\n", -1)
			return ""
		}
	case keyword == model_declare_tag:
		parserFunc = func(code string) string {
			h.result.ModelTypeName = strings.Join(strings.Split(code, " ")[1:], " ")
			return ""
		}
	case keyword == layout_declare_tag:
		parserFunc = func(code string) string {
			h.result.LayoutPath = strings.Join(strings.Split(code, " ")[1:], " ")
			h.result.LayoutPath = strings.Replace(h.result.LayoutPath, " ", "", -1)
			return ""
		}
	case keyword == else_tag:
		parserFunc = func(code string) string {
			return code
		}
	case keyword == else_if_tag:
		parserFunc = func(code string) string {
			return code
		}
	case keyword == end_server_tag:
		parserFunc = func(code string) string {
			return "}\n"
		}
	case keyword == out_tag:
		parserFunc = func(code string) string {
			code = strings.Replace(code, out_tag, "", -1)
			code = parser.Writeout_begin + code + parser.Writeout_end
			return code
		}
	case keyword == helper_tag:
		parserFunc = func(code string) string {
			code = strings.Replace(code, helper_tag, "", -1)
			code = strings.Replace(code, "func", "", -1)
			codes := strings.Split(code, "(")
			endCode := strings.Join(codes[1:], "")
			return codes[0] + ":=func(" + endCode
		}
	default:
		parserFunc = func(code string) string {
			return code
		}
	}
	return parserFunc
}

func readSkipBlank(text string, blankCount int) string {
	count := 0
	lastchar := ' '
	var chars = make([]rune, 1, 10)
	for _, c := range text {
		if lastchar != ' ' && c == ' ' {
			count++
		}
		chars = append(chars, c)
		if count > blankCount {
			break
		}
		lastchar = c
	}
	return string(chars)
}

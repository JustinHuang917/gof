// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/JustinHuang917/gof/goftool/parser"
	"strings"
)

const (
	openTag          = "<%"
	closeTag         = "%>"
	outTag           = "="
	endServerTag     = "}"
	importTag        = "import"
	modelDeclareTag  = "model"
	layoutDeclareTag = "layout"
	elseTag          = "}else{"
	elseIfTag        = "}elseif{"
	helperTag        = "helper"
)

func toChars(s string) []rune {
	var chars = make([]rune, 1, 10)
	for _, c := range s {
		chars = append(chars, c)
	}
	return chars
}

type HtmlParser struct {
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
	src := h.srcString
	for _, span := range strings.Split(src, openTag) {
		innerSpans := strings.Split(span, closeTag)
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
	if strings.Index(code, outTag) == 0 {
		keyword = outTag
	} else if strings.Index(code, endServerTag) == 0 {
		k := readSkipBlank(code, 2)
		k = strings.Replace(k, " ", "", -1)
		switch {
		case k == endServerTag:
			keyword = endServerTag
		case k == elseTag:
			keyword = elseTag
		case k == elseIfTag:
			keyword = elseIfTag
		}
	} else {
		keyword = ks[0]
	}
	return
}

func (h *HtmlParser) getParseFunc(keyword string) func(code string) string {
	var parserFunc func(code string) string
	switch {
	case keyword == importTag:
		parserFunc = func(code string) string {
			h.result.Imports = strings.Join(strings.Split(code, " ")[1:], " ")
			h.result.Imports = strings.Replace(h.result.Imports, ";", "\n", -1)
			return ""
		}
	case keyword == modelDeclareTag:
		parserFunc = func(code string) string {
			h.result.ModelTypeName = strings.Join(strings.Split(code, " ")[1:], " ")
			return ""
		}
	case keyword == layoutDeclareTag:
		parserFunc = func(code string) string {
			h.result.LayoutPath = strings.Join(strings.Split(code, " ")[1:], " ")
			h.result.LayoutPath = strings.Replace(h.result.LayoutPath, " ", "", -1)
			return ""
		}
	case keyword == elseTag:
		parserFunc = func(code string) string {
			return code
		}
	case keyword == elseIfTag:
		parserFunc = func(code string) string {
			return code
		}
	case keyword == endServerTag:
		parserFunc = func(code string) string {
			return "}\n"
		}
	case keyword == outTag:
		parserFunc = func(code string) string {
			code = strings.Replace(code, outTag, "", -1)
			code = parser.Writeout_begin + code + parser.Writeout_end
			return code
		}
	case keyword == helperTag:
		parserFunc = func(code string) string {
			code = strings.Replace(code, helperTag, "", -1)
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

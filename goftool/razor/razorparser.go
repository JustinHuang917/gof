// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package razor

import (
	"regexp"
	"strings"
)

const (
	Symbol = '@'
	String = iota
	Variable
	GoCodeBlock
)
const (
	VariableFirstCharRegex = `^[\\(_a-zA-Z]`
	ElseConditionRegex     = `^[\s\r\n\t]*else\s*{|[\s\t]+if\?`
	VariableRegex          = `^(?:(?:\()(?:new\s+)?[a-zA-Z0-9]+(?:\.|[-\+*\/^=<>?:]|[\[\(][^\]\)]*[\]\)]|[a-zA-Z0-9]+)*(?:\))|(?:new\s+)?[a-zA-Z0-9]+(?:\.|[\[\(][^\]\)]*[\]\)]|[a-zA-Z0-9]+)*)`
	ImportTag              = "import"
	ModelTag               = "model"
	LayoutTag              = "layout"
	HelperTag              = "helper"
	CommentTag             = '*'
)

func isVariableFirstChar(char rune) bool {
	s := string(char)
	flag, _ := regexp.MatchString(VariableFirstCharRegex, s)
	return flag
}

func isVariable(char rune) bool {
	s := string(char)
	flag, _ := regexp.MatchString(VariableRegex, s)
	return flag
}

func isElseCondition(str string) bool {
	flag, _ := regexp.MatchString(ElseConditionRegex, str)
	return flag
}

type Segement struct {
	SegementType int
	Content      string
}

type ParserModel struct {
	Segements                  []*Segement
	ImportParts                string
	LayoutPath                 string
	ModelName                  string
	segmentIndex               int
	conditionOpeningBraceCount int
}

func toChars(s string) []rune {
	var chars = make([]rune, 1, 10)
	for _, c := range s {
		chars = append(chars, c)
	}
	return chars
}

type RazorParserEngine struct {
	ParserModel *ParserModel
	src         []rune
	srcString   string
}

func NewRazorParserEngine(content string) *RazorParserEngine {
	p := &RazorParserEngine{}
	p.src = toChars(content)
	p.srcString = content
	p.ParserModel = &ParserModel{}
	p.ParserModel.conditionOpeningBraceCount = 0
	p.ParserModel.segmentIndex = 0
	p.ParserModel.Segements = make([]*Segement, 0)
	return p
}

func (p *RazorParserEngine) Parse() {
	i, l := 0, len(p.src)
	var char rune
	for ; i < l; i++ {
		char = p.src[i]
		if char == Symbol {
			p.handleString(p.ParserModel.segmentIndex, i-p.ParserModel.segmentIndex)
			nextChar := p.src[i+1]
			switch {
			case nextChar == Symbol:
				p.handleEscape(p.ParserModel.segmentIndex, Symbol)
				i = p.ParserModel.segmentIndex - 1
				break
			case nextChar == '}':
				p.handleEscape(p.ParserModel.segmentIndex, nextChar)
				i = p.ParserModel.segmentIndex - 1
				break
			case nextChar == '{':
				p.handleGoCodeBlock(i + 1)
				i = p.ParserModel.segmentIndex
				break
			case isVariableFirstChar(nextChar):
				codes := strings.Split(p.srcString[i:], " ")
				if codes[0] == "if" || codes[0] == "for" {
					p.handleConditionAndLoop(i)
				} else if codes[0] == HelperTag {
					p.handleHelperFunction(i)
				} else if codes[0] == ModelTag {
					p.handleModel(i)
					i = p.ParserModel.segmentIndex - 1
				} else if codes[0] == LayoutTag {
					p.handleLayout(i)
					i = p.ParserModel.segmentIndex - 1
				} else if codes[0] == ImportTag {
					p.handleImports(i)
					i = p.ParserModel.segmentIndex - 1
				} else {
					p.handleVariable(i)
				}
				i = p.ParserModel.segmentIndex - 1
				break
			case nextChar == CommentTag:
				next_nextChar := p.src[i+2]
				if next_nextChar == CommentTag {
					p.handleCommentLine(i)
					i = p.ParserModel.segmentIndex - 1
				} else {
					p.handleCommentSegment(i)
					i = p.ParserModel.segmentIndex - 1
				}
			}
		} else if char == '}' && p.ParserModel.conditionOpeningBraceCount > 0 {
			//	fmt.Println("i:", i)
			p.handleCloseBrace(i)
			stringRemain := p.srcString[p.ParserModel.segmentIndex:] //p.subStr(p.ParserModel.segmentIndex, l)
			if isElseCondition(stringRemain) {
				p.handleConditionAndLoop(i)
			}
			i = p.ParserModel.segmentIndex - 1
		}
	}
	if p.ParserModel.segmentIndex < l {
		p.handleString(p.ParserModel.segmentIndex, l-p.ParserModel.segmentIndex)
	}

}

func (p *RazorParserEngine) subStrUtill(start int, endChar rune) string {
	src := p.src[start:]
	str := make([]rune, 0, 0)
	for _, c := range src {
		if c != endChar {
			str = append(str, c)
		} else {
			break
		}
	}
	return string(str)
}

func (p *RazorParserEngine) handleString(startIndex, len int) {
	if len <= 0 {
		return
	}
	s := p.srcString[startIndex : startIndex+len-1]
	s = strings.Replace(s, `"`, `\"`, -1)
	seg := &Segement{String, s}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
	p.ParserModel.segmentIndex = startIndex + len
}

func (p *RazorParserEngine) handleEscape(startindex int, char rune) {
	p.ParserModel.segmentIndex = startindex + 2
	seg := &Segement{String, string(char)}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
}

func (p *RazorParserEngine) handleGoCodeBlock(startIndex int) {
	variableLength := p.getGoCodeBlockLength(startIndex)
	if variableLength == -1 {
		panic("No '}' matched ")
	}
	p.ParserModel.segmentIndex = startIndex + variableLength
	seg := &Segement{
		GoCodeBlock, p.srcString[startIndex : startIndex+variableLength-1],
	}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
	p.ParserModel.conditionOpeningBraceCount++
}

func (p *RazorParserEngine) getGoCodeBlockLength(startIndex int) int {
	openingBraceCount := 0
	l := len(p.src)
	for i := startIndex; i < l; i++ {
		currentChar := p.src[i]
		if currentChar == '{' {
			openingBraceCount++
		} else if currentChar == '}' {
			c := openingBraceCount - 1
			if c == 0 {
				return i - startIndex
			} else {
				openingBraceCount--
			}
		}
	}
	return -1
}

func (p *RazorParserEngine) handleConditionAndLoop(index int) {
	stringRemain := p.srcString[index:]
	openningBraceIndex := strings.Index(stringRemain, "{")
	s := stringRemain[0 : openningBraceIndex+1]
	seg := &Segement{GoCodeBlock, s}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
	p.ParserModel.segmentIndex = index + openningBraceIndex + 1
	p.ParserModel.conditionOpeningBraceCount++
}

func (p *RazorParserEngine) handleHelperFunction(index int) {
	stringRemain := p.srcString[index:]
	openningBraceIndex := strings.Index(stringRemain, "{")
	code := stringRemain[0 : openningBraceIndex+1]
	code = strings.Replace(code, "helper", "", -1)
	code = strings.Replace(code, "func", "", -1)
	codes := strings.Split(code, "(")
	endCode := strings.Join(codes[1:], "")
	code = codes[0] + ":=func(" + endCode
	seg := &Segement{GoCodeBlock, code}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
	p.ParserModel.segmentIndex = index + openningBraceIndex + 1
	p.ParserModel.conditionOpeningBraceCount++
}

func (p *RazorParserEngine) handleCloseBrace(index int) {
	l := index - p.ParserModel.segmentIndex
	p.handleString(p.ParserModel.segmentIndex, l)
	seg := &Segement{GoCodeBlock, string(p.src[index])}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
	p.ParserModel.segmentIndex = index + 1
	p.ParserModel.conditionOpeningBraceCount--
}

func (p *RazorParserEngine) handleVariable(index int) {
	stringRemain := p.srcString[index:]
	r, _ := regexp.Compile(VariableRegex)
	variableString := r.FindString(stringRemain)
	p.ParserModel.segmentIndex = index + len(variableString)
	seg := &Segement{Variable, variableString}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
}

func (p *RazorParserEngine) handleCommentLine(index int) {
	stringRemain := p.srcString[index:]
	endIndex := strings.Index(stringRemain, "\n")
	p.ParserModel.segmentIndex = index + endIndex + 1
	content := "<!--" + stringRemain[2:endIndex] + "-->"
	seg := &Segement{String, content}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
}

func (p *RazorParserEngine) handleCommentSegment(index int) {
	endTag := string(CommentTag) + string(Symbol)
	stringRemain := p.srcString[index:]
	endIndex := strings.Index(stringRemain, endTag)
	p.ParserModel.segmentIndex = index + endIndex + 2
	content := "<!--" + stringRemain[1:endIndex] + "-->"
	seg := &Segement{String, content}
	p.ParserModel.Segements = append(p.ParserModel.Segements, seg)
}

func (p *RazorParserEngine) getLine(startIndex int) (content string, endIndex int) {
	stringRemain := p.srcString[startIndex:]
	endIndex = strings.Index(stringRemain, "\n")
	content = stringRemain[0:endIndex]
	return
}

func (p *RazorParserEngine) handleImports(index int) {
	startIndex := index + len(ImportTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.ParserModel.ImportParts = strings.Replace(lineContent, ";", "\n", -1)
	p.ParserModel.segmentIndex = startIndex + endIndex
}

func (p *RazorParserEngine) handleModel(index int) {
	startIndex := index + len(ModelTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.ParserModel.ModelName = lineContent
	p.ParserModel.segmentIndex = startIndex + endIndex
}

func (p *RazorParserEngine) handleLayout(index int) {
	startIndex := index + len(LayoutTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.ParserModel.LayoutPath = strings.Trim(lineContent, " ")
	p.ParserModel.segmentIndex = startIndex + endIndex
}

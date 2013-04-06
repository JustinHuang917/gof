// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package razor

import (
	//	"fmt"
	"github.com/justinhuang917/gof/goftool/parser"
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
	variableFirstCharRegex = `^[\\(_a-zA-Z]`
	elseConditionRegex     = `^[\s\r\n\t]*else\s*{|[\s\t]+if\?`
	variableRegex          = `^(?:(?:\()(?:new\s+)?[a-zA-Z0-9]+(?:\.|[-\+*\/^=<>?:]|[\[\(][^\]\)]*[\]\)]|[a-zA-Z0-9]+)*(?:\))|(?:new\s+)?[a-zA-Z0-9]+(?:\.|[\[\(][^\]\)]*[\]\)]|[a-zA-Z0-9]+)*)`
	importTag              = "import"
	modelTag               = "model"
	layoutTag              = "layout"
	helperTag              = "helper"
	commentTag             = '*'
)

func isVariableFirstChar(char rune) bool {
	s := string(char)
	flag, _ := regexp.MatchString(variableFirstCharRegex, s)
	return flag
}

func isVariable(char rune) bool {
	s := string(char)
	flag, _ := regexp.MatchString(variableRegex, s)
	return flag
}

func isElseCondition(str string) bool {
	flag, _ := regexp.MatchString(elseConditionRegex, str)
	return flag
}

type segement struct {
	segementType int
	Content      string
}

type parserModel struct {
	segements                  []*segement
	importParts                string
	layoutPath                 string
	modelName                  string
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
	parserModel *parserModel
	src         []rune
	srcString   string
}

func NewRazorParserEngine(content string) *RazorParserEngine {
	p := &RazorParserEngine{}
	p.src = toChars(content)
	p.srcString = content
	p.parserModel = &parserModel{}
	p.parserModel.conditionOpeningBraceCount = 0
	p.parserModel.segmentIndex = 0
	p.parserModel.segements = make([]*segement, 0)
	return p
}

func genRazorOutput(segements []*segement) string {
	output := ""
	temp := ""
	for _, seg := range segements {
		temp = ""
		switch seg.segementType {
		case String:
			temp = strings.Replace(seg.Content, "\r", "\\r", -1)
			temp = strings.Replace(temp, "\n", "\\n", -1)
			temp = parser.Writeout_begin + "\"" + temp + "\"" + parser.Writeout_end
			temp += "\n"
		case GoCodeBlock:
			temp = seg.Content
			temp += "\n"
		case Variable:
			temp = parser.Writeout_begin + seg.Content + parser.Writeout_end
			temp += "\n"
		}
		output += temp
	}
	return output
}

func parserModelToParseResult(m *parserModel) *parser.ParseResult {
	result := &parser.ParseResult{}
	result.Imports = m.importParts
	result.LayoutPath = m.layoutPath
	result.ModelTypeName = m.modelName
	result.OutPutContent = genRazorOutput(m.segements)
	return result
}

func (p *RazorParserEngine) Parse() *parser.ParseResult {
	i, l := 0, len(p.src)
	var char rune
	for ; i < l; i++ {
		char = p.src[i]
		if char == Symbol {
			p.handleString(p.parserModel.segmentIndex, i-p.parserModel.segmentIndex)
			nextChar := p.src[i+1]
			switch {
			case nextChar == Symbol:
				p.handleEscape(p.parserModel.segmentIndex, Symbol)
				i = p.parserModel.segmentIndex - 1
				break
			case nextChar == '}':
				p.handleEscape(p.parserModel.segmentIndex, nextChar)
				i = p.parserModel.segmentIndex - 1
				break
			case nextChar == '{':
				p.handleGoCodeBlock(i + 1)
				i = p.parserModel.segmentIndex
				break
			case isVariableFirstChar(nextChar):
				codes := strings.Split(p.srcString[i:], " ")
				if codes[0] == "if" || codes[0] == "for" {
					p.handleConditionAndLoop(i)
				} else if codes[0] == helperTag {
					p.handleHelperFunction(i)
				} else if codes[0] == modelTag {
					p.handleModel(i)
					i = p.parserModel.segmentIndex - 1
				} else if codes[0] == layoutTag {
					p.handleLayout(i)
					i = p.parserModel.segmentIndex - 1
				} else if codes[0] == importTag {
					p.handleImports(i)
					i = p.parserModel.segmentIndex - 1
				} else {
					p.handleVariable(i)
				}
				i = p.parserModel.segmentIndex - 1
				break
			case nextChar == commentTag:
				next_nextChar := p.src[i+2]
				if next_nextChar == commentTag {
					p.handleCommentLine(i)
					i = p.parserModel.segmentIndex - 1
				} else {
					p.handleCommentSegment(i)
					i = p.parserModel.segmentIndex - 1
				}
			}
		} else if char == '}' && p.parserModel.conditionOpeningBraceCount > 0 {
			//	fmt.Println("i:", i)
			p.handleCloseBrace(i)
			stringRemain := p.srcString[p.parserModel.segmentIndex:] //p.subStr(p.parserModel.segmentIndex, l)
			if isElseCondition(stringRemain) {
				p.handleConditionAndLoop(i)
			}
			i = p.parserModel.segmentIndex - 1
		}
	}
	if p.parserModel.segmentIndex < l {
		p.handleString(p.parserModel.segmentIndex, l-p.parserModel.segmentIndex)
	}
	parseResult := parserModelToParseResult(p.parserModel)
	return parseResult
	//return genRazorOutput(p.parserModel.segements)
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
	seg := &segement{String, s}
	p.parserModel.segements = append(p.parserModel.segements, seg)
	p.parserModel.segmentIndex = startIndex + len
}

func (p *RazorParserEngine) handleEscape(startindex int, char rune) {
	p.parserModel.segmentIndex = startindex + 2
	seg := &segement{String, string(char)}
	p.parserModel.segements = append(p.parserModel.segements, seg)
}

func (p *RazorParserEngine) handleGoCodeBlock(startIndex int) {
	variableLength := p.getGoCodeBlockLength(startIndex)
	if variableLength == -1 {
		panic("No '}' matched ")
	}
	p.parserModel.segmentIndex = startIndex + variableLength
	seg := &segement{
		GoCodeBlock, p.srcString[startIndex : startIndex+variableLength-1],
	}
	p.parserModel.segements = append(p.parserModel.segements, seg)
	p.parserModel.conditionOpeningBraceCount++
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
	seg := &segement{GoCodeBlock, s}
	p.parserModel.segements = append(p.parserModel.segements, seg)
	p.parserModel.segmentIndex = index + openningBraceIndex + 1
	p.parserModel.conditionOpeningBraceCount++
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
	seg := &segement{GoCodeBlock, code}
	p.parserModel.segements = append(p.parserModel.segements, seg)
	p.parserModel.segmentIndex = index + openningBraceIndex + 1
	p.parserModel.conditionOpeningBraceCount++
}

func (p *RazorParserEngine) handleCloseBrace(index int) {
	l := index - p.parserModel.segmentIndex
	p.handleString(p.parserModel.segmentIndex, l)
	seg := &segement{GoCodeBlock, string(p.src[index])}
	p.parserModel.segements = append(p.parserModel.segements, seg)
	p.parserModel.segmentIndex = index + 1
	p.parserModel.conditionOpeningBraceCount--
}

func (p *RazorParserEngine) handleVariable(index int) {
	stringRemain := p.srcString[index:]
	r, _ := regexp.Compile(variableRegex)
	variableString := r.FindString(stringRemain)
	p.parserModel.segmentIndex = index + len(variableString)
	seg := &segement{Variable, variableString}
	p.parserModel.segements = append(p.parserModel.segements, seg)
}

func (p *RazorParserEngine) handleCommentLine(index int) {
	stringRemain := p.srcString[index:]
	endIndex := strings.Index(stringRemain, "\n")
	p.parserModel.segmentIndex = index + endIndex + 1
	content := "<!--" + stringRemain[2:endIndex] + "-->"
	seg := &segement{String, content}
	p.parserModel.segements = append(p.parserModel.segements, seg)
}

func (p *RazorParserEngine) handleCommentSegment(index int) {
	endTag := string(commentTag) + string(Symbol)
	stringRemain := p.srcString[index:]
	endIndex := strings.Index(stringRemain, endTag)
	p.parserModel.segmentIndex = index + endIndex + 2
	content := "<!--" + stringRemain[1:endIndex] + "-->"
	seg := &segement{String, content}
	p.parserModel.segements = append(p.parserModel.segements, seg)
}

func (p *RazorParserEngine) getLine(startIndex int) (content string, endIndex int) {
	stringRemain := p.srcString[startIndex:]
	endIndex = strings.Index(stringRemain, "\n")
	content = stringRemain[0:endIndex]
	return
}

func (p *RazorParserEngine) handleImports(index int) {
	startIndex := index + len(importTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.parserModel.importParts = strings.Replace(lineContent, ";", "\n", -1)
	p.parserModel.segmentIndex = startIndex + endIndex
}

func (p *RazorParserEngine) handleModel(index int) {
	startIndex := index + len(modelTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.parserModel.modelName = lineContent
	p.parserModel.segmentIndex = startIndex + endIndex
}

func (p *RazorParserEngine) handleLayout(index int) {
	startIndex := index + len(layoutTag)
	lineContent, endIndex := p.getLine(startIndex)
	p.parserModel.layoutPath = strings.Trim(lineContent, " ")
	p.parserModel.segmentIndex = startIndex + endIndex
}

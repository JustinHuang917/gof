package razor

import (
	//"fmt"
	"github.com/justinhuang917/gof/goftool/parser"
	"strings"
)

var (
	literalType        = "literal"
	codeType           = "code"
	expressionType     = "expression"
	import_tag         = "model"
	helper_tag         = "helper"
	layout_declare_tag = "layout"
	model_declare_tag  = "import"
	codeKeywords       = []string{"if", "for", "else"}
)

type Block struct {
	blockType string
	content   string
}
type Declaration struct {
	modelName string
	imports   []string
	layout    string
}

type RazorParser struct {
	Blocks       []Block
	Helpers      []Block
	codeParser   *codeParser
	markupParser *markupParser
	Declaration  *Declaration
	source       string
}

func NewRazorParser(content string) *RazorParser {
	p := &RazorParser{}
	p.markupParser = &markupParser{parser: p}
	p.codeParser = &codeParser{parser: p}
	p.Blocks = make([]Block, 0, 10)
	p.Helpers = make([]Block, 0, 10)
	p.Declaration = &Declaration{}
	p.source = content
	return p
}
func (p *RazorParser) parseMarkupBlcok(source string) {
	p.markupParser.parseBlock(source)
}

func (p *RazorParser) parseCodeBlock(source string) {
	p.codeParser.parseBlock(source)
}

func (p *RazorParser) pushBlock(blockType, content string) {
	b := &Block{blockType: blockType, content: content}
	p.Blocks = append(p.Blocks, *b)
}

func (p *RazorParser) pushHelper(blocks []Block) {
	for _, b := range blocks {
		p.Helpers = append(p.Helpers, b)
	}
}

func (p *RazorParser) genRazorOutput() string {
	output := ""
	temp := ""
	for _, hb := range p.Helpers {
		temp = p.parseBlock(hb)
		output += temp
	}

	for _, b := range p.Blocks {
		temp = p.parseBlock(b)
		output += temp
	}
	return output
}

func (p *RazorParser) parseBlock(b Block) string {
	temp := ""
	switch b.blockType {
	case literalType:
		temp = strings.Replace(b.content, "\r", "\\r", -1)
		temp = strings.Replace(temp, "\n", "\\n", -1)
		temp = strings.Replace(temp, `"`, `\"`, -1)
		temp = parser.Writeout_begin + "\"" + temp + "\"" + parser.Writeout_end
		temp += "\n"
	case codeType:
		temp = b.content
		temp += "\n"
	case expressionType:
		temp = parser.Writeout_begin + b.content + parser.Writeout_end
		temp += "\n"
	}
	return temp
}

func (p *RazorParser) Parse() *parser.ParseResult {
	p.markupParser.parseBlock(p.source)
	result := &parser.ParseResult{}
	result.Imports = strings.Join(p.Declaration.imports, "\n")
	result.LayoutPath = p.Declaration.layout
	result.ModelTypeName = p.Declaration.modelName
	result.OutPutContent = p.genRazorOutput()
	return result
}

type markupParser struct {
	parser *RazorParser
}

func (m *markupParser) isValidEmailChar(code uint8) bool {
	//code := int(char)
	if code >= 48 && code <= 57 {
		return true
	}
	if code >= 65 && code <= 90 {
		return true
	}
	if code >= 97 && code <= 122 {
		return true
	}
	return false
}

func (m *markupParser) isValidTransition(source string, index int) bool {
	if index == 0 {
		return true
	}
	if index == (len(source) - 1) {
		return false
	}

	if m.isValidEmailChar(source[index-1]) && m.isValidEmailChar(source[index+1]) {
		return false
	}
	if rune(source[index-1]) == '@' || rune(source[index+1]) == '@' {
		return false
	}
	return true
}

func (m *markupParser) nextTransition(source string) int {
	for i := 0; i < len(source); i++ {
		if rune(source[i]) == '@' && m.isValidTransition(source, i) {
			return i
		}
	}
	return -1
}

func (m *markupParser) parseBlock(source string) {
	var next = m.nextTransition(source)
	if next == -1 {
		m.parser.pushBlock("literal", source)
		return
	}

	var markup = source[0:next]
	m.parser.pushBlock(literalType, markup)
	m.parser.parseCodeBlock(source[next:])
}

type codeParser struct {
	parser *RazorParser
}

func (c *codeParser) isKeyword(word string) bool {
	for _, w := range codeKeywords {
		if w == word {
			return true
		}
	}
	return false
}
func (c *codeParser) nextChar(source string, findChar rune) int {
	for i, c := range source {
		if c == findChar {
			return i
		}
	}
	return -1
}

func (c *codeParser) endExplicitBlock(source string) int {
	return c.endBlock(source, '(', ')')
}

func (c *codeParser) endCodeBlock(source string) int {
	return c.endBlock(source, '{', '}')
}

func (c *codeParser) endBlock(source string, startChar, endChar rune) int {
	scope := 0
	var cur rune
	quoteChar := ' '
	for i := 0; i < len(source); i++ {
		cur = rune(source[i])
		if cur == '"' {
			if quoteChar == ' ' {
				quoteChar = cur
			} else if quoteChar == cur {
				quoteChar = ' '
			}
		}
		if cur == startChar && quoteChar == ' ' {
			scope++
		}
		if cur == endChar && quoteChar == ' ' {
			scope--
			if scope == 0 {
				return i
			}
		}
	}

	return -1
}

func appendStr(str string, char rune) string {
	return str + string(char)
}

func (c *codeParser) acceptBrace(source string, brace rune) string {
	if len(source) == 0 {
		return ""
	}
	sbrace := int(brace)
	if len(source) < 0 {
		return ""
	}
	qchr := 0
	ebrace := 0
	output := ""
	if sbrace == 40 {
		ebrace = 41
	}
	if sbrace == 91 {
		ebrace = 93
	}
	if rune(source[0]) != brace {
		return ""
	}
	scopes := 0
	//for i := 0; i < len(source); i++ {
	for _, cur := range source {
		//cur := rune(source[i])
		cde := int(cur)
		if cde == sbrace {
			if qchr == 0 {
				scopes++
			}
			output = appendStr(output, cur)
		} else if cde == ebrace {
			if qchr == 0 {
				scopes--
			}
			output = appendStr(output, cur)
			if scopes == 0 {
				break
			} else {
				if qchr == cde {
					qchr = 0
				} else if cde == 34 || cde == 39 {
					qchr = cde
				}
				output = appendStr(output, cur)
			}
		}
	}
	return output
}

func (c *codeParser) acceptIdentifier(source string) string {
	if len(source) == 0 {
		return ""
	}
	output := ""
	for i, cur := range source {
		//cur := rune(source[i])
		cde := int(cur)
		if i == 0 {
			if cde == 36 || cde == 95 || (cde >= 65 && cde <= 90) || (cde >= 97 && cde <= 122) { // $_A-Za-z
				output = appendStr(output, cur)
			} else {
				return ""
			}
		} else {
			if cde == 36 || cde == 95 || (cde >= 65 && cde <= 90) || (cde >= 97 && cde <= 122) { // $_A-Za-z
				output = appendStr(output, cur)
			} else {
				break
			}
		}
	}
	return output
}

func (c *codeParser) parseBlock(source string) {
	if source[0] != '@' {
		c.parser.parseMarkupBlcok(source)
		return
	}
	next := rune(source[1])
	if next == ':' {
		c.parseLine(source)
		return
	}
	if next == '(' {
		//e.g. @(model.name)
		c.parseExplicitExpression(source)
		return
	}

	if next == '{' {
		//e.g.@{ var age = 27; }
		c.parseCodeBlock(source)
		return
	}

	c.parseExpressionBlock(source)
}

func (c *codeParser) parseCodeBlock(source string) {
	end := c.endCodeBlock(source)
	if end == -1 {
		panic("Unterminated code block.")
	}
	code := source[2:end]
	c.parser.pushBlock(codeType, code)
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseLine(source string) {
	end := c.nextChar(source, '\n')
	if end == -1 {
		end = len(source) - 1
	}
	line := source[2:end]
	c.parser.parseMarkupBlcok(line)
	if end != -1 {
		c.parser.parseMarkupBlcok(source[end:])
	}

}
func (c *codeParser) parseExpression(source string) {
	block := source[1:]
	expr := c.readExpression(block)
	if expr == "" {
		c.parser.parseMarkupBlcok(block)
	} else {
		if expr == helper_tag {
			c.parseHelper(source)
		} else if expr == import_tag {
			c.parseModelDec(source)
		} else if expr == model_declare_tag {
			c.parseImportDec(source)

		} else if expr == layout_declare_tag {
			c.parseLayoutDec(source)
		} else {
			c.parser.pushBlock(expressionType, expr)
			c.parser.parseMarkupBlcok(source[len(expr)+1:])
		}
	}
}
func (c *codeParser) parseExplicitExpression(source string) {
	end := c.endExplicitBlock(source)
	if end == -1 {
		panic("Untermined explicit expression.")
	}
	expr := source[2:end]
	c.parser.pushBlock(expressionType, expr)
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseExpressionBlock(source string) {
	nextScope := c.nextChar(source, '{')
	if nextScope > -1 {
		identifier := strings.Split(source[1:], " ")[0]
		//fmt.Println("identifier:", identifier)
		if c.isKeyword(identifier) {
			c.parseKeyword(identifier, source)
			return
		}
	}
	c.parseExpression(source)
}

func (c *codeParser) parseKeyword(keyword string, source string) {
	switch keyword {
	case "if":
		c.parseIfBlock(source)
		break
	case "for":
		c.parseSimpleBlock(source)
		break
	case "else":
		c.parseSimpleBlock(source)
		break
	default:
		c.parser.parseMarkupBlcok(source[1:])
	}

}

func (c *codeParser) parseIfBlock(source string) {
	end := c.endCodeBlock(source)
	if end == -1 {
		panic("Unterminated if  block.")
	}
	start := c.nextChar(source, '{')
	statement := source[1 : start+1]
	c.parser.pushBlock(codeType, statement)
	innerBlock := source[start+1 : end]
	c.parser.parseMarkupBlcok(innerBlock)
	c.parser.pushBlock(codeType, "}")
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseSimpleBlock(source string) {
	end := c.endCodeBlock(source)
	if end == -1 {
		panic("Unterminated code block.")
	}
	start := c.nextChar(source, '{')
	statement := source[1 : start+1]
	c.parser.pushBlock(codeType, statement)
	innerBlock := source[start+1 : end]
	c.parser.parseMarkupBlcok(innerBlock)
	c.parser.pushBlock(codeType, "}")
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) readExpression(source string) string {
	if len(source) == 0 {
		return ""
	}
	output := ""
	state := 1
	i := 0
	for true {
		if state == 1 {
			id := c.acceptIdentifier(source[i:])
			if id == "" {
				break
			}
			output = output + id
			i = i + len(id)
			state = 2
		} else if state == 2 {
			if source[i] == '(' || source[i] == '[' {
				brace := c.acceptBrace(source[i:], rune(source[i]))
				if brace == "" {
					break
				}
				output = output + brace
				i = i + len(brace)
			} else {
				state = 3
			}
		} else if state == 3 {
			if source[i] == '.' {
				state = 4
				i++
				continue
			}
			break
		} else if state == 4 {
			id := c.acceptIdentifier(source[i:])
			if id == "" {
				break
			}
			output = appendStr(output, '.')
			state = 1
		}
	}
	return output
}

func (c *codeParser) parseHelper(source string) {
	end := c.endCodeBlock(source)
	if end == -1 {
		panic("Unterminated helper block.")
	}
	start := c.nextChar(source, '{')
	name := source[7 : start+1]
	l1 := len(c.parser.Blocks)
	names := strings.Split(name, "(")
	if len(names) != 2 {
		panic("Invalid helper func  Declaration:" + name)
	}
	funcName := names[0] + ":=func"
	funcDec := funcName + "(" + names[1]
	c.parser.pushBlock(codeType, funcDec)
	innerBlock := source[start+1 : end]
	c.parser.parseMarkupBlcok(innerBlock)
	c.parser.pushBlock(codeType, "}")
	l2 := len(c.parser.Blocks)
	helper := c.parser.Blocks[l1:l2]
	c.parser.pushHelper(helper)
	c.parser.Blocks = c.parser.Blocks[0:l1]
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseModelDec(source string) {
	end, dec := c.parseDeclaration(source, import_tag)
	c.parser.Declaration.modelName = dec
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseImportDec(source string) {
	end, dec := c.parseDeclaration(source, model_declare_tag)
	c.parser.Declaration.imports = strings.Split(dec, ",")
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseLayoutDec(source string) {
	end, dec := c.parseDeclaration(source, layout_declare_tag)
	c.parser.Declaration.layout = dec
	c.parser.parseMarkupBlcok(source[end+1:])
}

func (c *codeParser) parseDeclaration(source, decName string) (end int, dec string) {
	end = c.nextChar(source, '\n')
	l := len(layout_declare_tag) + 1
	dec = source[l:end]
	dec = strings.TrimSpace(dec)
	return
}

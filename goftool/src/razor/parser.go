package razor

import (
	"./util"
	//"fmt"
	//"strings"
)

const (
	open_tag           = "@{"
	close_tag          = "}"
	out_tag            = "="
	end_server_tag     = "}"
	import_tag         = "import"
	model_declare_tag  = "model"
	layout_declare_tag = "layout"
	else_tag           = "}else{"
	else_if_tag        = "}elseif{"
	helper_tag         = "helper"
	gohtml_ext         = ".gohtml"
	defaultModel       = "core.NilModel"
)

const (
	Html_Unkonwn = iota
	Html_Text
	Html_WhiteSpace
	Html_NewLine
	Html_OpenAngle    //'<'
	Html_Bang         //!
	Html_Solius       // '/'
	Html_QuestionMark // '?'
	Html_DoubleHyphen //'--'
	Html_LeftBracket  //[
	Html_CloseAngle   //>
	Html_RightBracket //]
	Html_Equals       //=
	Html_DoubleQuote  //"
	Html_SingleQuote  //'
	Html_Transition   //@
	Html_Colon
	RazorComment           // //
	RazorCommentStar       // */
	RazorCommentTransition //@*  
)

var (
	writeout_begin = "d.Writeout(out,"
	writeout_end   = ")\n"
	line           = 0
)

func main() {
	src := `@*{ 
				i:=0
				for i<3{
					<p>@i</p>
				}
			}*`

	t := NewTokenizer(src)
	//sm := &StateMechine{}
	t = t
}

type Tokenizer struct {
	StateMechine
	SrcString            string
	Src                  []rune
	Buffer               []rune
	CurrentStartPosition int
	CurrentPosition      int
	CurrentChar          rune
	Length               int
	Readed               []rune
}

func NewTokenizer(src string) *Tokenizer {
	t := &Tokenizer{}
	t.SrcString = src
	t.Src = util.ToChars(src)
	t.Buffer = make([]rune, 0, 10)
	t.CurrentPosition = 0
	t.CurrentStartPosition = 0
	t.Length = len(t.Src)
	return t
}

func (t *Tokenizer) PeekPos() int {
	j := t.CurrentPosition + 1
	if j > t.Length-1 {
		return -1
	}
	return j
}

func (t *Tokenizer) MoveNext() {
	t.Readed = append(t.Readed, t.CurrentChar)
}

func (t *Tokenizer) Peek() rune {
	t.MoveNext()
	return t.CurrentChar
}

func (t *Tokenizer) IsEnd() bool {
	return t.PeekPos() == -1
}

func (t *Tokenizer) TakeCurrent() {
	if t.IsEnd() {
		return
	}
	t.Buffer = append(t.Buffer, t.CurrentChar)
}

func (t *Tokenizer) WhiteSpace() *Symbol {
	for util.IsWhiteSpace(t.CurrentChar) {
		t.TakeCurrent()
	}
	return t.EndSymbol(Html_WhiteSpace)
}

func (t *Tokenizer) Newine() *Symbol {
	flag := t.CurrentChar == '\r'
	t.TakeCurrent()
	if flag && t.CurrentChar == '\n' {
		t.TakeCurrent()
	}
	return t.EndSymbol(Html_NewLine)
}

func (t *Tokenizer) StartSymbol() {
	t.Buffer = make([]rune, 0, 10)
	t.CurrentStartPosition = t.CurrentPosition
}

func (t *Tokenizer) HaveContent() bool {
	return len(t.Buffer) > 0
}

func (t *Tokenizer) EndSymbol(symbolType int) *Symbol {
	var sym *Symbol
	if t.HaveContent() {
		sym = CreatSymbol(t.CurrentStartPosition, string(t.Buffer), symbolType)
	}
	t.StartSymbol()
	return sym
}

func (t *Tokenizer) AtSymbol() bool {
	c := t.CurrentChar
	return c == '<' || c == '!' || c == '/' || c == '?' || c == '[' || c == '>' || c == '=' || c == '"' || c == '\'' || c == '@' || (c == '-' && t.Peek() == '-')
}

func (t *Tokenizer) HtmlSymbol() *Symbol {
	c := t.CurrentChar
	t.TakeCurrent()
	switch c {
	case '<':
		return t.EndSymbol(Html_OpenAngle)
	case '!':
		return t.EndSymbol(Html_Bang)
	case '/':
		return t.EndSymbol(Html_Solius)
	case '?':
		return t.EndSymbol(Html_QuestionMark)
	case '[':
		return t.EndSymbol(Html_LeftBracket)
	case '>':
		return t.EndSymbol(Html_CloseAngle)
	case ']':
		return t.EndSymbol(Html_RightBracket)
	case '=':
		return t.EndSymbol(Html_Equals)
	case '"':
		return t.EndSymbol(Html_DoubleQuote)
	case '\'':
		return t.EndSymbol(Html_SingleQuote)
	case '-':
		t.TakeCurrent()
		return t.EndSymbol(Html_DoubleHyphen)
	default:
		return t.EndSymbol(Html_Unkonwn)
	}
	return t.EndSymbol(Html_Unkonwn)
}

func (t *Tokenizer) HtmlText() *StateResult {
	var prev rune
	//prev = '\0'
	for t.IsEnd() && util.IsNewLineOrWhiteSpace(t.CurrentChar) && !t.AtSymbol() {
		prev = t.CurrentChar
		t.TakeCurrent()
	}
	if t.CurrentChar == '@' {
		next := t.Peek()
		if util.IsNumberOrLetter(prev) && util.IsNumberOrLetter(next) {
			t.TakeCurrent()
			return t.Stay()
		}
	}
	return t.Transition(t.EndSymbol(Html_Text), func() *StateResult {
		return t.Html()
	})
}

func (t *Tokenizer) Html() *StateResult {
	if util.IsWhiteSpace(t.CurrentChar) {
		return t.StaySymbol(t.WhiteSpace())
	} else if util.IsNewLine(t.CurrentChar) {
		return t.StaySymbol(t.Newine())
	} else if t.AtSymbol() {
		return t.StaySymbol(t.HtmlSymbol())
	}
	return t.TransitionNewState(func() *StateResult {
		return t.HtmlText()
	})

}

type Symbol struct {
	StartPosition int
	Content       string
	SymbolType    int
}

func CreatSymbol(start int, content string, symbolType int) *Symbol {
	return &Symbol{start, content, symbolType}
}

type State func() *StateResult

type StateResult struct {
	Next      State
	HasOutput bool
	Symbol    *Symbol
}

type StateMechine struct {
	CurrentState State
}

func (s *StateMechine) Turn() *Symbol {
	var result *StateResult
	if s.CurrentState != nil {
		for result != nil && result.HasOutput {
			result := s.CurrentState()
			s.CurrentState = result.Next
		}
		if result == nil {
			return nil
		}
		return result.Symbol
	}
	return nil
}

func (s *StateMechine) Stop() *StateResult {
	return nil
}

func (s *StateMechine) Stay() *StateResult {
	return s.TransitionNewState(s.CurrentState)
}

func (s *StateMechine) StaySymbol(symbol *Symbol) *StateResult {
	return &StateResult{s.CurrentState, true, symbol}
}

func (s *StateMechine) TransitionNewState(next State) *StateResult {
	r := &StateResult{}
	r.Next = next
	return r
}

func (s *StateMechine) Transition(symbol *Symbol, next State) *StateResult {
	return &StateResult{Next: next, HasOutput: true, Symbol: symbol}
}

package goftool

import (
	"errors"
	"fmt"
	"github.com/justinhuang917/gof/goftool/razor"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	gohtml_ext         = ".gohtml"
	gorazor_ext        = ".gorazor"
	defaultModel       = "gofcore.NilModel"
	viewDir            = "view/html"
)

var (
	writeout_begin = "d.Writeout(out,"
	writeout_end   = ")\n"
	line           = 0
)

type CompileResult struct {
	Imports       string
	ViewName      string
	LayoutPath    string
	ModelTypeName string
	OutPutContent string
	RouteName     string
}

func (c *CompileResult) isNeedLayout() bool {
	return c.LayoutPath != ""
}

type Comiler struct {
	FileName   string
	FilePath   string
	bytes      []byte
	LineCount  int
	OutputPath string
	Result     *CompileResult
}

func NewCompiler(path string, OutputPath string) (*Comiler, error) {
	compiler := &Comiler{}
	fi, err := os.Open(path)
	if err == nil {
		compiler.FileName = fi.Name()
		compiler.FilePath = path
		compiler.OutputPath = OutputPath
		compiler.Result = &CompileResult{}
		compiler.LineCount = 0
		compiler.bytes, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		compiler.Result.ViewName = getViewName(path)
		compiler.Result.LayoutPath = ""
		compiler.Result.ModelTypeName = defaultModel
	} else {
		return nil, err
	}
	defer fi.Close()
	return compiler, nil
}

func (c *Comiler) RazorCompile() {
	parserEngine := razor.NewRazorParserEngine(string(c.bytes))
	parserEngine.Parse()
	c.Result.OutPutContent = genRazorOutput(parserEngine.ParserModel.Segements)
	c.Result.ModelTypeName = parserEngine.ParserModel.ModelName
	c.Result.LayoutPath = parserEngine.ParserModel.LayoutPath
	c.Result.Imports = parserEngine.ParserModel.ImportParts

}

func genRazorOutput(segements []*razor.Segement) string {
	output := ""
	temp := ""
	for _, seg := range segements {
		temp = ""
		switch seg.SegementType {
		case razor.String:
			temp = strings.Replace(seg.Content, "\r", "\\r", -1)
			temp = strings.Replace(temp, "\n", "\\n", -1)
			temp = writeout_begin + "\"" + temp + "\"" + writeout_end
			temp += "\n"
		case razor.GoCodeBlock:
			temp = seg.Content
			temp += "\n"
		case razor.Variable:
			temp = writeout_begin + seg.Content + writeout_end
			temp += "\n"
		}
		output += temp
	}
	return output
}

func (c *Comiler) Compile(needOutput bool, isRazor bool) (err error) {
	if isRazor {
		c.RazorCompile()
	} else {
		var result = ""
		src := string(c.bytes)
		for _, span := range strings.Split(src, open_tag) {
			innerSpans := strings.Split(span, close_tag)
			l := len(innerSpans)
			if l == 1 {
				result += c.html(innerSpans[0])
			} else {
				code0 := innerSpans[0]
				code1 := innerSpans[1]
				result += c.logic(code0)
				if len(code1) > 0 {
					result += c.html(code1)
				}
			}
		}
		c.Result.OutPutContent = result
	}
	c.Result.RouteName = getRouteName(c.FilePath)
	var html string
	if b := c.Result.isNeedLayout(); b {
		html, err = c.Result.formatHtmlWithLayout(isRazor)
	} else {
		html, err = c.Result.formatHtmlNoLayout()
	}
	if err == nil {
		if needOutput {
			err = c.outputWriteFile(html)
		}
	}
	return
}

func (c *CompileResult) formatHtmlNoLayout() (html string, err error) {
	html = fmt.Sprintf(NoLayoutViewTemplate,
		c.Imports,
		c.ViewName,
		c.RouteName,
		c.ViewName,
		c.ViewName,
		c.ModelTypeName,
		c.OutPutContent)
	return
}

func (c *CompileResult) formatHtmlWithLayout(isRazor bool) (html string, err error) {
	layoutResult, err2 := compileLayout(c.LayoutPath, isRazor)
	if err2 == nil {
		layoutBodyResult := c.OutPutContent
		html = fmt.Sprintf(LayoutviewTemplate,
			c.Imports,
			c.ViewName,
			c.RouteName,
			c.ViewName,
			c.ViewName,
			c.ModelTypeName,
			layoutBodyResult,
			layoutResult.OutPutContent)

	} else {
		err = err2
	}
	return
}

func (c *Comiler) outputWriteFile(html string) error {
	os.Create(c.OutputPath)
	outputFile, err2 := os.OpenFile(c.OutputPath, os.O_RDWR, 0666)
	if err2 == nil {
		io.WriteString(outputFile, html)
	}
	defer outputFile.Close()
	return err2
}

func (c *Comiler) html(code string) string {
	c.LineCount += len(strings.Split("\n", code))
	code = strings.Replace(code, "\r", "\\r", -1)
	code = strings.Replace(code, "\n", "\\n", -1)
	code = strings.Replace(code, `"`, `\"`, -1)
	code = writeout_begin + "\"" + code + "\"" + writeout_end
	code = code + "\n"
	return code
}

func (c *Comiler) logic(code string) string {
	c.LineCount += len(strings.Split("\n", code))
	keyword := c.getKeyword(code)
	var parse = c.getParser(keyword)
	code = parse(code)
	code = code + "\n"
	return code
}

func (c *Comiler) getKeyword(code string) (keyword string) {
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

func (c *Comiler) getParser(keyword string) func(code string) string {
	var parser func(code string) string
	switch {
	case keyword == import_tag:
		parser = func(code string) string {
			c.Result.Imports = strings.Join(strings.Split(code, " ")[1:], " ")
			c.Result.Imports = strings.Replace(c.Result.Imports, ";", "\n", -1)
			return ""
		}
	case keyword == model_declare_tag:
		parser = func(code string) string {
			c.Result.ModelTypeName = strings.Join(strings.Split(code, " ")[1:], " ")
			return ""
		}
	case keyword == layout_declare_tag:
		parser = func(code string) string {
			c.Result.LayoutPath = strings.Join(strings.Split(code, " ")[1:], " ")
			c.Result.LayoutPath = strings.Replace(c.Result.LayoutPath, " ", "", -1)
			return ""
		}
	case keyword == else_tag:
		parser = func(code string) string {
			return code
		}
	case keyword == else_if_tag:
		parser = func(code string) string {
			return code
		}
	case keyword == end_server_tag:
		parser = func(code string) string {
			return "}\n"
		}
	case keyword == out_tag:
		parser = func(code string) string {
			code = strings.Replace(code, out_tag, "", -1)
			code = writeout_begin + code + writeout_end
			return code
		}
	case keyword == helper_tag:
		parser = func(code string) string {
			code = strings.Replace(code, helper_tag, "", -1)
			code = strings.Replace(code, "func", "", -1)
			codes := strings.Split(code, "(")
			endCode := strings.Join(codes[1:], "")
			return codes[0] + ":=func(" + endCode
		}
	default:
		parser = func(code string) string {
			return code
		}
	}
	return parser
}

func readUntil(text string, char rune) string {
	count := 0
	var chars = make([]rune, 1, 10)
	for _, c := range text {
		if c == char {
			break
		}
		chars = append(chars, c)
		count++
	}
	return string(chars)
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

var layoutResultCache map[string]*CompileResult

func init() {
	layoutResultCache = make(map[string]*CompileResult)
}

func compileLayout(layoutPath string, isRazor bool) (compileResult *CompileResult, err error) {
	if layoutResultCache[layoutPath] == nil {
		c, err1 := NewCompiler(layoutPath, "")
		if err1 != nil {
			err = err1
			return
		}
		err = c.Compile(false, isRazor)
		compileResult = c.Result
		layoutResultCache[layoutPath] = compileResult
	} else {
		compileResult = layoutResultCache[layoutPath]
	}
	return
}

type visitor struct{}

func (self *visitor) DoCompile(path string, outputDir string, f os.FileInfo, buidingArgs ...string) error {
	if f == nil {
		return nil
	}
	if f.Name() == "view.go" {
		return nil
	}
	if f.IsDir() {
		return nil
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		return nil
	} else {
		var err error
		extName := filepath.Ext(f.Name())
		if extName == gohtml_ext || extName == gorazor_ext {
			isRazor := false
			if extName == gorazor_ext {
				isRazor = true
			}
			OutputPath := filepath.Join(outputDir, getViewName(path)+".go")
			if OutputPath[0] != '.' {
				OutputPath = "./" + OutputPath
			}
			c, err1 := NewCompiler(path, OutputPath)
			if err1 != nil {
				return err1
			}
			fmt.Println("Parsing:", path)
			err = c.Compile(true, isRazor)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Building:", c.OutputPath)
				br := Building(c.OutputPath, buidingArgs...)
				if br != nil && br.Err != nil {
					//err = genCompileError(path, string(br.Out))
					errMsg := fmt.Sprintf("%s:%s", path, string(br.Out))
					err = errors.New(errMsg)
				}
			}
			if err == nil {
				fmt.Println("Success!!!")
				fmt.Println("------------------------------------------------------")
			}
			return err
		}
	}
	return nil
}

func Compile(dirPath string, outputDir string, buidingArgs ...string) error {
	begin := time.Now()
	v := &visitor{}
	//Clear()
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		err1 := v.DoCompile(path, outputDir, f, buidingArgs...)
		return err1
	})
	if err != nil {
		//Clear()
		fmt.Println("-----------------------------------------------")
		fmt.Println("Failed!!!")
	} else {
		fmt.Println("-----------------------------------------------")
		fmt.Println("Build All View Success!!!")
	}
	duration := time.Since(begin)
	fmt.Printf("Time Cost: %vs\n", duration.Seconds())
	return err
}

func Clear(dir string) error {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		} else if (f.Mode() & os.ModeSymlink) > 0 {
			return nil
		} else {
			if ok, _ := filepath.Match("V_*.go", f.Name()); ok {
				err = os.Remove(path)
				return err
			}
		}
		return nil
	})
	return err
}

func getViewName(filepath string) string {
	startIndex := strings.Index(filepath, viewDir)
	startIndex = startIndex + len(viewDir)
	f := filepath[startIndex:]
	viewName := strings.Replace(f, ".", "", -1)
	viewName = strings.Replace(viewName, "/", "_", -1)
	viewName = "V" + strings.Replace(viewName, gohtml_ext[1:], "", -1)
	viewName = strings.Replace(viewName, gorazor_ext[1:], "", -1)
	return viewName
}

func getRouteName(filepath string) string {
	startIndex := strings.Index(filepath, viewDir)
	startIndex = startIndex + len(viewDir)
	f := filepath[startIndex:]
	routeName := strings.Replace(f, ".", "", -1)
	routeName = strings.Replace(routeName, gohtml_ext[1:], "", -1)
	routeName = strings.Replace(routeName, gorazor_ext[1:], "", -1)

	return routeName
}

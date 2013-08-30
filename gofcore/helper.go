package gofcore

import (
	"fmt"
	"github.com/JustinHuang917/gof/gofcore/cfg"
	"strconv"
	"strings"
)

const (
	textboxTextFormat  = "<input type=\"input\" id=\"%s\" value=\"%s\" %s />"
	checkboxTextFormat = "<input type=\"check\" id=\"%s\" value=\"%s\" %s />"
)

func UrlControllerAction(controllerName, actionName string, values map[string]string) string {
	url := cfg.AppConfig.AppPath + controllerName + "/" + actionName
	if values != nil {
		args := make([]string, 0, len(values))
		for k, v := range values {
			kvStr := fmt.Sprintf("%s=%s", k, v)
			args = append(args, kvStr)
		}
		queryString := strings.Join(args, "&")
		url = fmt.Sprintf("%s?%s", url, queryString)
	}
	return url
}

func attrsToString(attrs map[string]string) string {
	s := ""
	for k, v := range attrs {
		s += fmt.Sprintf("%s=\"%s\"", k, v)
	}
	return s
}

func DisplayTextBox(id string, value string, attrs map[string]string) string {
	return fmt.Sprintf(textboxTextFormat, id, value, attrsToString(attrs))
}

func DisplayCheckbox(id string, value string, attrs map[string]string) string {
	return fmt.Sprintf(textboxTextFormat, id, value, attrsToString(attrs))
}
func DisplayInt(i int) string {
	return strconv.Itoa(i)
}

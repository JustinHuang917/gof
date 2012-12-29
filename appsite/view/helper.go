package view

import (
	"fmt"
	"strconv"
)

const (
	textboxTextFormat  = "<input type=\"input\" id=\"%s\" value=\"%s\" %s />"
	checkboxTextFormat = "<input type=\"check\" id=\"%s\" value=\"%s\" %s />"
)

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

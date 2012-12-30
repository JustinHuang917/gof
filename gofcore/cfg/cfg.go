package cfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppPath        string `json:"AppPath"`
	Port           string `json:"Port"`
	DefaultPath    string `json:"DefaultPath"`
	NotFoundPath   string `json:"NotFoundPath"`
	GofSessionId   string `json:GofSessionId`
	SessionMode    string `json:SessionMode`
	SessionExpires int    `json:SessionExpires`
	EnableSession  bool   `json:EnableSession`
}

var AppConfig *Config

func init() {
	fmt.Println("Init Config")
	err := load("./cfg.json")
	if err != nil {
		fmt.Println("Init Config Error:", err)
	}
}
func load(cfgPath string) error {
	file, err := os.Open(string(cfgPath))
	if err != nil {
		return err
	}
	AppConfig = &Config{}
	dec := json.NewDecoder(file)
	//fmt.Println(dec)
	if err = dec.Decode(AppConfig); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(AppConfig)
	//fmt.Println(AppConfig)
	return nil
}

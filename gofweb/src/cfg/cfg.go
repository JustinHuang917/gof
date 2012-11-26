package cfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppPath      string `json:"AppPath"`
	Port         string `json:"Port"`
	DefaultPath  string `json:"DefaultPath"`
	NotFoundPath string `json:"NotFoundPath"`
}

var AppConfig *Config

func init() {
	Load("cfg.json")
}
func Load(cfgPath string) error {
	file, err := os.Open(string(cfgPath))
	if err != nil {
		return err
	}
	AppConfig = &Config{}
	dec := json.NewDecoder(file)
	fmt.Println(dec)
	if err = dec.Decode(AppConfig); err != nil {
		return err
	}
	fmt.Println(AppConfig)
	return nil
}

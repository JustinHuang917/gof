// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {

	//Listen Ptah
	AppPath string `json:"AppPath"`
	//Listen Port
	Port string `json:"Port"`

	GofSessionId   string `json:GofSessionId`
	SessionMode    string `json:SessionMode`
	SessionExpires int    `json:SessionExpires`
	EnableSession  bool   `json:EnableSession`

	DebugMode string `json:DebugMode`

	HandlerSortings map[string]int `json:HandlerSortings`

	/*
		RouteRules is a list of RouteRule
		RouteRule is a map that  key is the route rule,and  value
		is a default values map
		route rule format:/Order/{id:[0-9]+}
	*/
	RouteRules []map[string]map[string]string `json:RouteRules`

	RootPath string `json:RootPath`

	StaticDirs []string `json:StaticDirs`

	//Special Settings
	AppSettings map[string]string `json:AppSettings`
}

var AppConfig *Config

func init() {
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
	if err = dec.Decode(AppConfig); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

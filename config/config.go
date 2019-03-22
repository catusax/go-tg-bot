package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/coolrc136/go-tg-bot/cmd"
)

type Conf struct {
	Hook      string `json:"hook"`
	Token     string `json:"token"`
	Tuling    string `json:"tuling"`
	Projectid string `json:"projectid"`
	Jsonfile  string `json:"jsonfile"`
	Lang      string `json:"lang"`
}

var Token, Webhook, Tuling_token, Projectid, Jsonfile, Lang string

func ReadConf() {
	var Config Conf
	configFile, err := os.Open(cmd.ConfPath)
	if err != nil {
		fmt.Printf("opening config file", err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&Config); err != nil {
		fmt.Printf("parsing config file", err.Error())
	}

	Token = Config.Token
	Webhook = Config.Hook
	Tuling_token = Config.Tuling
	Projectid = Config.Projectid
	Jsonfile = Config.Jsonfile
	Lang = Config.Lang
}

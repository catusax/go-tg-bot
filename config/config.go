package config

import (
    "os"
    "encoding/json"
    "fmt"
    "github.com/coolrc136/go-tg-bot/cmd"
)

type Conf struct{
    Hook	string `json:"hook"`
    Token   string `json:"token"`
    Tuling string `json:"tuling"`
}


var Token string
var Webhook string
var Tuling_token string

func ReadConf(){
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
}

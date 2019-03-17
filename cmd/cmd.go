package cmd

import (
    "os"
)

var ConfPath string = "conf/config.json"

func init(){

    if len(os.Args) > 2 && os.Args[1] == "-c" {
            ConfPath = os.Args[2]
        }
}

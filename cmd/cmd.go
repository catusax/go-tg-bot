package cmd

import (
	"flag"
    "os"
    "fmt"
)

var Sethook = flag.Bool("n", false, "reset webhook,default:false")
var Token = flag.String("token", "aaa", "your bot token")
var Webhook = flag.String("server", "https://www.google.com:8443", "your webhook address with port")
var Tuling_token = flag.String("tuling", "sgdger", "your tuling robot apikey")

func init(){
    flag.Parse() //解析参数

    if len(os.Args) == 1 || os.Args[1] == "-h" {
            flag.PrintDefaults()
            fmt.Println("  -h   help message")
            os.Exit(0)
        }
}

package tuling

import (
	"github.com/go-tg-bot/cmd"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PostMsg struct { //post数据
	Key    string `json:"key"`
	Info   string `json:"info"`
	Userid string `json:"userid"`
}
type Answer struct { //接收的回答
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Tuling(info string, userid string) string {

	var post PostMsg
	post.Key = *cmd.Tuling_token
	post.Info = info
	post.Userid = userid
	msg, _ := json.Marshal(post)
	fmt.Println(string(msg))

	res, err := http.Post("http://www.tuling123.com/openapi/api", "application/json;charset=utf-8", bytes.NewBuffer(msg)) //进行post请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	var ans Answer
	err = json.Unmarshal([]byte(content), &ans) //解析json

	fmt.Println(string(content))
	return ans.Text
}

/*
func main() {
	fmt.Println(tuling(os.Args[1], "47725"))
} */

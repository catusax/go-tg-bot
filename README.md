# go-tg-bot
telegram bot build with golang

# 使用前准备
从letsencrypt 申请证书，将`full_chained.pem`改为`cert.pem`,`privkey.pem`改为`key.pem`放入`cert`目录
如果需要使用自签名证书，请参考go-telegram-bot-api/telegram-bot-api修改代码

# 编译运行
```
go get github.com/coolrc136/go-tg-bot
go build github.com/coolrc136/go-tg-bot
```

可以使用help参数查看使用说明

```
./go-tg-bot -n=true -server="https://yourserver.com:8443/" -token="your-token" -tuling="tuling_robot_apikey"
```

# docker

```
docker run -itd --restart=always --name tgbot -p 8443:8443 -e SERVER="https://yourserver.com:8443/" -e TOKEN="your_token" -e TULING="tuling_robot_apikey" -v $PWD/cert:/bot/cert coolrc/tgbot
```

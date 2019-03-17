# go-tg-bot
telegram bot build with golang

# 使用前准备
从letsencrypt 申请证书，将`full_chained.pem`改为`cert.pem`,`privkey.pem`改为`key.pem`放入`cert`目录
如果需要使用自签名证书，请参考go-telegram-bot-api/telegram-bot-api修改代码

填写配置文件：`config.json`
获取dialogflow的认证文件：[https://dialogflow.com/docs/reference/v2-auth-setup](https://dialogflow.com/docs/reference/v2-auth-setup)
填写dialogflow的配置：`df.json`

# 编译运行
```
go get github.com/coolrc136/go-tg-bot
go build github.com/coolrc136/go-tg-bot
```

```
./go-tg-bot -c config.json
```

# docker

```
docker run -itd --restart=always --name tgbot -p 8443:8443 \
    -v $PWD/cert:/bot/cert \
    -v $PWD/conf:/bot/conf \
    coolrc/tgbot
```

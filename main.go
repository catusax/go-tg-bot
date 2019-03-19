package main

import (
	"fmt"
	"github.com/coolrc136/go-tg-bot/config"
	. "github.com/coolrc136/go-tg-bot/handle"
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)



func main() {
    config.ReadConf()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Webhook + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert/cert.pem", "cert/key.pem", nil)
	fmt.Println("listening")

	Handle(&updates, bot)
}
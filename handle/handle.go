package handle

import (
        "fmt"
        "log"
        "github.com/coolrc136/go-tg-bot/tuling"
        "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/coolrc136/go-tg-bot/dialogflow"
)

func Handle(updates *tgbotapi.UpdatesChannel,bot *tgbotapi.BotAPI) {
        for update := range *updates { //消息处理
                log.Printf("%+v\n", update)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
                if update.Message.IsCommand() {
                        //解析命令
                        switch update.Message.Command() {
                        case "start":
                                msg.Text = "my name is van."
                        case "sayhi":
                                msg.Text = "Hi :)"
                        case "status":
                                msg.Text = "I'm ok."
                        default:
                                msg.Text = "I don't know that command"
                        }

                } else {
            		dfresponse := dp.processNLP(update.Message.Text, fmt.Sprintf("%d", update.Message.Chat.ID))
            		if dfresponse.Intent != "Default Fallback Intent" {
                		msg.Text = dfresponse.Result
            		}
            		else {
                		msg.Text = tuling.Tuling(update.Message.Text, fmt.Sprintf("%d", update.Message.Chat.ID))
            		}
        		}

                if _, err := bot.Send(msg); err != nil {
                        log.Panic(err)
                }
        }
}


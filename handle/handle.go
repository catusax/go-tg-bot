package handle

import (
        "fmt"
        "log"

        "github.com/coolrc136/go-tg-bot/tuling"
        "github.com/coolrc136/go-tg-bot/config"
        df "github.com/coolrc136/go-tg-bot/dialogflow"

        "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handle(updates *tgbotapi.UpdatesChannel,bot *tgbotapi.BotAPI) {
        
        //init
        Tuling :=  tuling.NewApi(config.Tuling_token)

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
            		intent,dfmsg :=df.Df(update.Message.Text, fmt.Sprintf("%d", update.Message.Chat.ID))
            		if intent != "Default Fallback Intent" {
                		msg.Text = dfmsg
            		} else {
                		msg.Text = Tuling.GetMsg(update.Message.Text, fmt.Sprintf("%d", update.Message.Chat.ID))
            		}
        		}

                if _, err := bot.Send(msg); err != nil {
                        log.Panic(err)
                }
        }
}


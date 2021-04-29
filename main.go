package main

import (
	"github.com/Skycoin/git-telegram-bot/config"
	"github.com/Skycoin/git-telegram-bot/errutil"
	"github.com/Skycoin/git-telegram-bot/githandler"
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
)

func main() {
	l := log.New(os.Stdout, "skygit-bot", log.LstdFlags)

	cfg, err := config.NewBotConfig()
	if err != nil {
		l.Fatal(err)
	}

	bot, err := tb.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		l.Fatal(errutil.ErrCreatingBot.Desc(err))
	}

	bot.Debug = true
	updateConfig := tb.NewUpdate(0)

	updateConfig.Timeout = 10

	updates := bot.GetUpdatesChan(updateConfig)

	var chatId int64
	stopCh := make(chan struct{})
	var previousEventId string
	var currentEventId string
	ticker := time.NewTicker(10 * time.Second)

	for update := range updates {
		userIsAdmin := false
		chatConfig := update.Message.Chat.ChatConfig()
		admins, adminErr := bot.GetChatAdministrators(tb.ChatAdministratorsConfig{ChatConfig: chatConfig})
		if adminErr != nil {
			continue
		}
		for _, admin := range admins {
			if update.Message.From.ID == admin.User.ID {
				userIsAdmin = true
			}
		}
		if !userIsAdmin {
			continue
		}

		if update.Message.IsCommand() {
			chatId = update.Message.Chat.ID
			switch update.Message.Command() {
			case "start": // starts the poller
				msg := tb.NewMessage(chatId, "starting Skycoin poll github events...")
				if _, e := bot.Send(msg); err != nil {
					l.Printf("error sending start message: %v", e)
					continue
				}
				go func() {
					for {
						select {
						case <-stopCh:
							ticker.Stop()
							break
						case <-ticker.C:
							if err = githandler.HandleStartCommand(
								previousEventId,
								currentEventId,
								l, cfg.TargetOrgUrl,
								func(s string) error {
									msg = tb.NewMessage(chatId, s)
									if _, e := bot.Send(msg); err != nil {
										return e
									}
									return nil
								},
							); err != nil {
								l.Print(err)
								continue
							}
						}
					}
				}()
			case "stop": // stops it
				stopCh <- struct{}{}
				msg := tb.NewMessage(chatId, "stopping bot, you can use /reset then /start command to start it again")
				if _, err = bot.Send(msg); err != nil {
					l.Printf("error sending message: %v", err)
				}
			case "help": // displays help message
				msg := tb.NewMessage(chatId, `
Hi, here's my list of commands:
	/start: starts polling events from github
	/stop: stops the poller
	/reset: resets the poller, use with /start after /stop to restart polling event.
`)
				if _, err = bot.Send(msg); err != nil {
					l.Printf("error sending message: %v", err)
				}
			case "reset":
				ticker = time.NewTicker(10 * time.Second)
			}
		}

	}

}

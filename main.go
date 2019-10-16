package main

import (
	"bip-dev/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	srvc := service.NewService(bot)

	updates, err := srvc.UpdatesChan()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for {
		select {
		case <-sigs:
			wg.Wait()

		case update := <-updates:
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			wg.Add(1)
			go func(upd tgbotapi.Update) {
				if err := srvc.Handle(upd); err != nil {
					log.Println(err)
				}
				wg.Done()
			}(update)
		}
	}
}

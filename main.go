package main

import (
	"bip-dev/service"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	rds := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	if err := rds.Ping().Err(); err != nil {
		log.Fatal(err)
	}

	var pgql *sqlx.DB //todo

	languageBundle := i18n.NewBundle(language.English)
	languageBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.en.toml"))
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.ru.toml"))

	application := service.NewApplication(rds, pgql, languageBundle, log.New(os.Stdout, "app", log.LstdFlags|log.Lshortfile))

	var wg sync.WaitGroup
	go func() {
		for update := range updates {
			if update.Message == nil && update.CallbackQuery == nil {
				continue
			}
			wg.Add(1)
			go handle(application.NewFactory(update), bot, &wg)
		}
	}()

	<-sigs
	bot.StopReceivingUpdates()
	wg.Wait()
}

func handle(factory *service.AbstractFactory, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	defer wg.Done()
	factory.SetLocalizer()
	message := factory.CreateMessage()
	factory.SaveArgs()
	bot.Send(message)
	factory.SaveReply()
}

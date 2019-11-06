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

	langBundle := i18n.NewBundle(language.English)
	langBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	langBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.en.toml"))
	langBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.ru.toml"))

	logger := log.New(os.Stdout, "app", log.LstdFlags|log.Lshortfile)

	application := service.NewApplication(rds, pgql, langBundle, logger)

	var wg sync.WaitGroup
	go func() {
		for update := range updates {
			if isValidData(update) {
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

func isValidData(update tgbotapi.Update) bool {
	return (update.Message == nil || update.Message.From == nil) &&
		(update.CallbackQuery == nil || update.CallbackQuery.Data == "" ||
			update.CallbackQuery.Message == nil || update.CallbackQuery.Message.From == nil)
}

func handle(factory *service.AbstractFactory, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	defer wg.Done()
	factory.SetLocalizer()
	err := factory.Answer(bot)
	if err != nil {
		factory.Log(err)
		return
	}
	factory.SaveReply()
}

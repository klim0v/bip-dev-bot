package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var languageBundle = i18n.NewBundle(language.English)

func init() {
	languageBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.en.toml"))
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.ru.toml"))
}

type CommandHandler struct {
	ChatID    int64
	Localizer *i18n.Localizer
	Command   string
	Args      string
	Reply     string
}

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

	var wg sync.WaitGroup
	go func() {
		for update := range updates {
			wg.Add(1)
			if update.Message != nil {
				chatLang := "ru"
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				commandHandler := &CommandHandler{
					ChatID:    update.Message.Chat.ID,
					Localizer: i18n.NewLocalizer(languageBundle, chatLang, update.Message.From.LanguageCode),
					Command:   update.Message.Command(),
					Args:      update.Message.CommandArguments(),
					Reply:     "",
				}
				if !update.Message.IsCommand() {
					commandHandler.Command = "todo" //todo: get last bot message by update.Message.Chat.ID from redis
					commandHandler.Args = update.Message.Text
				}

				msg := commandHandler.createReplyMessage()

				//todo: save commandHandler.Reply to redis

				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}

			}
			wg.Done()

		}
	}()

	<-sigs
	bot.StopReceivingUpdates()
	wg.Wait()

}

func (command *CommandHandler) createReplyMessage() (msg tgbotapi.MessageConfig) {
	switch command.Command {
	case "help":
		command.Reply = "Help"
		msg = tgbotapi.NewMessage(
			command.ChatID,
			command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}),
		)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
				tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
				tgbotapi.NewInlineKeyboardButtonData("3", "3"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("4", "4"),
				tgbotapi.NewInlineKeyboardButtonData("5", "5"),
				tgbotapi.NewInlineKeyboardButtonData("6", "6"),
			),
		)
	}
	return
}

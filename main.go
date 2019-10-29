package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
)

var languageBundle = i18n.NewBundle(language.English)

func init() {
	languageBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.en.toml"))
	languageBundle.MustLoadMessageFile(filepath.Join(".", "locales", "actions.ru.toml"))
}

type Command struct {
	ChatID    int64
	Localizer *i18n.Localizer
	Command   string
	Args      string
	Reply     string
}

func formCommand(update tgbotapi.Update, rds *redis.Client) *Command {
	if update.Message != nil {
		lastCommandResult := rds.Get("lastCommand:" + strconv.Itoa(int(update.Message.Chat.ID)))
		lastCommand, err := lastCommandResult.Result()
		log.Println("get lastCommand", lastCommand)
		if err != nil {
			if err == redis.Nil {
				lastCommand = ""
			} else {
				log.Println(err)
			}
		}
		return createCommandByMessage(update.Message, lastCommand)
	}

	if update.CallbackQuery != nil {
		//todo: return createCommandByCallbackQuery(update.Message)
	}

	return nil
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

	rds := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	if ping := rds.Ping(); ping.Err() != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	go func() {
		for update := range updates {
			wg.Add(1)
			go handle(rds, update, bot, &wg)
		}
	}()

	<-sigs
	bot.StopReceivingUpdates()
	wg.Wait()
}

func handle(rds *redis.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	defer wg.Done()
	command := formCommand(update, rds)
	if command == nil {
		log.Println("command is nil")
		return
	}
	handleCommand(command, bot)
	rds.SetNX("lastCommand:"+strconv.Itoa(int(command.ChatID)), command.Reply, 0)
}

func createCommandByMessage(message *tgbotapi.Message, lastCommand string) *Command {
	chatLang := "ru"
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	command := &Command{
		ChatID:    message.Chat.ID,
		Localizer: i18n.NewLocalizer(languageBundle, chatLang, message.From.LanguageCode),
		Command:   message.Command(),
		Args:      message.CommandArguments(),
		Reply:     "",
	}
	if !message.IsCommand() {
		command.Command = lastCommand
		command.Args = message.Text
	}
	return command
}

func handleCommand(command *Command, bot *tgbotapi.BotAPI) {
	msg := command.createReplyMessage()
	if msg.ChatID == 0 {
		log.Println("unknown command")
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func (command *Command) createReplyMessage() tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	log.Println("createReplyMessage:command", command.Command)
	switch command.Command {
	case "":
		fallthrough
	case "help":
		command.Reply = "help"
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
	default:
		return msg
	}
	return msg
}

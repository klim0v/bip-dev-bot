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

func formCommand(message *tgbotapi.Message, rds *redis.Client) *Command {
	lastCommandResult := rds.Get("lastCommand:" + strconv.Itoa(int(message.Chat.ID)))
	lastCommand, err := lastCommandResult.Result()
	log.Println("get lastCommand", lastCommand)
	if err != nil {
		if err == redis.Nil {
			lastCommand = ""
		} else {
			log.Println(err)
		}
	}
	return createCommandByMessage(message, lastCommand)
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
	if update.Message != nil {
		command := formCommand(update.Message, rds)
		if command == nil {
			log.Println("command is nil")
			return
		}
		handleCommand(command, bot)
		rds.SetNX("lastCommand:"+strconv.Itoa(int(command.ChatID)), command.Reply, 0)
		return
	}

	//bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
	localizer := i18n.NewLocalizer(languageBundle, "ru", update.CallbackQuery.From.LanguageCode)
	message := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: update.CallbackQuery.Data}))
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check_sell"}), "check_sell"), //todo: make constants
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "cancel"),
		),
	)
	message.ReplyMarkup = &markup

	_, err := bot.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	//todo: save last action to redis
	return
}

func createCommandByMessage(message *tgbotapi.Message, lastCommand string) *Command {
	chatLang := "ru" //todo
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
	//log.Println("createReplyMessage:command", command.Command)
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
				tgbotapi.NewInlineKeyboardButtonData("by_coin", "by_coin"),
				tgbotapi.NewInlineKeyboardButtonData("sell_coin", "sell_coin"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("my_orders", "my_orders"),
			),
		)
	case "by_coin":
	case "sell_coin":
	case "my_orders":
	default:
		return msg
	}
	return msg
}

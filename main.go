package main

import (
	"fmt"
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
	"strings"
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
	if err := rds.Ping().Err(); err != nil {
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
		rds.Set("lastCommand:"+strconv.Itoa(int(command.ChatID)), command.Reply, 0)
		return
	}

	if update.CallbackQuery.Data != "" {
		fields := strings.Fields(update.CallbackQuery.Data)
		var args string
		if len(fields) > 1 {
			args = strings.Join(fields[1:], " ")
		}
		command := Command{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			Localizer: i18n.NewLocalizer(languageBundle, "ru", update.CallbackQuery.From.LanguageCode),
			Command:   fields[0],
			Args:      args,
			Reply:     "",
		}

		var message tgbotapi.Chattable
		switch command.Command {
		case "by_coin":
			command.Reply = "send_minter_address"
			msg := tgbotapi.NewEditMessageText(command.ChatID, update.CallbackQuery.Message.MessageID, command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}))
			markup := sendMinterAddressMarkup(command.Localizer, minterAddresses())
			msg.ReplyMarkup = &markup
			msg.ParseMode = "markdown"
			message = msg
		//case strings.HasPrefix(update.CallbackQuery.Data, "send_minter_address") && len(strings.Fields(update.CallbackQuery.Data)) == 2:
		case "use_minter_address":
			//todo save command.Args

			command.Reply = "send_email_address"
			msg := tgbotapi.NewEditMessageText(command.ChatID, update.CallbackQuery.Message.MessageID, command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}))
			markup := sendEmailAddressMarkup(command.Localizer, emailAddresses())
			msg.ReplyMarkup = &markup
			msg.ParseMode = "markdown"
			message = msg
		case "use_email_address":
			//todo save command.Args

			command.Reply = "send_btc"
			msg := tgbotapi.NewEditMessageText(command.ChatID, update.CallbackQuery.Message.MessageID,
				fmt.Sprintf(command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}), 0.0184, -24.28, 516841, 4.00, btcAddresses()),
			)
			markup := sendBTCAddressMarkup(command.Localizer)
			msg.ReplyMarkup = &markup
			msg.ParseMode = "markdown"
			message = msg
		case "help":
			command.Reply = "help"
			msg := tgbotapi.NewEditMessageText(command.ChatID, update.CallbackQuery.Message.MessageID,
				command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}),
			)
			markup := helpMarkup(command.Localizer)
			msg.ReplyMarkup = &markup
			msg.ParseMode = "markdown"
			message = msg
		default:
			return
		}

		if err := rds.Set("lastCommand:"+strconv.Itoa(int(command.ChatID)), command.Reply, 0).Err(); err != nil {
			log.Println(err)
			return
		}

		_, err := bot.Send(message)
		if err != nil {
			log.Fatal(err)
		}

		//todo: save last action to redis
		return
	}
}

func emailAddresses() []string {
	return []string{"klim0v-sergey@yandex.ru"}
}

func minterAddresses() []string {
	return []string{"Mx00000000000000000000000000000987654321"}
}

func btcAddresses() string {
	return "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN"
}

func helpMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "by_coin"}), "by_coin"),
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "sell_coin"}), "sell_coin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "my_orders"}), "my_orders"),
		),
	)
}

func sendBTCAddressMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), "check_sell"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "by_coin"), //todo get next step from relations map
		),
	)
}

func sendEmailAddressMarkup(localizer *i18n.Localizer, addresses []string) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address, fmt.Sprintf("use_email_address %s", address))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "by_coin"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func sendMinterAddressMarkup(localizer *i18n.Localizer, addresses []string) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address, fmt.Sprintf("use_minter_address %s", address))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "help"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
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

func (command *Command) createReplyMessage() (msg tgbotapi.MessageConfig) {
	switch command.Command {
	case "":
		fallthrough
	case "help":
		command.Reply = "help"
		msg = tgbotapi.NewMessage(
			command.ChatID,
			command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}),
		)
		markup := helpMarkup(command.Localizer)
		msg.ReplyMarkup = &markup
	case "send_minter_address":
		//todo save and use command.Args
		command.Reply = "send_email_address"
		msg = tgbotapi.NewMessage(
			command.ChatID,
			command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}),
		)
		msg.ReplyMarkup = sendMinterAddressMarkup(command.Localizer, emailAddresses())
	case "send_email_address":
		//todo save and use command.Args
		command.Reply = "send_btc"
		msg = tgbotapi.NewMessage(
			command.ChatID,
			fmt.Sprintf(command.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: command.Reply}), 0.0184, -24.28, 516841, 4.00, btcAddresses()),
		)
		msg.ReplyMarkup = sendBTCAddressMarkup(command.Localizer)
	case "sell_coin":
	case "my_orders":
	default:
		return msg
	}
	msg.ParseMode = "markdown"
	return msg
}

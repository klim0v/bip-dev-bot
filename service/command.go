package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CommandFactory struct {
	Message
	Command    string
	Args       string
	Repository *Repository
}

func (command *CommandFactory) SaveArgs() error {
	var err error

	switch command.Command {
	case "send_minter_address":
		//todo save and use command.Args
		return err
	case "send_email_address":
		//todo save and use command.Args
		return err
	}

	return nil
}

func (command *CommandFactory) CreateMessage() tgbotapi.Chattable {
	var msg tgbotapi.MessageConfig
	switch command.Command {
	case "":
		fallthrough
	case "help":
		command.Message.reply = "help"
		msg = tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
		)
		msg.ReplyMarkup = helpMarkup(command.Localizer())
	case "send_minter_address":
		command.Message.reply = "send_email_address"
		msg = tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
		)
		msg.ReplyMarkup = sendEmailAddressMarkup(command.Localizer(), command.Repository.emailAddresses())
	case "send_email_address":
		command.Message.reply = "send_btc"
		msg = tgbotapi.NewMessage(
			command.ChatID(),
			fmt.Sprintf(command.translateReply(), 0.0184, -24.28, 516841, 4.00, command.Repository.btcAddresses()),
		)
		msg.ReplyMarkup = sendBTCAddressMarkup(command.Localizer())
	default:

		return msg
	}
	msg.ParseMode = "markdown"

	return msg
}

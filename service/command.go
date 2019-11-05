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

type HelpCommandFactory struct {
	CommandFactory
}

func (command *HelpCommandFactory) Answer() (tgbotapi.Chattable, error) {
	command.Message.reply = "help"
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
	)
	msg.ReplyMarkup = helpMarkup(command.Localizer())
	msg.ParseMode = "markdown"
	return msg, nil
}

type SendMinterAddressCommandFactory struct {
	CommandFactory
}

func (command *SendMinterAddressCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidEmailAddress(command.Args) {
		command.Message.reply = "send_email_address" // todo: may be to move saveReply there and remove this line
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: "send_email_address_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	command.Message.reply = "send_email_address"
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
	)
	msg.ReplyMarkup = sendEmailAddressMarkup(command.Localizer(), command.Repository.emailAddresses())
	msg.ParseMode = "markdown"
	return msg, nil
}

type SendEmailAddressCommandFactory struct {
	CommandFactory
}

func (command *SendEmailAddressCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidMinterAddress(command.Args) {
		command.Message.reply = "send_minter_address" // todo: may be to move saveReply there and remove this line
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: "send_minter_address_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	command.Message.reply = "send_btc"
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.translate(command.reply), 0.0184, -24.28, 516841, 4.00, command.Repository.btcAddresses()),
	)
	msg.ReplyMarkup = sendBTCAddressMarkup(command.Localizer())
	msg.ParseMode = "markdown"
	return msg, nil
}

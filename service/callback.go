package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CallbackFactory struct {
	Message
	MessageUpdateID int
	Command         string
	Args            string
	Repository      *Repository
}

func (callback *CallbackFactory) SaveArgs() error {
	var err error

	switch callback.Command {
	case "use_minter_address":
		//todo save and use command.Args
		return err
	case "use_email_address":
		//todo save and use command.Args
		return err
	}

	return nil
}

func (callback *CallbackFactory) CreateMessage() tgbotapi.Chattable {
	var message tgbotapi.Chattable
	switch callback.Command {
	case "by_coin":
		callback.Message.reply = "send_minter_address"

		msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translateReply())
		markup := sendMinterAddressMarkup(callback.Localizer(), callback.Repository.minterAddresses())
		msg.ReplyMarkup = &markup
		msg.ParseMode = "markdown"
		message = msg

	case "use_minter_address":
		callback.Message.reply = "send_email_address"

		msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translateReply())
		markup := sendEmailAddressMarkup(callback.Localizer(), callback.Repository.emailAddresses())
		msg.ReplyMarkup = &markup
		msg.ParseMode = "markdown"
		message = msg

	case "use_email_address":
		callback.Message.reply = "send_btc"

		sprintf := fmt.Sprintf(callback.translateReply(), 0.0184, -24.28, 516841, 4.00, callback.Repository.btcAddresses())
		msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, sprintf)
		markup := sendBTCAddressMarkup(callback.Localizer())
		msg.ReplyMarkup = &markup
		msg.ParseMode = "markdown"
		message = msg

	case "help":
		callback.Message.reply = "help"

		msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translateReply())
		markup := helpMarkup(callback.Localizer())
		msg.ReplyMarkup = &markup
		msg.ParseMode = "markdown"
		message = msg

	default:
		return nil
	}

	return message
}

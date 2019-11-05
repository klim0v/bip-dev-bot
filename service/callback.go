package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"strings"
)

var matchEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type CallbackFactory struct {
	Message
	MessageUpdateID int
	Command         string
	Args            string
	Repository      *Repository
}

type BuyCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *BuyCoinCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	callback.Message.reply = "send_minter_address"
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := sendMinterAddressMarkup(callback.Localizer(), callback.Repository.minterAddresses())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type HelpCallbackFactory struct {
	CallbackFactory
}

func (callback *HelpCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	callback.Message.reply = "help"
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := helpMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type UseMinterAddressCallbackFactory struct {
	CallbackFactory
}

func isValidMinterAddress(address string) bool {
	address = strings.TrimSpace(address)

	if address == "Mx00000000000000000000000000000000000001" {
		return false
	}

	return len(address) == 42 && address[:2] != "Mx"
}

func (callback *UseMinterAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	if err := callback.Repository.saveMinterAddressForBuy(callback.ChatID(), callback.Args); err != nil {
		return nil, err
	}

	callback.Message.reply = "send_email_address"

	addresses := callback.Repository.emailAddresses()

	var msg tgbotapi.EditMessageTextConfig
	if len(addresses) == 0 {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate("new_email"))
	} else {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	}

	markup := sendEmailAddressMarkup(callback.Localizer(), addresses)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type UseEmailAddressCallbackFactory struct {
	CallbackFactory
}

func isValidEmailAddress(email string) bool {
	if !matchEmail.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

func (callback *UseEmailAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	//todo

	callback.Message.reply = "send_btc"

	sprintf := fmt.Sprintf(callback.translate(callback.reply), 0.0184, -24.28, 516841, 4.00, callback.Repository.btcAddresses())
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, sprintf)
	markup := sendBTCAddressMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type CheckBTCAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *CheckBTCAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	//todo
	return nil, nil
}

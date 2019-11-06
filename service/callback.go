package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strconv"
)

type CallbackFactory struct {
	Message
	MessageUpdateID int
	Command         string
	Args            string
	Repository      *Repository
}

type SellCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *SellCoinCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	callback.Message.reply = enterCoinName
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := selectCoinNameMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type BuyCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *BuyCoinCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	callback.Message.reply = selectMinterAddress
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := selectMinterAddressMarkup(callback.Localizer(), callback.Repository.minterAddresses())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type HelpCallbackFactory struct {
	CallbackFactory
}

func (callback *HelpCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	callback.Message.reply = help
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := helpMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type UseMinterAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseMinterAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {

	if err := callback.Repository.saveMinterAddressForSell(callback.ChatID(), callback.Args); err != nil {
		return nil, err
	}

	callback.Message.reply = selectEmailAddress

	addresses := callback.Repository.emailAddresses()

	var msg tgbotapi.EditMessageTextConfig
	if len(addresses) == 0 {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate("new_email"))
	} else {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	}

	markup := selectEmailAddressMarkup(callback.Localizer(), addresses)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type UseEmailAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseEmailAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	emailAddressID, err := strconv.Atoi(callback.Args)
	if err != nil {
		return nil, err
	}

	if err := callback.Repository.saveEmailAddressForBuy(callback.ChatID(), emailAddressID); err != nil {
		return nil, err
	}

	callback.Message.reply = sendDepositForBuyBIP

	sprintf := fmt.Sprintf(callback.translate(callback.reply), 0.0184, -24.28, 516841, 4.00, "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN")
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, sprintf)
	markup := sendBTCAddressMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"
	return msg, nil
}

type CheckSellCallbackFactory struct {
	CallbackFactory
}

func (callback *CheckSellCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	return nil, nil
}

type UseBitcoinAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseBitcoinAddressCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidBitcoinAddress(callback.Args) {
		callback.Message.reply = useBitcoinAddress
		msg := tgbotapi.NewMessage(
			callback.ChatID(),
			callback.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: callback.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	callback.Message.reply = sendYourCoins
	msg := tgbotapi.NewMessage(
		callback.ChatID(),
		fmt.Sprintf(callback.translate(callback.reply), "BIP", "BIP", "www.example.com"),
	)
	//msg.ReplyMarkup = todo
	msg.ParseMode = "markdown"
	return msg, nil
}

type ToDoCallbackFactory struct {
	CallbackFactory
}

func (callback *ToDoCallbackFactory) Answer() (tgbotapi.Chattable, error) {
	panic("todo")
	return nil, nil
}

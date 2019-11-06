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

func (callback *SellCoinCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.reply = enterCoinName
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := selectCoinNameMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type BuyCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *BuyCoinCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.reply = selectMinterAddress
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := selectMinterAddressMarkup(callback.Localizer(), callback.Repository.minterAddresses())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type HelpCallbackFactory struct {
	CallbackFactory
}

func (callback *HelpCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.reply = help
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.translate(callback.reply))
	markup := helpMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type UseMinterAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseMinterAddressCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {

	if err := callback.Repository.saveMinterAddressForSell(callback.ChatID(), callback.Args); err != nil {
		return err
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

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type UseEmailAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseEmailAddressCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	emailAddressID, err := strconv.Atoi(callback.Args)
	if err != nil {
		return err
	}

	if err := callback.Repository.saveEmailAddressForBuy(callback.ChatID(), emailAddressID); err != nil {
		return err
	}

	callback.Message.reply = sendDepositForBuyBIP

	sprintf := fmt.Sprintf(callback.translate(callback.reply), 0.0184, -24.28, 516841, 4.00, "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN")
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, sprintf)
	markup := sendBTCAddressMarkup(callback.Localizer())
	msg.ReplyMarkup = &markup
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type CheckSellCallbackFactory struct {
	CallbackFactory
	QueryID string
}

func (callback *CheckSellCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.reply = waitDepositCoin
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, fmt.Sprintf(callback.translate(callback.reply), "BIP"))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

type UseBitcoinAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseBitcoinAddressCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !isValidBitcoinAddress(callback.Args) {
		callback.Message.reply = useBitcoinAddress
		msg := tgbotapi.NewMessage(
			callback.ChatID(),
			callback.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: callback.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}

	link := "www.example.com"

	callback.Message.reply = sendYourCoins
	msg1 := tgbotapi.NewMessage(
		callback.ChatID(),
		fmt.Sprintf(callback.translate(callback.reply), "BIP", "BIP", "www.example.com"),
	)
	msg1.ReplyMarkup = shareMarkup(callback.Localizer(), link)
	msg1.ParseMode = "markdown"

	if _, err := bot.Send(msg1); err != nil {
		return err
	}

	msg2 := tgbotapi.NewMessage(
		callback.ChatID(),
		"`Mx36f2a491683f7667006f9208bfd6220d551c05fd`",
	)
	msg2.ReplyMarkup = sendYourCoinsMarkup(callback.Localizer())
	msg2.ParseMode = "markdown"

	if _, err := bot.Send(msg2); err != nil {
		return err
	}

	return nil
}

type CheckSendDepositCallbackFactory struct {
	CallbackFactory
	QueryID string
}

func (callback *CheckSendDepositCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.reply = waitDepositCoin
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, fmt.Sprintf(callback.translate(callback.reply), "BIP"))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

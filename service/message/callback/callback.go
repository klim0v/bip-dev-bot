package callback

import (
	"bip-dev/service/message"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type CallbackFactory struct {
	message.Message
	MessageUpdateID int
	Command         string
	Args            string
	Repository      *message.Repository
}

type SellCoinCallbackFactory struct {
	CallbackFactory
}

func (callback *SellCoinCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.EnterCoinName)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.SelectCoinNameMarkup(callback.Localizer())
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
	callback.Message.SetReply(message.EnterMinterAddress)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.SelectMinterAddressMarkup(callback.Localizer(), callback.Repository.MinterAddresses())
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
	callback.Message.SetReply(message.Help)
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	markup := message.HelpMarkup(callback.Localizer())
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

	minterID, err := strconv.Atoi(callback.Args)
	if err != nil {
		return err
	}

	if err := callback.Repository.SaveBuyMinterAddress(callback.ChatID(), minterID); err != nil {
		return err
	}

	callback.Message.SetReply(message.EnterEmailAddress)

	addresses := callback.Repository.EmailAddresses()

	var msg tgbotapi.EditMessageTextConfig
	if len(addresses) == 0 {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate("new_email"))
	} else {
		msg = tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, callback.Translate(callback.Reply()))
	}

	markup := message.SelectEmailAddressMarkup(callback.Localizer(), addresses)
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

	if err := callback.Repository.SaveEmailAddressForBuy(callback.ChatID(), emailAddressID); err != nil {
		return err
	}

	callback.Message.SetReply(message.SendDepositForBuyBIP)

	sprintf := fmt.Sprintf(callback.Translate(callback.Reply()), 0.0184, -24.28, 516841, 4.00, "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN")
	msg := tgbotapi.NewEditMessageText(callback.ChatID(), callback.MessageUpdateID, sprintf)
	markup := message.SendBTCAddressMarkup(callback.Localizer())
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
	callback.Message.SetReply(message.WaitDepositCoin)
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, fmt.Sprintf(callback.Translate(callback.Reply()), "BIP"))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

type UseBitcoinAddressCallbackFactory struct {
	CallbackFactory
}

func (callback *UseBitcoinAddressCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	bitcoinAddressID, err := strconv.Atoi(callback.Args)
	if err != nil {
		return err
	}

	if err := callback.Repository.SaveSellBitcoinAddress(callback.ChatID(), bitcoinAddressID); err != nil {
		return err
	}

	coinName, err := callback.Repository.SellCoinName(callback.ChatID())
	if err != nil {
		if err != redis.Nil {
			//todo: logging
		}

		callback.Message.SetReply(message.EnterCoinName)
		msg := tgbotapi.NewMessage(
			callback.ChatID(),
			callback.Translate(callback.Message.Reply()+"_invalid"),
		)
		msg.ParseMode = "markdown"

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}

	link := "www.example.com"

	callback.Message.SetReply(message.SendYourCoins)
	msg1 := tgbotapi.NewMessage(
		callback.ChatID(),
		fmt.Sprintf(callback.Translate(callback.Reply()), coinName, link),
	)
	msg1.ReplyMarkup = message.ShareMarkup(callback.Localizer(), link)
	msg1.ParseMode = "markdown"

	if _, err := bot.Send(msg1); err != nil {
		return err
	}

	msg2 := tgbotapi.NewMessage(
		callback.ChatID(),
		"`Mx36f2a491683f7667006f9208bfd6220d551c05fd`",
	)
	msg2.ReplyMarkup = message.SendYourCoinsMarkup(callback.Localizer())
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
	callback.Message.SetReply(message.WaitDepositCoin)
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, fmt.Sprintf(callback.Translate(callback.Reply()), "BIP"))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

package callback

import (
	"bip-dev/service/message"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

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

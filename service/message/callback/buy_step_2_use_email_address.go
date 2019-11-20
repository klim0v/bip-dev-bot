package callback

import (
	"bip-dev/service/message"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

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

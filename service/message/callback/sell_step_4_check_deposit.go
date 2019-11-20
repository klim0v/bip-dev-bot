package callback

import (
	"bip-dev/service/message"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CheckSellCallbackFactory struct {
	CallbackFactory
	QueryID string
}

func (callback *CheckSellCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.WaitDepositBtc)
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, fmt.Sprintf(callback.Translate(callback.Reply()), "BIP"))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

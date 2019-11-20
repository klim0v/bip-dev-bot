package callback

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CheckSendDepositCallbackFactory struct {
	CallbackFactory
	QueryID string
}

func (callback *CheckSendDepositCallbackFactory) Answer(bot *tgbotapi.BotAPI) error {
	callback.Message.SetReply(message.WaitDepositCoin)
	msg := tgbotapi.NewCallbackWithAlert(callback.QueryID, callback.Translate(callback.Reply()))

	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		return err
	}

	return nil
}

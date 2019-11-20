package command

import (
	"bip-dev/service/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HelpCommandFactory struct {
	CommandFactory
}

func (command *HelpCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	command.Message.SetReply(message.Help)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Translate(command.Message.Reply()),
	)
	msg.ReplyMarkup = message.HelpMarkup(command.Localizer())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

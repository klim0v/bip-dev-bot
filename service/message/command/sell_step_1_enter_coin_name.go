package command

import (
	"bip-dev/service/message"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EnterCoinNameCommandFactory struct {
	CommandFactory
}

func (command *EnterCoinNameCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	if !IsValidCoinName(command.Args) {
		command.Message.SetReply(message.EnterCoinName)
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Translate(command.Message.Reply()+"_invalid"),
		)
		msg.ParseMode = "markdown"

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}

	if err := command.Repository.SaveSellCoinName(command.ChatID(), command.Args); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterPriceCoin)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply())),
	)
	//msg.ReplyMarkup = todo
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

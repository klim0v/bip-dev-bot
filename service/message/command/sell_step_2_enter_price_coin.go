package command

import (
	"bip-dev/service/message"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EnterPriceCoinCommandFactory struct {
	CommandFactory
}

func (command *EnterPriceCoinCommandFactory) Answer(bot *tgbotapi.BotAPI) error {
	coinName, err := command.Repository.SellCoinName(command.ChatID())
	if err != nil {
		if err != redis.Nil {
			//todo: logging
		}

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

	if !IsValidPriceCoin(coinName, command.Args) {
		command.Message.SetReply(message.EnterPriceCoin)
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

	if err := command.Repository.SaveSellPrice(command.ChatID(), command.Args); err != nil {
		return err
	}

	command.Message.SetReply(message.EnterBitcoinAddress)
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.Translate(command.Reply())),
	)
	msg.ReplyMarkup = message.SelectBitcoinMarkup(command.Localizer(), command.Repository.BtcAddresses())
	msg.ParseMode = "markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

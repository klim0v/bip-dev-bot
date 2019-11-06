package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"regexp"
	"strings"
)

var matchCoinName = regexp.MustCompile("^[0-9-A-Z-a-z]{3,10}$")
var matchEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidPriceCoin(address string) bool {
	//todo
	return true
}
func isValidMinterAddress(address string) bool {
	address = strings.TrimSpace(address)

	if address == "Mx00000000000000000000000000000000000001" {
		return false
	}

	return len(address) == 42 && address[:2] != "Mx"
}

func isValidEmailAddress(email string) bool {
	if !matchEmail.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

func isValidCoinName(coinName string) bool {
	if !matchCoinName.MatchString(coinName) {
		return false
	}
	return true
}

func isValidBitcoinAddress(address string) bool {
	//todo
	return true
}

type CommandFactory struct {
	Message
	Command    string
	Args       string
	Repository *Repository
}

type HelpCommandFactory struct {
	CommandFactory
}

func (command *HelpCommandFactory) Answer() (tgbotapi.Chattable, error) {
	command.Message.reply = help
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
	)
	msg.ReplyMarkup = helpMarkup(command.Localizer())
	msg.ParseMode = "markdown"
	return msg, nil
}

type SelectMinterAddressCommandFactory struct {
	CommandFactory
}

func (command *SelectMinterAddressCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidMinterAddress(command.Args) {
		command.Message.reply = selectMinterAddress
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	emailID, err := command.Repository.addMinterAddress(command.ChatID(), command.Args)
	if err != nil {
		return nil, err
	}

	if err := command.Repository.saveEmailAddressForBuy(command.ChatID(), emailID); err != nil {
		return nil, err
	}

	command.Message.reply = selectEmailAddress
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply}),
	)
	msg.ReplyMarkup = selectEmailAddressMarkup(command.Localizer(), command.Repository.emailAddresses())
	msg.ParseMode = "markdown"
	return msg, nil
}

type SelectEmailAddressCommandFactory struct {
	CommandFactory
}

func (command *SelectEmailAddressCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidEmailAddress(command.Args) {
		command.Message.reply = selectEmailAddress
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	emailID, err := command.Repository.addEmailAddress(command.ChatID(), command.Args)
	if err != nil {
		return nil, err
	}

	if err := command.Repository.saveEmailAddressForBuy(command.ChatID(), emailID); err != nil {
		return nil, err
	}

	command.Message.reply = selectBitcoinAddress
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.translate(command.reply), 0.0184, -24.28, 516841, 4.00, command.Repository.btcAddresses()),
	)
	msg.ReplyMarkup = sendBTCAddressMarkup(command.Localizer())
	msg.ParseMode = "markdown"
	return msg, nil
}

type EnterCoinNameCommandFactory struct {
	CommandFactory
}

func (command *EnterCoinNameCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidCoinName(command.Args) {
		command.Message.reply = enterCoinName
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	if err := command.Repository.saveCoinNameForSell(command.ChatID(), command.Args); err != nil {
		return nil, err
	}

	command.Message.reply = enterPriceCoin
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.translate(command.reply)),
	)
	//msg.ReplyMarkup = todo
	msg.ParseMode = "markdown"
	return msg, nil
}

type EnterPriceCoinCommandFactory struct {
	CommandFactory
}

func (command *EnterPriceCoinCommandFactory) Answer() (tgbotapi.Chattable, error) {
	if !isValidPriceCoin(command.Args) {
		command.Message.reply = enterPriceCoin
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	if err := command.Repository.savePriceForSell(command.ChatID(), command.Args); err != nil {
		return nil, err
	}

	command.Message.reply = selectBitcoinAddress
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.translate(command.reply)),
	)
	msg.ReplyMarkup = selectBitcoinMarkup(command.Localizer(), command.Repository.btcAddresses())
	msg.ParseMode = "markdown"
	return msg, nil
}

type SelectBitcoinAddressCommandFactory struct {
	CommandFactory
}

func (command *SelectBitcoinAddressCommandFactory) Answer() (tgbotapi.Chattable, error) { //todo make []Chattable
	if !isValidBitcoinAddress(command.Args) {
		command.Message.reply = selectBitcoinAddress
		msg := tgbotapi.NewMessage(
			command.ChatID(),
			command.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: command.Message.reply + "_invalid"}),
		)
		msg.ParseMode = "markdown"
		return msg, nil
	}

	command.Message.reply = sendYourCoins
	msg := tgbotapi.NewMessage(
		command.ChatID(),
		fmt.Sprintf(command.translate(command.reply), "BIP", "BIP", "www.example.com"),
	)
	//msg.ReplyMarkup = todo
	msg.ParseMode = "markdown"
	return msg, nil
}

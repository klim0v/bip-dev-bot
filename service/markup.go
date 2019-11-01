package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func helpMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "by_coin"}), "by_coin"),
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "sell_coin"}), "sell_coin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "my_orders"}), "my_orders"),
		),
	)
}

func sendBTCAddressMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), "check_sell"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "by_coin"),
		),
	)
}

func sendEmailAddressMarkup(localizer *i18n.Localizer, addresses []string) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address, fmt.Sprintf("use_email_address %s", address))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "by_coin"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func sendMinterAddressMarkup(localizer *i18n.Localizer, addresses []string) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address, fmt.Sprintf("use_minter_address %s", address))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), "help"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

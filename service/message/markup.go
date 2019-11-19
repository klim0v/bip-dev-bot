package message

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func HelpMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: BuyCoin}), BuyCoin),
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: SellCoin}), SellCoin),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: MyOrders}), MyOrders),
		),
	)
}

func SendBTCAddressMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), CheckSell),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), BuyCoin),
		),
	)
}

func SendYourCoinsMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "check"}), CheckSell),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), EnterPriceCoin),
		),
	)
}

func ShareMarkup(localizer *i18n.Localizer, link string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "share"}), link),
		),
	)
}

func SelectBitcoinMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("%s %d", UseBitcoinAddress, address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), EnterPriceCoin),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func SelectEmailAddressMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("%s %d", UseEmailAddress, address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), BuyCoin),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func SelectMinterAddressMarkup(localizer *i18n.Localizer, addresses []Address) tgbotapi.InlineKeyboardMarkup {
	var addressesKeyboard [][]tgbotapi.InlineKeyboardButton
	for _, address := range addresses {
		addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(address.Value, fmt.Sprintf("use_minter_address %d", address.ID))))
	}
	addressesKeyboard = append(addressesKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), Help),
	))
	return tgbotapi.NewInlineKeyboardMarkup(addressesKeyboard...)
}

func SelectCoinNameMarkup(localizer *i18n.Localizer) tgbotapi.InlineKeyboardMarkup {
	var keyboards [][]tgbotapi.InlineKeyboardButton
	keyboards = append(keyboards, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "cancel"}), Help),
	))
	return tgbotapi.NewInlineKeyboardMarkup(keyboards...)
}

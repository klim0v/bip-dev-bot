package message

import "github.com/nicksnyder/go-i18n/v2/i18n"

const (
	CheckSendDeposit     = "check_send_deposit"
	SendDepositForBuyBIP = "send_deposit"
	BuyCoin              = "buy_coin"
	WaitDepositBtc       = "wait_deposit_btc"

	EnterEmailAddress   = "enter_email_address"
	EnterMinterAddress  = "enter_minter_address"
	EnterBitcoinAddress = "enter_bitcoin_address"
	NewEmailAddress     = "new_email_address"
	NewMinterAddress    = "new_minter_address"
	NewBitcoinAddress   = "new_bitcoin_address"
	UseEmailAddress     = "use_email_address"
	UseMinterAddress    = "use_minter_address"
	UseBitcoinAddress   = "use_bitcoin_address"

	SellCoin       = "sell_coin"
	EnterCoinName  = "enter_coin_name"
	EnterPriceCoin = "enter_price_coin"
	CheckSell      = "check_sell"

	MyOrders = "my_orders"

	SendYourCoins   = "send_your_coins"
	Help            = "help"
	WaitDepositCoin = "wait_deposit_coin"
)

type Message struct {
	chatID      int64
	messageLang string
	localizer   *i18n.Localizer
	reply       string
}

func (message *Message) Reply() string {
	return message.reply
}

func (message *Message) SetReply(reply string) {
	message.reply = reply
}

func (message *Message) SetChatID(chatID int64) {
	message.chatID = chatID
}

func (message *Message) ChatID() int64 {
	return message.chatID
}

func (message *Message) SetLocalizer(localizer *i18n.Localizer) {
	message.localizer = localizer
}

func (message *Message) Localizer() *i18n.Localizer {
	return message.localizer
}

func (message *Message) SetMessageLang(messageLang string) {
	message.messageLang = messageLang
}

func (message *Message) MessageLang() string {
	return message.messageLang
}

func (message *Message) Translate(text string) string {
	return message.Localizer().MustLocalize(&i18n.LocalizeConfig{MessageID: text})
}

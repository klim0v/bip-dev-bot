package callback

import (
	"bip-dev/service/message"
)

type CallbackFactory struct {
	message.Message
	MessageUpdateID int
	Command         string
	Args            string
	Repository      *message.Repository
}

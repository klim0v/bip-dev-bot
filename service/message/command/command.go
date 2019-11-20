package command

import (
	"bip-dev/service/message"
)

type CommandFactory struct {
	message.Message
	Command    string
	Args       string
	Repository *message.Repository
}

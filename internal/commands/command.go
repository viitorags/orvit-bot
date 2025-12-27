package commands

import (
	"context"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var RegisteredCommands = make(map[string]Command)

type Command interface {
	Execute(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error
	Name() string
}

func Register(cmd Command) {
	RegisteredCommands[cmd.Name()] = cmd
}

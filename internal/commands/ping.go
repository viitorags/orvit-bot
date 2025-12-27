package commands

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type PingCommand struct{}

func (p *PingCommand) Name() string        { return "ping" }
func (p *PingCommand) Description() string { return "Responde com Pong!" }

func (p *PingCommand) Execute(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error {
	msg := &waE2E.Message{
		Conversation: proto.String("!pong"),
	}

	_, err := client.SendMessage(ctx, evt.Info.Chat, msg)
	return err
}

func init() {
	Register(&PingCommand{})
}

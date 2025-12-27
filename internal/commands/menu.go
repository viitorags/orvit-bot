package commands

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type MenuCommand struct{}

func (p *MenuCommand) Name() string        { return "menu" }
func (p *MenuCommand) Description() string { return "Mostra o menu de comandos" }

func (p *MenuCommand) Execute(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error {
	if evt.Info.IsGroup {
		msg := &waE2E.Message{
			Conversation: proto.String(`*Menu de Comandos*
> !fig  (Cria figurinhas(Videos, Imagens, Gifs))
> !ping (Faz !pong)
> !info (Traz algumas informações sobre o grupo)
> !menu (Mostra o menu de comandos)
				`),
		}
		_, err := client.SendMessage(ctx, evt.Info.Chat, msg)
		return err
	} else {
		return nil
	}
}

func init() {
	Register(&MenuCommand{})
}

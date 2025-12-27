package commands

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type InfoCommand struct{}

func (p *InfoCommand) Name() string        { return "info" }
func (p *InfoCommand) Description() string { return "Traz as informações do grupo" }

func (p *InfoCommand) Execute(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error {
	if evt.Info.IsGroup {
		infoGroup, _ := client.GetGroupInfo(ctx, evt.Info.Chat)

		msg := &waE2E.Message{
			Conversation: proto.String(fmt.Sprintf(`*Informações do Grupo*
Nome Grupo: %v

Grupo Criado: %v

Descrição: %v

Membros: %v
				`, infoGroup.Name, infoGroup.GroupCreated.Format("02/01/2006"), infoGroup.Topic, infoGroup.ParticipantCount)),
		}

		_, err := client.SendMessage(ctx, evt.Info.Chat, msg)
		return err
	} else {

	}

	return nil
}

func init() {
	Register(&InfoCommand{})
}

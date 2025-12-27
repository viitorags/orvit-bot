package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/viitorags/orvit/internal/commands"
	"github.com/viitorags/orvit/internal/services"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (b *Bot) HandleMessage(evt *events.Message) {
	godotenv.Load()

	botPrefix := os.Getenv("BOT_PREFIX")

	if evt.Info.IsFromMe {
		return
	}

	if time.Since(evt.Info.Timestamp) > 2*time.Minute {
		return
	}

	var text string

	msg := evt.Message
	if msg.GetConversation() != "" {
		text = msg.GetConversation()
	} else if msg.ExtendedTextMessage != nil {
		text = msg.ExtendedTextMessage.GetText()
	} else if msg.ImageMessage != nil {
		text = msg.ImageMessage.GetCaption()
	} else if msg.VideoMessage != nil {
		text = msg.VideoMessage.GetCaption()
	}

	if text == "" {
		return
	}

	hasPrefix, err := regexp.MatchString(fmt.Sprintf(`\%s\b`, botPrefix), strings.ToLower(text))
	if err != nil {
		log.Fatal(err)
	}

	if hasPrefix {
		parts := strings.Split(strings.TrimSpace(text), " ")
		cmdName := strings.ToLower(strings.TrimPrefix(parts[0], "!"))

		fmt.Printf("-> Tentando executar comando: %s\n", cmdName)

		if cmd, ok := commands.RegisteredCommands[cmdName]; ok {
			go func() {
				err := cmd.Execute(context.Background(), b.Client, evt)
				if err != nil {
					fmt.Printf("Erro no comando %s: %v\n", cmdName, err)
				} else {
					fmt.Printf("%s executado com sucesso!\n", cmdName)
				}
			}()
		} else {
			fmt.Printf("Comando '!%s' não encontrado no mapa.\n", cmdName)
			participant := fmt.Sprintf("%s@%s", evt.Info.Sender.User, evt.Info.Sender.Server)

			extMsg := &waE2E.ExtendedTextMessage{
				Text: proto.String(fmt.Sprintf("O comando: !%v não existe", cmdName)),
				ContextInfo: &waE2E.ContextInfo{
					StanzaID:      proto.String(evt.Info.ID),
					Participant:   &participant,
					QuotedMessage: evt.Message,
				},
			}

			msg := &waE2E.Message{
				ExtendedTextMessage: extMsg,
			}

			b.Client.SendMessage(context.Background(), evt.Info.Chat, msg)
		}
	} else {
		if evt.Info.IsGroup {
			containsBotName, _ := regexp.MatchString(`(?i)\borvit\b`, text)
			isMentioned, _ := regexp.MatchString(`(?i)\@255525405093910\b`, text)

			repliedToBot := false
			if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.ContextInfo != nil {
				ctx := evt.Message.ExtendedTextMessage.ContextInfo

				if ctx.Participant != nil {
					quotedParticipant := *ctx.Participant

					myPhone := b.Client.Store.ID.User
					myLID := b.Client.Store.LID.User

					fmt.Printf("[DEBUG] Citado: %s | Meus IDs: %s / %s\n", quotedParticipant, myPhone, myLID)

					if strings.Contains(quotedParticipant, myPhone) || (myLID != "" && strings.Contains(quotedParticipant, myLID)) {
						repliedToBot = true
					}
				}
			}

			fmt.Printf("BotName: %v | Replied: %v | Text: %s | Lid: %v\n", containsBotName, repliedToBot, text, evt.Info.Sender.String())

			if containsBotName || repliedToBot || isMentioned {
				fmt.Println("Processando mensagem...")

				textFormated := text

				for _, jidStr := range evt.Message.GetExtendedTextMessage().GetContextInfo().GetMentionedJID() {
					parts := strings.Split(jidStr, "@")
					idSomenteNumeros := parts[0]
					nomeParaSubstituir := idSomenteNumeros

					if jid, err := types.ParseJID(jidStr); err == nil {
						if contato, errStore := b.Client.Store.Contacts.GetContact(context.Background(), jid); errStore == nil {
							if contato.FullName != "" {
								nomeParaSubstituir = contato.FullName
							} else if contato.PushName != "" {
								nomeParaSubstituir = contato.PushName
							}
						}
					}

					textFormated = strings.ReplaceAll(textFormated, "@"+idSomenteNumeros, nomeParaSubstituir)
				}

				senderName := evt.Info.PushName
				if senderName == "" {
					senderName = "Usuário"
				}

				textFormated = strings.ReplaceAll(textFormated, "@"+evt.Info.Sender.User, senderName)

				fmt.Println("Texto Formatado:", textFormated)

				textMsg := services.HuggingFace(textFormated)

				participantJID := evt.Info.Sender.String()

				extMsg := &waE2E.ExtendedTextMessage{
					Text: proto.String(textMsg),
					ContextInfo: &waE2E.ContextInfo{
						StanzaID:      proto.String(evt.Info.ID),
						Participant:   proto.String(participantJID),
						QuotedMessage: evt.Message,
					},
				}

				msg := &waE2E.Message{ExtendedTextMessage: extMsg}
				_, err := b.Client.SendMessage(context.Background(), evt.Info.Chat, msg)
				if err != nil {
					fmt.Printf("Erro ao enviar: %v\n", err)
				}
			}
		} else {
			if evt.Info.IsFromMe {
				return
			}

			textMsg := services.HuggingFace(text)

			msg := &waE2E.Message{
				Conversation: proto.String(textMsg),
			}

			b.Client.SendMessage(context.Background(), evt.Info.Chat, msg)
			fmt.Println(evt.Info.Sender, " :", textMsg)
		}
	}
}

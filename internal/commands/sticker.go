package commands

import (
	"context"
	"fmt"

	"github.com/viitorags/orvit/internal/helpers"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type StickerCommand struct{}

func (p *StickerCommand) Name() string        { return "fig" }
func (p *StickerCommand) Description() string { return "Cria figurinhas a partir de imagens" }

func (p *StickerCommand) Execute(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error {
	jid := evt.Info.Chat
	extractMedia := func(m *waE2E.Message) (img *waE2E.ImageMessage, vid *waE2E.VideoMessage) {
		if m == nil {
			return
		}

		img = m.GetImageMessage()
		vid = m.GetVideoMessage()

		if img == nil && m.GetViewOnceMessage().GetMessage().GetImageMessage() != nil {
			img = m.GetViewOnceMessage().GetMessage().GetImageMessage()
		}
		if vid == nil && m.GetViewOnceMessage().GetMessage().GetVideoMessage() != nil {
			vid = m.GetViewOnceMessage().GetMessage().GetVideoMessage()
		}

		if img == nil && m.GetDeviceSentMessage().GetMessage().GetImageMessage() != nil {
			img = m.GetDeviceSentMessage().GetMessage().GetImageMessage()
		}
		return
	}

	var imgMsg *waE2E.ImageMessage
	var videoMsg *waE2E.VideoMessage

	imgMsg, videoMsg = extractMedia(evt.Message)

	if imgMsg == nil && videoMsg == nil {
		quotedMsg := evt.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
		if quotedMsg != nil {
			imgMsg, videoMsg = extractMedia(quotedMsg)
		}
	}

	var msgToDownload whatsmeow.DownloadableMessage
	if imgMsg != nil {
		msgToDownload = imgMsg
	} else if videoMsg != nil {
		msgToDownload = videoMsg
		if videoMsg.GetGifPlayback() {
			fmt.Println("O arquivo é um GIF, mas o download é igual ao de vídeo")
		}
	} else {
		return fmt.Errorf("nenhuma mídia encontrada para download")
	}

	data, err := client.Download(ctx, msgToDownload)
	if err != nil {
		return fmt.Errorf("erro ao baixar imagem: %w", err)
	}

	senderName := evt.Info.PushName
	if senderName == "" {
		senderName = "Usuário"
	}

	stickerData, err := helpers.ProcessSticker(data, senderName)
	if err != nil {
		return err
	}

	resp, err := client.Upload(ctx, stickerData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("erro no upload: %w", err)
	}

	participant := fmt.Sprintf("%s@%s", evt.Info.Sender.User, evt.Info.Sender.Server)
	_, err = client.SendMessage(ctx, jid, &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(evt.Info.ID),
				Participant:   &participant,
				QuotedMessage: evt.Message,
			},
			URL:           proto.String(resp.URL),
			DirectPath:    proto.String(resp.DirectPath),
			MediaKey:      resp.MediaKey,
			Mimetype:      proto.String("image/webp"),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(stickerData))),
		},
	})

	return err
}

func init() {
	Register(&StickerCommand{})
}

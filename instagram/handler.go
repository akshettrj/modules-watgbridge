package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"watgbridge/modules"
	"watgbridge/state"
	"watgbridge/utils"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

func InstagramModuleWhatsAppEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.Timestamp.UTC().Before(state.State.StartTime) {
			// Old events
			return
		}

		var chat waTypes.JID
		if v.Info.IsIncomingBroadcast() && !v.Info.IsGroup {
			chat = v.Info.MessageSource.Sender
		} else {
			chat = v.Info.Chat
		}

		if !v.Info.IsFromMe && !slices.Contains(instaConfig.WhatsAppAllowedGroups, chat.User) {
			return
		}

		text := ""
		if extendedMessageText := v.Message.GetExtendedTextMessage().GetText(); extendedMessageText != "" {
			text = extendedMessageText
		} else {
			text = v.Message.GetConversation()
		}

		if text == "" {
			return
		}

		textSplit := strings.Split(text, " \n\t")
		for _, token := range textSplit {
			if IsSupportedLink(token) && !IsStoriesLink(token) {
				downloadLink(token, v, chat)
			} else if IsInstagramLink(token) {
				tryUserProfile(token, v, chat)
			}
		}
	}
}

func downloadLink(link string, v *events.Message, chat waTypes.JID) {
	waClient := state.State.WhatsAppClient

	req, _ := http.NewRequest("GET", link, nil)

	AddCookies(req)
	AddHeaders(req)
	AddQueries(req)

	res, err := client.Do(req)
	if err != nil {
		utils.WaSendText(chat, fmt.Sprintf("Could not get JSON data:\n\n%s",
			err.Error()), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
			v.Message, true)
		return
	}
	defer res.Body.Close()
	SaveCookies(res)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.WaSendText(chat, fmt.Sprintf("Could not read response body:\n\n%s",
			err.Error()), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
			v.Message, true)
		return
	}

	mediaType := GetMediaType(body)

	switch mediaType {

	case MediaTypeCarousel:
		var ic InstagramCarousel
		err := json.Unmarshal(body, &ic)
		if err != nil {
			utils.WaSendText(chat, fmt.Sprintf("Could not parse body into InstagramCarousel:\n\n%s",
				err.Error()), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
				v.Message, true)
			return
		}

		var (
			caption           = ic.Caption()
			successfulUploads = 0
		)

		for _, item := range ic.Items[0].CarouselMedia {
			mediaLink := item.DownloadLink()

			req, _ := http.NewRequest("GET", mediaLink, nil)
			itemBytes, err := DownloadFile(req)
			if err != nil {
				continue
			}

			if item.MediaType == MediaTypeImage {
				uploadedImage, err := waClient.Upload(context.Background(), itemBytes, whatsmeow.MediaImage)
				if err != nil {
					continue
				}

				msgToSend := &waProto.Message{
					ImageMessage: &waProto.ImageMessage{
						Url:               proto.String(uploadedImage.URL),
						DirectPath:        proto.String(uploadedImage.DirectPath),
						MediaKey:          uploadedImage.MediaKey,
						MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
						Mimetype:          proto.String(http.DetectContentType(itemBytes)),
						FileEncSha256:     uploadedImage.FileEncSHA256,
						FileSha256:        uploadedImage.FileSHA256,
						FileLength:        proto.Uint64(uint64(len(itemBytes))),
						Height:            proto.Uint32(uint32(item.Height)),
						Width:             proto.Uint32(uint32(item.Width)),
						ContextInfo: &waProto.ContextInfo{
							StanzaId:      proto.String(v.Info.ID),
							Participant:   proto.String(v.Info.MessageSource.Sender.ToNonAD().String()),
							QuotedMessage: v.Message,
						},
					},
				}

				_, err = waClient.SendMessage(context.Background(), chat, msgToSend)
				if err == nil {
					successfulUploads += 1
				}
			} else if mediaType == MediaTypeVideo {
				uploadedVideo, err := waClient.Upload(context.Background(), itemBytes, whatsmeow.MediaVideo)
				if err != nil {
					continue
				}

				msgToSend := &waProto.Message{
					VideoMessage: &waProto.VideoMessage{
						Url:           proto.String(uploadedVideo.URL),
						DirectPath:    proto.String(uploadedVideo.DirectPath),
						MediaKey:      uploadedVideo.MediaKey,
						Mimetype:      proto.String(http.DetectContentType(itemBytes)),
						FileEncSha256: uploadedVideo.FileEncSHA256,
						FileSha256:    uploadedVideo.FileSHA256,
						FileLength:    proto.Uint64(uint64(len(itemBytes))),
						Seconds:       proto.Uint32(uint32(item.VideoDuration)),
						GifPlayback:   proto.Bool(false),
						Height:        proto.Uint32(uint32(item.Height)),
						Width:         proto.Uint32(uint32(item.Width)),
						ContextInfo: &waProto.ContextInfo{
							StanzaId:      proto.String(v.Info.ID),
							Participant:   proto.String(v.Info.MessageSource.Sender.ToNonAD().String()),
							QuotedMessage: v.Message,
						},
					},
				}

				_, err = waClient.SendMessage(context.Background(), chat, msgToSend)
				if err == nil {
					successfulUploads += 1
				}
			}
		}

		if successfulUploads > 0 {
			utils.WaSendText(chat, caption, v.Info.ID, v.Info.MessageSource.Sender.
				ToNonAD().String(), v.Message, true)
			return
		}

	case MediaTypeVideo:
		var ir InstagramReel
		err := json.Unmarshal(body, &ir)
		if err != nil {
			utils.WaSendText(chat, fmt.Sprintf("Could not parse body into InstagramReel:\n\n%s",
				err.Error()), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
				v.Message, true)
			return
		}

		var (
			caption   = ir.Caption()
			mediaLink = ir.DownloadLink()
		)

		req, _ := http.NewRequest("GET", mediaLink, nil)
		videoBytes, err := DownloadFile(req)
		if err != nil {
			return
		}

		uploadedVideo, err := waClient.Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)
		if err != nil {
			return
		}

		msgToSend := &waProto.Message{
			VideoMessage: &waProto.VideoMessage{
				Caption:       proto.String(caption),
				Url:           proto.String(uploadedVideo.URL),
				DirectPath:    proto.String(uploadedVideo.DirectPath),
				MediaKey:      uploadedVideo.MediaKey,
				Mimetype:      proto.String(http.DetectContentType(videoBytes)),
				FileEncSha256: uploadedVideo.FileEncSHA256,
				FileSha256:    uploadedVideo.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(videoBytes))),
				Seconds:       proto.Uint32(uint32(ir.Items[0].VideoDuration)),
				GifPlayback:   proto.Bool(false),
				Height:        proto.Uint32(uint32(ir.Items[0].Height)),
				Width:         proto.Uint32(uint32(ir.Items[0].Width)),
				ContextInfo: &waProto.ContextInfo{
					StanzaId:      proto.String(v.Info.ID),
					Participant:   proto.String(v.Info.MessageSource.Sender.ToNonAD().String()),
					QuotedMessage: v.Message,
				},
			},
		}

		waClient.SendMessage(context.Background(), chat, msgToSend)

	case MediaTypeImage:
		var ii InstagramImage
		err := json.Unmarshal(body, &ii)
		if err != nil {
			utils.WaSendText(chat, fmt.Sprintf("Could not parse body into InstagramImage:\n\n%s",
				err.Error()), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
				v.Message, true)
			return
		}

		var (
			caption   = ii.Caption()
			mediaLink = ii.DownloadLink()
		)

		req, _ := http.NewRequest("GET", mediaLink, nil)
		imageBytes, err := DownloadFile(req)
		if err != nil {
			return
		}

		uploadedImage, err := waClient.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
		if err != nil {
			return
		}

		msgToSend := &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				Caption:           proto.String(caption),
				Url:               proto.String(uploadedImage.URL),
				DirectPath:        proto.String(uploadedImage.DirectPath),
				MediaKey:          uploadedImage.MediaKey,
				MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
				Mimetype:          proto.String(http.DetectContentType(imageBytes)),
				FileEncSha256:     uploadedImage.FileEncSHA256,
				FileSha256:        uploadedImage.FileSHA256,
				FileLength:        proto.Uint64(uint64(len(imageBytes))),
				Height:            proto.Uint32(uint32(ii.Items[0].Height)),
				Width:             proto.Uint32(uint32(ii.Items[0].Width)),
				ContextInfo: &waProto.ContextInfo{
					StanzaId:      proto.String(v.Info.ID),
					Participant:   proto.String(v.Info.MessageSource.Sender.ToNonAD().String()),
					QuotedMessage: v.Message,
				},
			},
		}

		waClient.SendMessage(context.Background(), chat, msgToSend)

	default:
		utils.WaSendText(chat, fmt.Sprintf("Unkown media type:\n\n[%v]",
			mediaType), v.Info.ID, v.Info.MessageSource.Sender.ToNonAD().String(),
			v.Message, true)
		return

	}
}

func tryUserProfile(link string, v *events.Message, chat waTypes.JID) {
	waClient := state.State.WhatsAppClient

	req, _ := http.NewRequest("GET", link, nil)
	AddCookies(req)
	AddHeaders(req)
	AddQueries(req)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	SaveCookies(res)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var iup InstagramUserProfile
	err = json.Unmarshal(body, &iup)
	if err != nil {
		return
	}

	if iup.Graphql.User.ID == "" {
		return
	}

	dpLink := iup.ProfilePicURLHD()
	req, _ = http.NewRequest("GET", dpLink, nil)
	dpBytes, err := DownloadFile(req)
	if err != nil {
		return
	}

	uploadedImage, err := waClient.Upload(context.Background(), dpBytes, whatsmeow.MediaImage)
	if err != nil {
		return
	}

	msgToSend := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:           proto.String(iup.Caption()),
			Url:               proto.String(uploadedImage.URL),
			DirectPath:        proto.String(uploadedImage.DirectPath),
			MediaKey:          uploadedImage.MediaKey,
			MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
			Mimetype:          proto.String(http.DetectContentType(dpBytes)),
			FileEncSha256:     uploadedImage.FileEncSHA256,
			FileSha256:        uploadedImage.FileSHA256,
			FileLength:        proto.Uint64(uint64(len(dpBytes))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      proto.String(v.Info.ID),
				Participant:   proto.String(v.Info.MessageSource.Sender.ToNonAD().String()),
				QuotedMessage: v.Message,
			},
		},
	}

	waClient.SendMessage(context.Background(), chat, msgToSend)
}

func init() {
	modules.WhatsAppHandlers = append(modules.WhatsAppHandlers,
		InstagramModuleWhatsAppEventHandler)
}

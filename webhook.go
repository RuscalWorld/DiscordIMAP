package main

import (
	"crypto/md5"
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/gtuk/discordwebhook"
)

func PublishMessage(message *imap.Message, body string) error {
	authorName := message.Envelope.From[0].Address()
	authorAvatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(authorName)))

	discordMessage := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{
			{
				Title:       &message.Envelope.Subject,
				Description: &body,
				Author: &discordwebhook.Author{
					Name:    &authorName,
					IconUrl: &authorAvatarURL,
				},
			},
		},
	}

	return discordwebhook.SendMessage(discordWebhookURL, discordMessage)
}

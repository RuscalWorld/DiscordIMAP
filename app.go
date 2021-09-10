package main

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
)

var (
	imapHost          = requireEnv("IMAP_HOST")
	imapUser          = requireEnv("IMAP_USER")
	imapPassword      = requireEnv("IMAP_PASSWORD")
	imapMailbox       = getEnvOrDefault("IMAP_MAILBOX", "INBOX")
	imapTLS           = getEnvOrDefault("IMAP_TLS", "true")
	discordWebhookURL = requireEnv("DISCORD_WEBHOOK_URL")
)

func getEnvOrDefault(key, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	} else {
		return def
	}
}

func requireEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	panic(fmt.Sprintf("%s environment variable must be set", key))
}

func main() {
	var imapClient *client.Client
	var err error

	if imapTLS == "true" || imapTLS == "1" {
		imapClient, err = client.DialTLS(imapHost, nil)
	} else {
		imapClient, err = client.Dial(imapHost)
	}

	if err != nil {
		log.Fatalln("Unable to connect:", err)
		return
	}

	err = imapClient.Login(imapUser, imapPassword)
	if err != nil {
		log.Fatalln("Unable to login:", err)
		return
	}

	HandleMailbox(imapClient)
}

package main

import (
	"bytes"
	"github.com/mattn/godown"
	"log"
	"net/mail"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func HandleMailbox(client *client.Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		_ = <-ticker.C

		_, err := client.Select(imapMailbox, false)
		if err != nil {
			log.Println("Unable to fetch mailbox status:", err)
			continue
		}

		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}

		nums, err := client.Search(criteria)
		if err != nil {
			log.Println("Unable to find messages:", err)
			continue
		}

		if len(nums) == 0 {
			continue
		}

		set := new(imap.SeqSet)
		set.AddNum(nums...)

		section := &imap.BodySectionName{}
		messages := make(chan *imap.Message, 10)
		err = client.Fetch(set, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, imap.FetchBodyStructure, section.FetchItem()}, messages)
		if err != nil {
			log.Println("Unable to fetch messages:", err)
			continue
		}

		log.Println(len(messages))
		for message := range messages {
			reader, err := mail.ReadMessage(message.GetBody(section))
			if err != nil {
				log.Println("Unable to read message:", err)
				continue
			}

			buffer := &bytes.Buffer{}
			err = godown.Convert(buffer, reader.Body, &godown.Option{})
			if err != nil {
				log.Println("Unable to convert message body:", err)
				continue
			}

			_ = PublishMessage(message, string(buffer.Bytes()))
		}
	}
}

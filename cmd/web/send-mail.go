package main

import (
	"github.com/psanodiya94/gobooking.com/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
	} else {
		email := mail.NewMSG()
		email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

		if m.Template == "" {
			email.SetBody(mail.TextPlain, m.Content)
		} else {
			data, err := os.ReadFile(
				filepath.Join("./templates/emails", m.Template),
			)
			if err != nil {
				app.ErrorLog.Println(err)
			} else {
				mailTemplate := string(data)
				msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
				email.SetBody(mail.TextHTML, msgToSend)
			}
		}

		err = email.Send(client)
		if err != nil {
			app.ErrorLog.Println(err)
		} else {
			app.InfoLog.Println("Email sent!")
		}
	}
}

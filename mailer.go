package main

import (
	"net/smtp"
)

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func SendMail(content *string) error {
	from := "From: " + config.MailFromText + " <" + config.MailFrom + ">"
	body := from + "\r\nSubject: " + config.MailSubject + "\r\n" + MIME + "\r\n" + *content

	err := smtp.SendMail(config.MailHost+":"+config.MailPort,
		smtp.PlainAuth("", config.MailFrom, config.MailPass, config.MailHost),
		config.MailFrom, config.MailTo, []byte(body))

	return err
}

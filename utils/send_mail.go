package utils

import (
	"fmt"
	"net/smtp"
)

const (
	Server  = "smtp.masterplac.com"
	Sender  = "info@masterplac.com"
	Subject = "Subject: Mensaje de PAPYRUS\r\n\r\n"
	Passwd  = "Mail.2012+"
	Port    = 465
)

func SendMail(to, msg string) error {
	auth := smtp.PlainAuth("", Sender, Passwd, Server)

	sliceTo := []string{to}
	byteMsg := []byte(fmt.Sprintf("%v%v", Subject, msg))

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", Server, 465),
		auth,
		Sender,
		sliceTo,
		byteMsg,
	)
}

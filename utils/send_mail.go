package utils

import (
	"errors"
	"fmt"
	"net/smtp"
)

const ()

var (
	Enabled = false
	Server  string
	Sender  string
	Subject string
	Passwd  string
	Port    int
)

func SendMail(to, msg string) error {
	if !Enabled {
		return errors.New("Email sending not enabled in config.json file")
	}
	logger := NewAggregatedLogger("EMAIL", "EMAIL")
	auth := smtp.PlainAuth("", Sender, Passwd, Server)

	sliceTo := []string{to}
	byteMsg := []byte(fmt.Sprintf("%v%v", Subject, msg))

	logger.Inf("Sending email...")
	err := smtp.SendMail(
		fmt.Sprintf("%v:%v", Server, Port),
		auth,
		Sender,
		sliceTo,
		byteMsg,
	)
	if err != nil {
		logger.Err("ERROR:", err)
		return err
	}

	logger.Inf("Email sent to", to)
	return err
}

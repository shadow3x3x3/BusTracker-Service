package main

import (
	"encoding/json"
	"net/smtp"
	"os"
)

const (
	mailConfigFile = "mail_config.json"
)

type mailConfig struct {
	From       string `json:"from"`
	Pass       string `json:"password"`
	SMTPServer string `json:"smtp_server"`
}

// MailSender can be used to send mail.
type MailSender struct {
	auth smtp.Auth
}

// Init method can initialize MailSender from config file.
func (m *MailSender) Init() (err error) {
	return initMailConfig(m)
}

// Send method can send email by context
func (m *MailSender) Send(context string, to string) (err error) {
	mailContent := makeMailStructure(
		to,
		"subject",
		context,
	)

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		m.auth,
		"shadow3x3x3@gmail.com",
		[]string{"shadow3x3x3@gmail.com"},
		[]byte(mailContent),
	)

	return nil
}

func initMailConfig(m *MailSender) error {
	config := mailConfig{}

	if err := readMailConfig(&config); err != nil {
		return err
	}

	m.auth = smtp.PlainAuth(
		"",
		config.From,
		config.Pass,
		config.SMTPServer,
	)

	return nil
}

func readMailConfig(c *mailConfig) error {
	file, err := os.Open(mailConfigFile)
	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(&c)
}

func makeMailStructure(
	to string,
	subject string,
	context string) string {

	return "From: " + "shadow3x3x3@gmail.com" + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		context
}

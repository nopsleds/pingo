package impl

import "net/smtp"

type Instance struct {
	SmtpAuth smtp.Auth
}

func ParseSmtpAuth(cfg ConfigSmtp) smtp.Auth {
	return smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
}

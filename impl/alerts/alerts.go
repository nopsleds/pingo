package alerts

import "net/smtp"
improt "../config"

type Instance struct {
	SmtpAuth smtp.Auth
}

func (cfg config.Smtp) smtp.Auth {
	return smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
}




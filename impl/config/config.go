package config

type Config struct {
	Web     WebConfig
	Targets map[string]Target
	Alerts  Alerts
}

type WebConfig struct {
	Port int
}

type Target struct {
	Type        string
	Hostname    string
	Port        int
	IntervalSec int

	HttpUrl            string
	HttpExpectedStatus int
}

type Alerts struct {
	Emails   []string
	Webhooks map[string]Webhook
}

type Webhook struct {
	Method       string
	UrlTemplate  string
	BodyTemplate string
}

const TypeHttp = "http"
const TypeHttps = "https"
const TypeTcp = "tcp"
const TypeTls = "tls"

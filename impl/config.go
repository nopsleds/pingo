package impl

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Web     ConfigWeb
	Smtp    ConfigSmtp
	Targets map[string]ConfigTarget
	Alerts  Alerts
}

type ConfigWeb struct {
	Port int
}
type ConfigSmtp struct {
	Username string
	Password string
	Host     string
}
type ConfigTarget struct {
	Type            string
	Hostname        string
	Port            int
	PollingInterval string

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

func LoadOrInitFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file not found > creating with default value")
			f, err = initDefaultConfigFile(path)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("error while openig configfile: %+v", err)
			return nil, err
		}
	}

	log.Println("file not found > creating with default value")

	var cfg Config
	_, err = toml.DecodeReader(f, &cfg)
	return &cfg, err
}

func initDefaultConfigFile(path string) (*os.File, error) {
	err := ioutil.WriteFile(path, []byte(DEFAULT_CONFIG_FILE_CONTENT), 0666)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

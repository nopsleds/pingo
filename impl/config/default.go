package config

const DEFAULT_CONFIG_FILE_CONTENT = `
[web]
port = 8081

[targets.google-fr]
type = "http"
pollingInterval = "5s"
httpURL = "http://google.fr"

[targets.bad]
type = "http"
pollingInterval = "1s"
httpURL = "http://goosdfsdfefsegle.fr"

[alerts]
emails = [ "test@email.com" ]

[alerts.webhooks.charly]
method = "GET"
urlTemplate = "http://test.com/something"
`

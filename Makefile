install:
	go get github.com/BurntSushi/toml
build:
	go build -o cmd/pingo/pingo cmd/pingo/main.go
run:
	./cmd/pingo/pingo -f res/config.toml
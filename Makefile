.PHONY:
.SILENT:

build:
	go build -o ./bin/bot cmd/bot/main.go

run: build
	./bin/bot -t ${BOT_TOKEN}
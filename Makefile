TAG=default
TOKEN=

debug: ./debug/bot.go
	go run ./debug/bot.go -token=$(TOKEN)

init:
	go mod tidy

publish:
	./scripts/publish.sh $(TAG)

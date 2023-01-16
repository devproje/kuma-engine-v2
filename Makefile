TAG=default
TOKEN=

debug: ./debug/bot.go
	go run ./debug/bot.go -token=$(TOKEN)

publish:
	./scripts/publish.sh $(TAG)

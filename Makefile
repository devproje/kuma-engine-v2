TAG=default

init:
	go mod tidy

publish:
	./publish.sh $(TAG)
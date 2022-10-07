TAG = "default"

init:
	go mod tidy

publish:
	git tag $(TAG)
	git push origin $(TAG)
	GOPROXY=proxy.golang.org go list -m github.com/devproje/kuma-engine@$(TAG)
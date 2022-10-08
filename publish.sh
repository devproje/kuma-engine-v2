#!/bin/sh
MODULE=github.com/devproje/kuma-engine@$(TAG)
PROXY_URL=proxy.golang.org

if [ "$1" == "default" ]; then
	echo "Please type some tag"
else
	git tag $(TAG)
	git push origin $(TAG)
	GOPROXY=$PROXY_URL go list -m $MODULE
fi

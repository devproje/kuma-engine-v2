#!/bin/sh
MODULE=github.com/devproje/kuma-engine@$1
PROXY_URL=proxy.golang.org

if [ "$1" == "default" ]; then
	echo "Please type some tag"
else
	git tag $1
	git push origin $1
	GOPROXY=$PROXY_URL go list -m $MODULE
fi

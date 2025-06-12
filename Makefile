.PHONY: build clean dev

build:
	-mkdir target
	go build -o target/senanomusic
	cp .env.production target/.env

clean:
	-rm -rf target
	-rm __debug*

dev:
	-mkdir tmp
	go run main.go -env=.env.dev

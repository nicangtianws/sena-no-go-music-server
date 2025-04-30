.PHONY: build clean dev

build:
	go build -o bin/senanomusic
	cp .env.production bin/.env

clean:
	rm -rf bin/
	rm __debug*
	rm Music.db

dev:
	go run main.go -env=.env.dev

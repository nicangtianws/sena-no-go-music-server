.PHONY: build clean dev package clean-pkg

build:
	-mkdir -p target/bin
	go build -o target/bin/senanomusic
	cp .env.production target/bin/.env

package:
	$(MAKE) build
	-mkdir target/senanomusic-server
	cp -r target/bin/ target/senanomusic-server/
	cp -r script/ target/senanomusic-server/
	tar -zcf target/senanomusic-server.tar.gz -C target/ senanomusic-server/
	rm -rf target/senanomusic-server/

clean:
	-rm -rf target
	-rm __debug*

clean-pkg:
	-rm -rf target/senanomusic-server
	-rm target/senanomusic-server.tar.gz

dev:
	-mkdir tmp
	go run main.go -env=.env.dev

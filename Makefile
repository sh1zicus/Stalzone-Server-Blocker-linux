BINARY  := stalzone-blocker
MODULE  := github.com/sh1zicus/stalzone-server-blocker
CMD     := ./cmd/stalzone-blocker

.PHONY: build setcap run clean

build:
	go build -o $(BINARY) $(CMD)

setcap: build
	sudo setcap cap_net_admin+ep ./$(BINARY)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

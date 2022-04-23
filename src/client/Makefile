all:
	$(MAKE) -C proto
	go build -o client main.go
	chmod +x client
	cp client ../../test/client/client

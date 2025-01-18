all:
	$(MAKE) -C proto
	go build -o client main.go
	chmod +x client
	mv client ../../script/client/client

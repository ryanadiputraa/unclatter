server:
	air
env:
	cp config/config.example.yml config/config.yml
test:
	go test -v ./...

.PHONY: test

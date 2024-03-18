server:
	air
env:
	cp config/config.example.yml config/config.yml
compose-up:
	docker compose up -d
compose-down:
	docker compose down
test:
	go test -v ./...

.PHONY: test

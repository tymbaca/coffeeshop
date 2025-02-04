run:
	docker compose up -d
	cd waiter && go run ./cmd/guiwaiter &
	cd barista && go run .

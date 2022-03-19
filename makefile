.PHONY: run
run:
	go run ./cmd/chesspal/main.go

frontend:
	cd web/vue-frontend && npm run serve -- --port 3000
.PHONY: run
run:
	go run ./cmd/chesspal/main.go

frontend:
	cd web/vue-frontend && npm run serve -- --port 3000

build-frontend:
	cd web/vue-frontend && npm run build
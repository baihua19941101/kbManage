.PHONY: backend-test frontend-test test lint fmt

backend-test:
	cd backend && go test ./...

frontend-test:
	cd frontend && npm test

test: backend-test frontend-test

lint:
	cd frontend && npm run lint

fmt:
	cd backend && gofmt -w ./cmd ./internal

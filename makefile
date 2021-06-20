run:
	go run ./cmd/app

test:
	go test -bench . -benchmem -benchtime 3s

tidy:
	go mod tidy
	go mod vendor
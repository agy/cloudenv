bin/cloudenv: cmd/cloudenv/main.go vet test fmt
	go build -o $@ cmd/cloudenv/main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -race -cover

.PHONY: fmt
fmt:
	go fmt ./...

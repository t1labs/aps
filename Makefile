test:
	go test -cover ./...

security:
	go get github.com/securego/gosec/cmd/gosec/...
	gosec -exclude=G104 ./...

vendor:
	GO111MODULE=on go mod vendor

vet:
	go vet ./...

.PHONY: test security vendor vet
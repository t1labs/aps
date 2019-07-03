test:
	go test -cover ./...

vendor:
	GO111MODULE=on go mod vendor

.PHONY: vendor
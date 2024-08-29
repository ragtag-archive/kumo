bin/kumo: $(shell find . -name '*.go')
	mkdir -p bin/
	CGO_ENABLED=0 go build \
		-ldflags "-extldflags '-static' -s -w" \
		-tags netgo \
		-o bin/kumo ./main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: docker
docker:
	docker build -t kumo .

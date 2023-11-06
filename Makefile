# docker
.PHONY: run/docker
run/docker:
	docker-compose up -d

# run consumer and producer
.PHONY: run/producer
run/producer:
	go run cmd/producer/producer.go

# run consumer and producer
.PHONY: run/consumer
run/consumer:
	go run cmd/consumer/consumer.go

####### utils #####
## fmt: format code
.PHONY: fmt
fmt:
	go fmt ./...

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

# build: build all files
.PHONY: build
build:
	go build ./...

# clean:
.PHONY: clean
clean:
	go clean -modcache

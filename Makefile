BIN=./bin/prservice
PRSERVICE_MAIN=./cmd/prservice/main.go
GOSOURCES=$(shell find . -name '*.go')
GOOS=linux
CGO_ENABLED=0

.PHONY: all clean run

all: ${BIN} 

${BIN}: ${PRSERVICE_MAIN} ${GOSOURCES}
	go mod tidy 
	go mod download 
	go test -v ./...
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} go build -o $@ $< 

run: ${BIN}
	$<

clean:
	rm -rf ./bin

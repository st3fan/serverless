.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/submit cmd/submit/submit.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/execute cmd/execute/execute.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/api cmd/api/api.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

.PHONY: build clean deploy gomodgen

build: clean
	env GOOS=linux go build -ldflags="-s -w" -o .bin/hello hello/main.go

clean:
	rm -rf .bin .artifact.zip

artifact: clean
	zip -r .artifact.zip .

deploy: build
	sls deploy --verbose

.PHONY: build-app run clean

build-app:
	go build -o bin/flexcreek ./cmd/flexcreek.go

run: build-app
	@./bin/flexcreek

clean:
	@rm -rf bin

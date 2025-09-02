build:
	go build -v ./cmd/snippetbox

.PHONEY: build

run:
	./snippetbox -addr=":8080"

.PHONEY: run

.DEFAULT_GOAL := build
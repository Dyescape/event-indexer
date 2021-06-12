build: compile

compile:
	@echo Compiling binary
	@go build .

run: build
	./event-indexer run
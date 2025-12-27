BUF := go run github.com/bufbuild/buf/cmd/buf@latest

generate:
	$(BUF) generate
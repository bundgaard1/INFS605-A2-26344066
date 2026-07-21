.PHONY: tools generate dev

# Tjekker om buf er installeret, ellers installeres den automatisk
tools:
	@which buf > /dev/null || (echo "Installerer Buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)

# Kører tool-tjek og genererer derefter kode
generate: tools
	@echo "Genererer Protobuf & gRPC kode..."
	@cd proto && buf generate

dev: generate
	@cd frontend && go run cmd/main.go


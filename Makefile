.PHONY: tools generate dev

# Checks if buf is installed, otherwise it is installed automatically
tools:
	@which buf > /dev/null || (echo "Installing Buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)
	@which templ > /dev/null || (echo "Installing templ..." && go install github.com/a-h/templ/cmd/templ@latest)

# Runs the tool check and then generates code
generate: tools
	@echo "Generating Protobuf & gRPC code..."
	@cd proto && buf generate
	@echo "Generating templ components..."
	@cd frontend && templ generate

dev: generate
	@cd frontend && go run cmd/main.go

clean:
	@echo "Cleaning up generated files..."
	@rm -rf **/gen
	@find frontend -name '*_templ.go' -delete

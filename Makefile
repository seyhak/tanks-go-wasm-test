@echo "wasm -> compile to wasm"

wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm

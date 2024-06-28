wasm:
	env GOOS=js GOARCH=wasm go build -o build/town.wasm github.com/mikelangelon/town-sweet-town

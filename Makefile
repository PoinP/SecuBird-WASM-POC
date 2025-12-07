build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/main.wasm

build:
	GOOS=js GOARCH=wasm go build -o bin/main.wasm
	
	cp $(shell go env GOROOT)/lib/wasm/wasm_exec.js bin/wasm_exec.js
	cp $(shell go env GOROOT)/lib/wasm/wasm_exec_node.js bin/wasm_exec_node.js
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.html bin/index.html

	sed -i "" 's/test.wasm/main.wasm/g' bin/index.html
	sed -i "" 's#../../lib/wasm/wasm_exec.js#./wasm_exec.js#g' bin/index.html

serve:
	go run ./server/main.go ./bin

run:
	go run main.go

# WebAssembly

Compile and run a simple hello world game using WebAssembly.

Ref: https://ebitengine.org/en/documents/webassembly.html

## Compile

```bash
env GOOS=js GOARCH=wasm go build -o helloworld.wasm github.com/gvicentin/webassembly
```

This will generate a `helloworld.wasm` file.

## Copy the wasm_exec.js file

The `wasm_exec.js` file is needed to run the WebAssembly code. It is located in the Go installation directory.

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

## Run

```bash
python3 -m http.server
```

Open your browser and go to `http://localhost:8000/`. Open the console to see the output.

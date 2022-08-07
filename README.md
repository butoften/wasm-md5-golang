# wasm
wasm孵化器

### macos
* go本身生成wasm
```go
GOOS=js GOARCH=wasm go build -o xxx.wasm
GOOS=js GOARCH=wasm go build -o src/wasm/md5.wasm
```
* tinygo生成wasm
```go
brew tap tinygo-org/tools
brew install tinygo

tinygo version
tinygo build -o xxx-tiny.wasm
```

* 原生编译
```go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

* tinygo编译
```go
cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" ./wasm_exec_tiny.js
```



### windows
临时设置 golang 环境参数（仅作用于当前CMD）
* set GOOS=js 
* set GOARCH=wasm
* go build -o main.wasm

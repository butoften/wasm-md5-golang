# wasm
wasm孵化器

### 启动前端静态服务器
```
yarn
yarn server
```
如果遇到报错：env: node\r: No such file or directory 
说明是文件回车换行问题，我的mac系统解决方案是，用你的编辑器比如vscode 打开node_module里live-server文件夹
把index.js与live-server.js的 crlf改成lf
再执行yarn server

### macos
* go本身生成wasm
```go
GOOS=js GOARCH=wasm go build -o xxx.wasm
GOOS=js GOARCH=wasm go build -o ../src/wasm/md5.wasm

gzip --best md5.wasm
```
* tinygo生成wasm
```go
brew tap tinygo-org/tools
brew install tinygo

tinygo version
tinygo build -o xxx-tiny.wasm
tinygo build -o ../src/wasm/md5-tiny.wasm
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
临时设置 golang 环境参数 (本人没测试，window系统建议使用git bash来操作，git bash会模拟linux环境)
* set GOOS=js 
* set GOARCH=wasm
* go build -o main.wasm

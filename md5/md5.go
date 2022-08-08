package main

import (
	"crypto/md5"
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	c := make(chan struct{})
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	// js.Global().Set("sayHello", js.FuncOf(sayHello))
	/* js.Global().Get("document").
	Call("getElementById", "brightness").
	Call("addEventListener", "change", brightnessCb) */
	js.Global().Set("getFileMd5", js.FuncOf(getFileMd5))
	js.Global().Set("wasmMd5Add", js.FuncOf(wasmMd5Add))
	js.Global().Set("wasmMd5End", js.FuncOf(wasmMd5End))
}

func getFileMd5(value js.Value, args []js.Value) interface{} {
	startTime := time.Now()
	array := args[0]
	byteLength := array.Get("byteLength").Int()
	var buffer []uint8 = make([]uint8, byteLength)
	js.CopyBytesToGo(buffer, array)
	md5hashByteArr := md5.Sum(buffer)
	md5hash := fmt.Sprintf("%x", md5hashByteArr)
	// endTime := time.Now()
	// elapsed := endTime.Sub(startTime).Seconds()
	elapsed := time.Since(startTime)
	fmt.Printf("elapsed: %v\n", elapsed)
	return md5hash
}

var md5hash = md5.New()

func wasmMd5Add(value js.Value, args []js.Value) interface{} {
	array := args[0]
	byteLength := array.Get("byteLength").Int()
	var buffer []uint8 = make([]uint8, byteLength)
	js.CopyBytesToGo(buffer, array)
	md5hash.Write(buffer)
	return nil
}
func wasmMd5End(value js.Value, args []js.Value) interface{} {
	finalMd5Hash := fmt.Sprintf("%x", md5hash.Sum(nil))
	return finalMd5Hash
}

package main

import (
	"crypto/md5"
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	// js.Global().Set("sayHello", js.FuncOf(sayHello))
	/* js.Global().Get("document").
	Call("getElementById", "brightness").
	Call("addEventListener", "change", brightnessCb) */
	js.Global().Set("getFileMd5", js.FuncOf(getFileMd5))
	js.Global().Set("initWasmMd5", js.FuncOf(initWasmMd5))
	js.Global().Set("wasmMd5Add", js.FuncOf(wasmMd5Add))
	js.Global().Set("wasmMd5End", js.FuncOf(wasmMd5End))
}

func sayHello(value js.Value, args []js.Value) interface{} {

	/* fmt.Printf("args: %v\n", len(args))
	a := args[0]
	b := args[1]
	*/
	// h := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	// h.Len *= 4
	// h.Cap *= 4
	// fmt.Printf("(*(*[]byte)(unsafe.Pointer(h))): %v\n", *(*[]byte)(unsafe.Pointer(h)))

	/* fmt.Printf("a: %v-%T\n", a, a)
	fmt.Printf("b: %v-%T\n", b, a) */
	return []any{1, 1}
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

/* var totalSize int //总字节数
var sliceSize int //分段字节数
var count int = 0; */

func initWasmMd5(value js.Value, args []js.Value) interface{} {
	/* totalSize = args[0].Int()
	sliceSize = args[1].Int() */
	// fmt.Printf("js传过来的 totalSize: %v\n", totalSize)
	// fmt.Printf("js传过来的 sliceSize: %v\n", sliceSize)
	return nil
}

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

var brightnessCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	delta := args[0].Get("target").Get("valueAsNumber").Float()
	start := time.Now()
	fmt.Printf("start: %v\n", start)
	fmt.Printf("delta: %v\n", delta)
	return nil
})
var onInputChange = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// quick return if no source image is yet uploaded
	delta := args[0].Get("target").Get("valueAsNumber").Float()
	start := time.Now()
	fmt.Printf("start: %v\n", start)
	fmt.Printf("delta: %v\n", delta)
	/* s.updateImage(res, start)
	args[0].Call("preventDefault") */
	return nil
})

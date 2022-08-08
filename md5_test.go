package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var path = flag.String("path", "/Users/wuhao/Downloads/测试视频/600M.mp4", "")

//4b26c974ebb21f5bdfe3ae04a8d322b3
// var path = flag.String("path", "/Users/wuhao/Downloads/test.jpeg", "")

func TestMd5(t *testing.T) {
	aaa() //15.14s 8.744s 6.161s 7.683s 9.903s 5.938s
	//bbb() //3.136s 3.138s 2.999s 3.125s 4.582s 3.398s
}

func aaa() {
	f, err := os.Open(*path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ReadAll", err)
		return
	}

	// md5.Sum(body)
	fmt.Printf("%x\n", md5.Sum(body))
}

func bbb() {
	f, err := os.Open(*path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}

	defer f.Close()

	md5hash := md5.New()
	fmt.Printf("f: %v\n", f)
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return
	}

	md5hash.Sum(nil)
	fmt.Printf("%x\n", md5hash.Sum(nil))
}

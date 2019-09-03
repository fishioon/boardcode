package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestImage2text(t *testing.T) {
	tess, _ := NewTess()
	data, _ := ioutil.ReadFile("/tmp/hello.png")
	res, _ := tess.Image2text(data)
	log.Printf("res:%s", res)
}

func Hello() {
	log.Printf("hello, world")
}

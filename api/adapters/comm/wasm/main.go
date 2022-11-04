//go:build js && wasm && go1.18 && !cgo
// +build js,wasm,go1.18,!cgo

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	js.Global().Set("showDate", js.FuncOf(showDate))
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	h2Element := js.Global().Get("document").
		Call("getElementById", "description")
	msg := fmt.Sprintf(
		"Sample 2 | Modifying the HTML (%s)",
		currentTime)
	h2Element.Set("innerHTML", msg)
	c := make(chan int)
	<-c
}

func showDate(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	h2 := document.Call("createElement", "h2")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("Current date & time is %s", currentTime)
	h2.Set("innerHTML", msg)

	div := document.Call("getElementById", "main")
	div.Call("appendChild", h2)
	return nil
}

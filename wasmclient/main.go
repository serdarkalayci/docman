package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall/js"

	"github.com/serdarkalayci/docman/client/dto"
)

func main() {
	js.Global().Set("getFolder", js.FuncOf(getFolder))
	js.Global().Set("getDocument", js.FuncOf(getDocument))
	getFolderContent("1")
	c := make(chan int)
	<-c
}

func getFolder(this js.Value, args []js.Value) interface{} {
	folderID := args[0].String()
	getFolderContent(folderID)
	return nil
}

func getFolderContent(folderID string) {
	document := js.Global().Get("document")
	navdiv := document.Call("getElementById", "nav")
	h2 := document.Call("createElement", "h2")
	if folderID == "" {
		folderID = "1"
	} else {
		if len(strings.Split(folderID, "/")) > 1 {
			folderID = strings.Split(folderID, "/")[1]
		}
	}
	go func() {
		response, err := http.Get(fmt.Sprintf("http://localhost:5550/folder/%s", folderID))
		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		var responseObject dto.FolderResponseDTO
		json.Unmarshal(responseData, &responseObject)

		h2.Set("innerHTML", responseObject.Name)
		navdiv.Set("innerHTML", "")
		navdiv.Call("appendChild", h2)
		if responseObject.Folders != nil {
			for i := 0; i < len(responseObject.Folders); i++ {
				spn := document.Call("createElement", "span")
				spn.Set("className", "linkbutton")
				img := document.Call("createElement", "input")
				img.Set("type", "image")
				img.Set("src", "images/folder.ico")
				img.Set("className", "btnimg")
				img.Set("innerHTML", responseObject.Folders[i].Name)
				img.Set("onclick", js.FuncOf(getFolder).Call("bind", img, responseObject.Folders[i].ID))
				btn := document.Call("createElement", "button")
				btn.Set("className", "linkbutton")
				btn.Set("innerHTML", responseObject.Folders[i].Name)
				btn.Set("onclick", js.FuncOf(getFolder).Call("bind", btn, responseObject.Folders[i].ID))
				spn.Call("appendChild", img)
				spn.Call("appendChild", btn)
				navdiv.Call("appendChild", spn)
			}
		}
		if responseObject.Documents != nil {
			for i := 0; i < len(responseObject.Documents); i++ {
				spn := document.Call("createElement", "span")
				spn.Set("className", "linkbutton")
				img := document.Call("createElement", "input")
				img.Set("type", "image")
				img.Set("src", "images/document.png")
				img.Set("className", "btnimg")
				img.Set("innerHTML", responseObject.Documents[i].Name)
				img.Set("onclick", js.FuncOf(getDocument).Call("bind", img, responseObject.Documents[i].ID))
				img.Set("src", "images/document.png")
				btn := document.Call("createElement", "button")
				btn.Set("className", "linkbutton")
				btn.Set("innerHTML", responseObject.Documents[i].Name)
				btn.Set("onclick", js.FuncOf(getDocument).Call("bind", btn, responseObject.Documents[i].ID))
				spn.Call("appendChild", img)
				spn.Call("appendChild", btn)
				navdiv.Call("appendChild", spn)
			}
		}
	}()
}

func getDocument(this js.Value, args []js.Value) interface{} {
	documentID := args[0].String()
	getDocumentContent(documentID)
	return nil
}

func getDocumentContent(documentID string) {
	document := js.Global().Get("document")
	condiv := document.Call("getElementById", "content")
	h2 := document.Call("createElement", "h2")
	if documentID == "" {
		return
	}
	if len(strings.Split(documentID, "/")) > 1 {
		documentID = strings.Split(documentID, "/")[1]
	}
	go func() {
		response, err := http.Get(fmt.Sprintf("http://localhost:5550/document/%s", documentID))
		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		condiv.Set("innerHTML", "")
		var responseObject dto.DocumentResponseDTO
		json.Unmarshal(responseData, &responseObject)

		h2.Set("innerHTML", responseObject.Name)
		spn := document.Call("createElement", "span")
		spn.Set("innerHTML", responseObject.Content)
		condiv.Call("appendChild", h2)
		condiv.Call("appendChild", spn)
	}()
}

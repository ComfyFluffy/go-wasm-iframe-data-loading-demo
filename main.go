package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"syscall/js"
)

var dataChannel = make(chan []byte)

func loadData(this js.Value, args []js.Value) interface{} {
	fmt.Println("Loading data...")
	inputArray := args[0]

	// Convert Uint8Array to Go byte slice
	length := inputArray.Get("length").Int()
	fmt.Println("Data Size:", length)
	goBytes := make([]byte, length)
	js.CopyBytesToGo(goBytes, inputArray)

	fmt.Println("Data Received! First 10 bytes:", goBytes[:10])

	dataChannel <- goBytes
	return nil
}

func main() {
	js.Global().Set("goLoadData", js.FuncOf(loadData))
	fmt.Println("Hello, WebAssembly!")

	data := <-dataChannel
	fmt.Println("Data loaded!")
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	// Read 1.txt from the zip file
	for _, file := range zipReader.File {
		fmt.Println("File:", file.Name)
		rc, err := file.Open()
		if err != nil {
			fmt.Println("Error:", err)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(rc)
		fmt.Println("File Content:", buf.String())
		rc.Close()
	}
	if err != nil {
		fmt.Println("Error:", err)
	}

}

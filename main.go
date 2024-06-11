package main

import (
	"fmt"
	"log"
	"os"
)

var (
	osCreate = os.Create
	file     *os.File
	fileRead = file.Read
)

func main() {
	file, err := osCreate("./foo")
	if err != nil {
		log.Panic("error creating file", err)
	}
	_, err = fileRead([]byte{})
	if err != nil {
		log.Panic("error reading file", err)
	}
	fmt.Println("file:", file)
}

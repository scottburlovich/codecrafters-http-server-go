package main

import (
	"fmt"
	"os"
)

func handleError(message string, err error) {
	fmt.Println(message+": ", err.Error())
	os.Exit(1)
}

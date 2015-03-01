// FileServer.go
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("WebFile at :9510")
	fileServer := http.FileServer(http.Dir("./"))
	err := http.ListenAndServe(":9510", fileServer)
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error", err.Error())
		os.Exit(1)
	}
}

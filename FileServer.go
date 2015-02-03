// FileServer.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("F:/Github/goproject/css"))
	err := http.ListenAndServe(":9998", fileServer)
	if err != nil {
		fmt.Println(err)
	}
}

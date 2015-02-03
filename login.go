package main

import (
	"fmt"
	"net/http"
)

func main() {
	h := http.FileServer(http.Dir("F:/Github/goproject"))
	err := http.ListenAndServe(":9999", h)
	if err != nil {
		fmt.Println("fuck")
	}
}

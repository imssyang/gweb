package main

import (
	"net/http"
)

func main() {
	dir := http.Dir("/root/.blog/demo")

	fileServer := http.FileServer(dir)

	http.Handle("/", fileServer)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ответ от сервера на порту 9002")
	})

	http.ListenAndServe(":9002", nil)
}

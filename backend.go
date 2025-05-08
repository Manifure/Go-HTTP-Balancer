package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Запуск сервера для тестирования, пример запуска: PORT=9001 go run backend.go
func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ответ от backend на порту %s\n", port)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

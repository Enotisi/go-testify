package main

import (
	"fmt"
	"net/http"

	"github.com/Enotisi/go-testify/internal/handlers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/cafe", handlers.MainHandle)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

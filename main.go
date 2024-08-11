package main

import (
	"fmt"
	"net/http"
	"os"
	"xerus/internal/middleware"
	"xerus/internal/view"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	mux.HandleFunc("GET /static/", view.ServeStaticFiles)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, view.Home)
	})

	fmt.Printf("server is running on port %s\n", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}

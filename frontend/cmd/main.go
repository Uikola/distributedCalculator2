package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir("frontend")))

	http.ListenAndServe(os.Getenv("FRONTEND_PORT"), nil)
}

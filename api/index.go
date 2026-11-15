package main

import (
	"MokutonCoffee-API/server"
	"net/http"
)

// Handler dipakai oleh Vercel Go runtime sebagai serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	server.New().ServeHTTP(w, r)
}

// func main kosong agar `go build ./...` lolos; Vercel memakai Handler.
func main() {}

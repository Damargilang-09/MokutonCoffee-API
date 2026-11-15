package main

import (
	"MokutonCoffee-API/server"
	"log"
	"os"
)

func main() {
	if os.Getenv("VERCEL") != "" {
		// Di Vercel engine dibangun oleh handler api/index.go per cold start.
		return
	}

	log.Println("MokutonCoffee-API listening on :8000")
	if err := server.New().Run(":8000"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

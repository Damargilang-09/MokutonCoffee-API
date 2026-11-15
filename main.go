package main

import (
	"MokutonCoffee-API/server"
	"log"
	"os"
)

func main() {
	// Vercel (executable Go runtime) menjalankan binary ini sebagai
	// long-running server yang listen di $PORT; lokal default :8000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("MokutonCoffee-API listening on :" + port)
	if err := server.New().Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}

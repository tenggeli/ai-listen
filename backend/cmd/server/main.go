package main

import (
	"log"

	"listen/backend/internal/bootstrap"
)

func main() {
	server := bootstrap.NewServer()
	log.Printf("listen backend started on :%s", server.Port)
	if err := server.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

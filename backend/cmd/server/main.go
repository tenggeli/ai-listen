package main

import (
	"log"

	"ai-listen/backend/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap app failed: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("run app failed: %v", err)
	}
}

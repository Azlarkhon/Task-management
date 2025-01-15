package main

import (
	"log"
	"task-management/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"url-shortener/internal/app"
)

func main() {
	env := flag.String("env", "local", "environment")
	flag.Parse()

	if env != nil && *env == "local" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		filename := fmt.Sprintf("%s\\.env.%s", dir, *env)
		if err = godotenv.Load(filename); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	app.Run(*env)
}

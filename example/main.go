package main

import (
	server "github.com/ArturChopikian/csv_http_server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading env file: ", err)
	}

	nameOfFolder := os.Getenv("FOLDER_WITH_FILES")
	address := os.Getenv("ADDRESS")

	s, err := server.NewCSVServer(nameOfFolder, address)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.MapHandlers(); err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

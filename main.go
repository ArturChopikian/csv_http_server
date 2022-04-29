package main

import (
	"csv_server/project"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading env file: ", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("error while getting pwd: ", err)
	}

	nameOfFolder := os.Getenv("FOLDER_WITH_FILES")
	port := os.Getenv("PORT")

	pathToFolder := pwd + project.PathSeparator + nameOfFolder

	err = filepath.Walk(pathToFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// check if file is scv
		file := strings.Split(info.Name(), ".")
		if len(file) > 1 && file[1] == "csv" {
			// create handler for each csv file
			http.HandleFunc("/"+info.Name(), getFile(path))
			fmt.Printf("Create handler: localhost%s/%s\n", port, info.Name())
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CSV server start")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func getFile(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		fileBytes, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal("error while reading the file:", err)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")

		_, err = w.Write(fileBytes)
		if err != nil {
			log.Println("error while writing in response:", err)
		}
	}
}

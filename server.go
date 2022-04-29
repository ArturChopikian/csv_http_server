package cvs_http_server

import (
	"errors"
	"fmt"
	"github.com/ArturChopikian/csv_http_server/project"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CSVServer - contain to fields
// folderWithFiles - save file where we save our
// address - address where server will be run
type CSVServer struct {
	folderWithFiles string
	address         string
}

func NewCSVServer(folder, address string) (*CSVServer, error) {
	if folder == "" {
		return &CSVServer{}, errors.New("folder can not be blank")
	}
	if address == "" {
		return &CSVServer{}, errors.New("address can not be blank")
	}

	return &CSVServer{
		folderWithFiles: folder,
		address:         address,
	}, nil
}

// Run - run the server and define simple handler for checking if server live
func (s *CSVServer) Run() error {
	http.HandleFunc("/", hello)

	return http.ListenAndServe(s.address, nil)
}

// MapHandlers - walks through all files of the picked folder
// and call function which create handler for each csv file
func (s *CSVServer) MapHandlers() error {

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error while getting pwd: %v", err)
	}

	nameOfFolder := os.Getenv("FOLDER_WITH_FILES")

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
			fmt.Printf("Create handler: %s/%s\n", s.address, info.Name())
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// getFile - wrapped handler which take name of csv file and create handler with this name
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

// hello - simple handler to check if server is working
func hello(w http.ResponseWriter, _ *http.Request) {
	if _, err := io.WriteString(w, "CSV server work"); err != nil {
		log.Println(err)
	}
}

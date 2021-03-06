package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/adamluzsi/frameless/example/channels"
	"github.com/adamluzsi/frameless/example/usecases"
	"github.com/adamluzsi/frameless/storages"
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storages.NewLocal(filepath.Join(wd, "db"))
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	cases := usecases.NewUseCases(storage)
	handler := channels.NewHTTPHandler(cases)
	server := &http.Server{Addr: ":8080", Handler: handler}
	cli := channels.NewCLI(os.Stdout, cases, server)

	if err := cli.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"log"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/handlers"

	_ "github.com/lib/pq"
)

func main() {

	srv := handlers.SetupServer()

	log.Printf("Serving on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}

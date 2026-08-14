package main

import (
	"log"

	"github.com/geneowak/file-storage-s3-golang/handlers"

	_ "github.com/lib/pq"
)

func main() {

	cfg := handlers.InitConfig()
	srv := handlers.SetupServer(&cfg)

	log.Printf("Serving on: http://localhost:%s/app/\n", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}

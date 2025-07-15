package main

import (
	"flag"
	"log"

	"github.com/gnojus/wedl/api"
)

func main() {
	port := flag.String("port", "8080", "Port to run the API server on")
	flag.Parse()

	server := api.NewServer(*port)
	log.Fatal(server.Start())
}
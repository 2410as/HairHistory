package main

import (
	"log"
	"net/http"

	"github.com/annasakai/hairhistorymemo/apps/main/app/controller"
)

func main() {
	handler := controller.NewRouter(wireDeps())

	addr := ":8080"
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

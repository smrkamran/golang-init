package main

import (
	"fmt"
	"go-stocks/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	fmt.Println("Starting server on the port 9000...")

	log.Fatal(http.ListenAndServe(":9000", r))
}

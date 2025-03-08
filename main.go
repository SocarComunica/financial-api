package main

import (
	"log"

	"github.com/socarcomunica/financial-api/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/brutella/geizhals"
)

func main() {
	pr, err := geizhals.GetProduct("1696985")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s %.2f - %.2f\n", pr.Id, pr.MinPrice, pr.MaxPrice)
}

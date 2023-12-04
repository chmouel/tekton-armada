package main

import (
	"log"
	"os"

	armada "github.com/chmouel/armadas/pkg"
)

func main() {
	if err := armada.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

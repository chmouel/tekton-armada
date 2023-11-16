package main

import (
	"log"
	"os"

	andromede "github.com/chmouel/andromede/pkg"
)

func main() {
	if err := andromede.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

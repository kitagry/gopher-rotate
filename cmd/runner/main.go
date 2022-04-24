package main

import (
	"log"

	"github.com/kitagry/gopherrotate"
)

func main() {
	mascot := gopherrotate.NewMascot()
	if err := mascot.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/kitagry/gopherrotate"
)

func main() {
	mascot := gopherrotate.NewMascot()
	mascot.Say("Hello")
	if err := mascot.Run(); err != nil {
		log.Fatal(err)
	}
}

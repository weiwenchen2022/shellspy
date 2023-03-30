package main

import (
	"log"
	"os"
	"shellspy"
)

func init() {
	log.SetFlags(0)
}

func main() {
	s := shellspy.NewSession(os.Stdin, os.Stdout)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

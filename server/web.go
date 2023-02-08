package server

import (
	"log"
)

func RunWeb(listen string) {
	if err := NewRouter().Run(listen); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/roarkanize/qsync"
	"log"
)

func main() {
	log.Println("initializing database")
	db, err := qsync.InitDB("qsync.db")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("closing database")
	err = db.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

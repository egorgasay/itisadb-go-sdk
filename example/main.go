package main

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
)

// main to run this test, itisadb must be run on :8888.
func main() {
	db, err := itisadb.New(":8888")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetOne(ctx, "qwe", "111", false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetOne(ctx, "qwe")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "111" {
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", get)
	}
}

package main

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
)

// main to run this test, itisadb must be run on :8888.
func main() {
	ctx := context.TODO()

	db, err := itisadb.New(ctx, ":8888")
	if err != nil {
		return
	}

	err = db.SetOne(ctx, "qwe", "111").Error()
	if err != nil {
		log.Fatalln(err)
	}

	res := db.GetOne(ctx, "qwe")
	if err != nil {
		log.Fatalln(err)
	}

	if x := res.Unwrap(); x != "111" {
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", x)
	}
}

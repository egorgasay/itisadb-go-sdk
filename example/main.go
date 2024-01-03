package main

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
)

// main to run this test, itisadb must be run on :8888.
func main() {
	ctx := context.TODO()

	db := itisadb.New(ctx, ":8888").Unwrap()

	db.SetOne(ctx, "qwe", "111").Unwrap()

	if x := db.GetOne(ctx, "qwe").Unwrap().Value; x != "111" {
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", x)
	}
}

package main

import (
	"context"
	grpcisclient "github.com/egorgasay/gRPCis-client"
	"log"
)

// main to run this test, gRPCis must be run on :800.
func main() {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetOne(ctx, "qwe", "111")
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOne(ctx, "qwe")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "111" {
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", get)
	}
}

package grpcisclient_test

import (
	"context"
	grpcisclient "grpcis-client"
	"log"
	"testing"
)

// TestSetGet to run this test, grpcis must be run on :800.
func TestSetGet(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetOnce(ctx, "qwe", "111")
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOnce(ctx, "qwe")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "111" {
		t.Fatal("Wrong value")
	}
}

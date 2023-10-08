package itisadb

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk/api"
	"testing"
)

func TestAuth(t *testing.T) {
	ctx := context.Background()
	cl, err := New(ctx, ":8888")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// method that requires auth
	resp, err := cl.cl.Servers(withAuth(ctx), &api.ServersRequest{})
	if err != nil {
		t.Fatalf("Servers: %v", err)
	}

	t.Log("Auth OK")
	t.Log(resp.ServersInfo)
}

// TestSetToGetFrom to run this test, itisadb must be run on :8888.
func TestSetToGetFrom(t *testing.T) {
	db, err := New(context.TODO(), ":8888")
	if err != nil {
		return
	}

	var snum int32 = 1

	ctx := context.TODO()
	err = db.set(ctx, "fff", "qqq", SetOptions{
		Server:   &snum,
		Uniques:  false,
		ReadOnly: false,
	}).Err()
	if err != nil {
		t.Fatal(err)
	}

	get, err := db.get(ctx, "fff", GetOptions{Server: &snum}).ValueAndErr()
	if err != nil {
		t.Fatal(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

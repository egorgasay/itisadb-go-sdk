package itisadb

import (
	"context"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"testing"
)

func TestAuth(t *testing.T) {
	ctx := context.Background()
	cl := New(ctx, ":8888").Unwrap()

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
	db := New(context.TODO(), ":8888").Unwrap()

	var snum int32 = 1

	ctx := context.TODO()
	err := db.set(ctx, "fff", "qqq", SetOptions{
		Server:   &snum,
		Unique:   false,
		ReadOnly: false,
	}).Error()
	if err != nil {
		t.Fatal(err)
	}

	get := db.get(ctx, "fff", GetOptions{Server: &snum}).Unwrap()

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

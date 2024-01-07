package itisadb

import (
	"context"
	"testing"
)

// TestSetToGetFrom to run this test, itisadb must be run on :8888.
func TestSetToGetFrom(t *testing.T) {
	db := New(context.TODO(), ":8888").Unwrap()

	var snum int32 = 1

	ctx := context.TODO()
	err := db.set(ctx, "fff", "qqq", SetOptions{
		Server:   snum,
		Unique:   false,
		ReadOnly: false,
	}).Error()
	if err != nil {
		t.Fatal(err)
	}

	get := db.get(ctx, "fff", GetOptions{Server: snum}).Unwrap()

	if get.Value != "qqq" {
		t.Fatal("Wrong value")
	}
}

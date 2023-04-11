package grpcisclient_test

import (
	"context"
	grpcisclient "grpcis-client"
	"log"
	"reflect"
	"testing"
)

// TestSetGetOne to run this test, grpcis must be run on :800.
func TestSetGetOne(t *testing.T) {
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
		t.Fatal("Wrong value")
	}
}

// TestSetToGetFrom to run this test, grpcis must be run on :800.
func TestSetToGetFrom(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetTo(ctx, "fff", "qqq", 1)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetFrom(ctx, "fff", 1)
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetToDBGetFromDB to run this test, grpcis must be run on :800.
// TODO: Add edge cases.
func TestSetToDBGetFromDB(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetToDB(ctx, "db_key", "qqq")
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetFromDB(ctx, "db_key")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetToAllGet to run this test, grpcis must be run on :800.
// TODO: Add edge cases.
func TestSetToAllGet(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetToAll(ctx, "all_key", "qqq")
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOne(ctx, "all_key")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetManyGetMany to run this test, grpcis must be run on :800.
// TODO: Add edge cases.
func TestSetManyGetMany(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	m := map[string]string{
		"m1": "k1",
		"m2": "k2",
		"m3": "k3",
		"m4": "k4",
		"m5": "k5",
	}

	ctx := context.TODO()
	err = grpcis.SetMany(ctx, m)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOne(ctx, "m2")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "k2" {
		t.Fatal("Wrong value")
	}

	k := []string{"m1", "m2", "m3", "m4", "m5"}
	res, err := grpcis.GetMany(ctx, k)
	if err != nil {
		log.Fatalln(err)
	}

	if !reflect.DeepEqual(res, m) {
		t.Fatal("Wrong value")
	}
}

// TestSetManyOptsGetManyOpts to run this test, grpcis must be run on :800.
// TODO: Add edge cases.
func TestSetManyOptsGetManyOpts(t *testing.T) {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	me := map[string]string{
		"mo1": "k1",
		"mo2": "k2",
		"mo3": "k3",
		"mo4": "k4",
		"mo5": "k5",
	}

	m := map[string]grpcisclient.Value{
		"mo1": {Value: "k1", Opts: grpcisclient.Opts{Server: 1}},
		"mo2": {Value: "k2", Opts: grpcisclient.Opts{Server: 1}},
		"mo3": {Value: "k3", Opts: grpcisclient.Opts{Server: -1}},
		"mo4": {Value: "k4", Opts: grpcisclient.Opts{Server: -2}},
		"mo5": {Value: "k5", Opts: grpcisclient.Opts{Server: -3}},
	}

	ctx := context.TODO()
	err = grpcis.SetManyOpts(ctx, m)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOne(ctx, "mo2")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "k2" {
		t.Fatal("Wrong value")
	}

	k := []grpcisclient.Key{
		{Key: "mo1", Opts: grpcisclient.Opts{Server: 1}},
		{Key: "mo2", Opts: grpcisclient.Opts{Server: 1}},
		{Key: "mo3", Opts: grpcisclient.Opts{Server: -1}},
		{Key: "mo4", Opts: grpcisclient.Opts{Server: 0}},
		{Key: "mo5", Opts: grpcisclient.Opts{Server: 0}},
	}

	res, err := grpcis.GetManyOpts(ctx, k)
	if err != nil {
		log.Fatalln(err)
	}

	if !reflect.DeepEqual(res, me) {
		t.Fatal("Wrong value")
	}
}

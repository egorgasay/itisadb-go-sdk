package itisadb_test

import (
	"bufio"
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"modernc.org/strutil"
)

// TestSetGetOne to run this test, itisadb must be run on :800.
func TestSetGetOne(t *testing.T) {
	db, err := itisadb.New(":800")
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
		t.Fatal("Wrong value")
	}
}

// TestSetToGetFrom to run this test, itisadb must be run on :800.
func TestSetToGetFrom(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetTo(ctx, "fff", "qqq", 1, false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetFrom(ctx, "fff", 1)
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetToDBGetFromDB to run this test, itisadb must be run on :800.
// TODO: Add edge cases.
func TestSetToDBGetFromDB(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetToDB(ctx, "db_key", "qqq", false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetFromDB(ctx, "db_key")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetToAllGet to run this test, itisadb must be run on :800.
// TODO: Add edge cases.
func TestSetToAllGet(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetToAll(ctx, "all_key", "qqq", false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetOne(ctx, "all_key")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetManyGetMany to run this test, itisadb must be run on :800.
// TODO: Add edge cases.
func TestSetManyGetMany(t *testing.T) {
	db, err := itisadb.New(":800")
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
	err = db.SetMany(ctx, m, false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetOne(ctx, "m2")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "k2" {
		t.Fatal("Wrong value")
	}

	k := []string{"m1", "m2", "m3", "m4", "m5"}
	res, err := db.GetMany(ctx, k)
	if err != nil {
		log.Fatalln(err)
	}

	if !reflect.DeepEqual(res, m) {
		t.Fatal("Wrong value")
	}
}

// TestSetManyOptsGetManyOpts to run this test, itisadb must be run on :800.
// TODO: Add edge cases.
func TestSetManyOptsGetManyOpts(t *testing.T) {
	db, err := itisadb.New(":800")
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

	m := map[string]itisadb.Value{
		"mo1": {Value: "k1", Opts: itisadb.Opts{Server: 1}},
		"mo2": {Value: "k2", Opts: itisadb.Opts{Server: 1}},
		"mo3": {Value: "k3", Opts: itisadb.Opts{Server: -1}},
		"mo4": {Value: "k4", Opts: itisadb.Opts{Server: -2}},
		"mo5": {Value: "k5", Opts: itisadb.Opts{Server: -3}},
	}

	ctx := context.TODO()
	err = db.SetManyOpts(ctx, m, false)
	if err != nil {
		log.Fatalln(err)
	}

	get, err := db.GetOne(ctx, "mo2")
	if err != nil {
		log.Fatalln(err)
	}

	if get != "k2" {
		t.Fatal("Wrong value")
	}

	k := []itisadb.Key{
		{Key: "mo1", Opts: itisadb.Opts{Server: 1}},
		{Key: "mo2", Opts: itisadb.Opts{Server: 1}},
		{Key: "mo3", Opts: itisadb.Opts{Server: -1}},
		{Key: "mo4", Opts: itisadb.Opts{Server: 0}},
		{Key: "mo5", Opts: itisadb.Opts{Server: 0}},
	}

	res, err := db.GetManyOpts(ctx, k)
	if err != nil {
		log.Fatalln(err)
	}

	if !reflect.DeepEqual(res, me) {
		t.Fatal("Wrong value")
	}
}

// Benchmark test for SetOne.
func BenchmarkSetOne(b *testing.B) {
	db, err := itisadb.New(":800")
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.TODO()
	j := 1000
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k := fmt.Sprint(j)
		j++
		b.StartTimer()
		err = db.SetOne(ctx, k, "value", false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestSetGetOneToIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	index, err := db.Index(context.TODO(), "User")
	if err != nil {
		log.Fatalln(err)
	}

	err = index.Set(context.TODO(), "Name", "Max", false)
	if err != nil {
		log.Fatalln(err)
	}

	value, err := index.Get(context.TODO(), "Name")
	if err != nil {
		log.Fatalln(err)
	}

	if value != "Max" {
		t.Fatalf("Wrong value [%s] wanted [%s]\n", value, "Max")
	}

	/// CAR

	car, err := db.Index(context.TODO(), "Car")
	if err != nil {
		log.Fatalln(err)
	}

	err = car.Set(context.TODO(), "Name", "MyCar", false)
	if err != nil {
		log.Fatalln(err)
	}

	value, err = car.Get(context.TODO(), "Name")
	if err != nil {
		log.Fatalln(err)
	}

	if value != "MyCar" {
		t.Fatal("Wrong value")
	}

	/// WHEEL

	wheel, err := car.Index(context.TODO(), "Wheel")
	if err != nil {
		log.Fatalln(err)
	}

	err = wheel.Set(context.TODO(), "Color", "Black", false)
	if err != nil {
		log.Fatalln(err)
	}

	value, err = wheel.Get(context.TODO(), "Color")
	if err != nil {
		log.Fatalln(err)
	}

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	/// TRAILER

	trailer, err := car.Index(context.TODO(), "Trailer")
	if err != nil {
		log.Fatalln(err)
	}

	err = trailer.Set(context.TODO(), "Color", "Red", false)
	if err != nil {
		log.Fatalln(err)
	}

	/// TEST WHEEL & TRAILER AREAS ARE STILL WORKING

	value, err = wheel.Get(context.TODO(), "Color")
	if err != nil {
		log.Fatalln(err)
	}

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	value, err = trailer.Get(context.TODO(), "Color")
	if err != nil {
		log.Fatalln(err)
	}

	if value != "Red" {
		t.Fatal("Wrong value")
	}
}

// TestIsIndex
func TestIsIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	name := fmt.Sprintf("%d", time.Now().Unix())
	index, err := db.Index(context.TODO(), name)
	if err != nil {
		log.Fatalf("Index error: %s\n", err)
	}

	time.Sleep(time.Second)

	ctx := context.TODO()
	if ok, err := db.IsIndex(ctx, index.GetName()); err != nil || !ok {
		t.Fatal("Not an index, but should be")
	}

	name = fmt.Sprintf("%d", time.Now().Unix())
	newIndex, err := index.Index(ctx, index.GetName())
	if err != nil {
		log.Fatalf("Index error: %s\n", err)
	}

	time.Sleep(time.Second)

	if ok, err := db.IsIndex(ctx, newIndex.GetName()); err != nil || !ok {
		t.Fatal("Not an index, but should be")
	}

	if ok, _ := db.IsIndex(ctx, fmt.Sprintf("%d", time.Now().Unix())); ok {
		t.Fatal("Index, but should not be")
	}
}

func TestSize(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	index, err := db.Index(context.TODO(), fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		log.Fatalln(err)
	}

	size, err := index.Size(context.TODO())
	if err != nil {
		t.Fatalf("Error %v\n", err)
	} else if size != 0 {
		t.Fatalf("Wrong size %d\n", size)
	}

	for i := 0; i < 100; i++ {
		k := fmt.Sprint(i)
		v := fmt.Sprint(i)
		if err = index.Set(context.TODO(), k, v, false); err != nil {
			t.Fatalf("Set error %v\n", err)
		}
	}

	if size, err := index.Size(context.TODO()); err != nil || size != 100 {
		t.Fatalf("Wrong size %d\n", size)
	}
}

// TestGetIndex
func TestGetIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	name := fmt.Sprintf("%d", time.Now().Unix())
	index, err := db.Index(context.TODO(), name)
	if err != nil {
		log.Fatalf("Index error: %s\n", err)
	}

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	for k, v := range data {
		if err = index.Set(context.TODO(), k, v, false); err != nil {
			t.Fatalf("Set error %v\n", err)
		}
	}

	ctx := context.TODO()
	m, err := index.GetIndex(ctx)
	if err != nil {
		t.Fatalf("GetIndex error %v\n", err)
	}

	if !reflect.DeepEqual(data, m) {
		t.Fatalf("Wrong data %v\n", m)
	}
}

// Test to define rps for SetOne.
func TestSetOneRPS(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 25000

	log.Println("Total actions:", gnum)
	log.Println("RPS:", maxRPS)

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		ints[i] = fmt.Sprint(i)
	}

	log.Println("Hops:", gnum/maxRPS)

	var total time.Duration
	for tt := gnum / maxRPS; tt > 0; tt-- {
		var wg sync.WaitGroup
		wg.Add(maxRPS)

		var wgSent sync.WaitGroup
		wgSent.Add(maxRPS)
		ctx := context.TODO()
		for i := 0; i < maxRPS; i++ {
			wg.Done()
			go func(i int) {
				wg.Wait()
				db.GetOne(ctx, ints[i])
				wgSent.Done()
			}(i)

		}
		wg.Wait()
		start := time.Now()
		wgSent.Wait()
		total += time.Since(start)
		t.Log(time.Since(start))
	}
	t.Log(total)
}

func TestDistinct(t *testing.T) {
	f, err := os.Open("/tmp/log14/transactionLogger")
	if err != nil {
		t.Fatal(err)
	}

	var keys = make(map[string]struct{}, 16000)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		action := scanner.Text()
		decode, err := strutil.Base64Decode([]byte(action))
		if err != nil {
			return
		}

		split := strings.Split(string(decode), " ")
		key := split[1]
		keys[key] = struct{}{}
	}

	t.Log(len(keys))
}

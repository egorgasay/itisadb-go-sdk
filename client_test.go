package itisadb_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
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

func TestDelete(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.TODO()

	num := rand.Int31()
	n := fmt.Sprint(num)
	err = db.SetOne(ctx, "key_for_delete"+n, "value_for_delete", false)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Del(ctx, "key_for_delete"+n)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.GetOne(ctx, "key_for_delete"+n)
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Key should be deleted, but %v", err)
	}
}

func TestDeleteIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteIndex" + n

	indx, err := db.Index(ctx, name)
	if err != nil {
		log.Fatalln(err)
	}
	err = indx.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalln(err)
	}

	err = indx.DeleteIndex(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	indx, err = db.Index(ctx, name)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Index should be deleted")
	}

	// TEST DELETE INNER INDEX

	name = "inner_index"
	inner, err := indx.Index(ctx, name)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	err = inner.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	err = inner.DeleteIndex(ctx)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	inner, err = indx.Index(ctx, name)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Inner Index should be deleted")
	}
}

func TestDeleteAttr(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteAttr" + n

	indx, err := db.Index(ctx, name)
	if err != nil {
		log.Fatalln(err)
	}
	err = indx.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalln(err)
	}

	err = indx.DeleteAttr(ctx, "key_for_delete")
	if err != nil {
		log.Fatalln(err)
	}

	indx, err = db.Index(ctx, name)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Key should be deleted")
	}

	// TEST DELETE ATTR INNER INDEX

	name = "inner_index"
	inner, err := indx.Index(ctx, name)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	err = inner.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	err = inner.DeleteAttr(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	inner, err = indx.Index(ctx, name)
	if err != nil {
		log.Fatalf("Inner index %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Inner Index key should be deleted")
	}

	// TEST DELETE ATTR (INNER INDEX) KEY DOES NOT EXIST

	err = inner.DeleteAttr(ctx, "key_for_delete_does_not_exist")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Inner Index key shouldn't be deleted: %v", err)
	}

}

func TestAttachIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.TODO()
	originalIndex := "TestAttachIndex"
	indx, err := db.Index(ctx, originalIndex)
	if err != nil {
		log.Fatalln(err)
	}

	attachedIndex := "TestAttachIndex2"
	inner, err := db.Index(ctx, attachedIndex)
	if err != nil {
		log.Fatalln(err)
	}

	err = inner.Set(ctx, "key_for_attach", "value_for_attach", false)
	if err != nil {
		log.Fatalf("set Inner index %v: %v", attachedIndex, err)
	}

	err = indx.Attach(ctx, inner.GetName())
	if err != nil {
		log.Fatalln(err)
	}

	innerCopy, err := indx.Index(ctx, attachedIndex)
	if err != nil {
		log.Fatalf("main switch index %v: %v", attachedIndex, err)
	}

	err = innerCopy.Set(ctx, "key_for_attach3", "value_for_attach3", false)
	if err != nil {
		log.Fatalf("set Inner index %v: %v", attachedIndex, err)
	}

	err = inner.Set(ctx, "key_for_attach4", "value_for_attach4", false)
	if err != nil {
		log.Fatalf("set Inner index %v: %v", attachedIndex, err)
	}

	originalAttached, err := inner.GetIndex(ctx)
	if err != nil {
		log.Fatalf("get Inner index %v: %v", attachedIndex, err)
	}

	copiedAttached, err := innerCopy.GetIndex(ctx)
	if err != nil {
		log.Fatalf("get Inner index %v: %v", attachedIndex, err)
	}

	if !reflect.DeepEqual(originalAttached, copiedAttached) {
		t.Fatalf("Inner index not equal original index:  %v != %v", originalAttached, copiedAttached)
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

	time.Sleep(time.Second)

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

func TestClient_StructToIndex(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	type Wheel struct {
		Color string
		Size  string
	}

	type Trailer struct {
		Color string
		Size  string
	}

	type Car struct {
		Name    string
		Wheel   *Wheel
		Trailer *Trailer
	}

	type User struct {
		Name  string
		Age   int
		Email string
		Car   *Car
	}

	user := User{
		Name:  "Max",
		Age:   25,
		Email: "max@mail.ru",
		Car: &Car{
			Name: "MyCar",
			Wheel: &Wheel{
				Color: "Black",
				Size:  "Big",
			},
			Trailer: &Trailer{
				Color: "Red",
				Size:  "Big",
			},
		},
	}

	index, err := db.StructToIndex(context.TODO(), "User1", user)
	if err != nil {
		t.Fatal(err)
	}

	us := &User{}
	err = db.IndexToStruct(context.TODO(), "User1", us)
	if err != nil {
		t.Fatal(err)
	}

	mu := map[string]string{
		"Name":  "Max",
		"Car":   "\n\tName: MyCar\n\tWheel: \n\t\tSize: Big\n\t\tColor: Black\n\tTrailer: \n\t\t\tColor: Red\n\t\t\tSize: Big",
		"Age":   "25",
		"Email": "max@mail.ru",
	}

	mi, err := index.GetIndex(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	splittedMi := strings.Split(strings.Replace(strings.Replace(mi["Car"], "\t", "", -1), "\n", "", -1), "")
	splittedMu := strings.Split(strings.Replace(strings.Replace(mu["Car"], "\t", "", -1), "\n", "", -1), "")

	if !IsTheSameArray(splittedMi, splittedMu) {
		t.Fatalf("want\n%v,\ngot \n%v", splittedMu, splittedMi)
	}

	delete(mu, "Car")
	delete(mi, "Car")

	if !reflect.DeepEqual(mi, mu) {
		t.Fatalf("want\n%v,\ngot \n%v", mu, mi)
	}
}
func IsTheSameArray[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	tmp := make(map[T]T)
	for _, el := range a {
		tmp[el] = el
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
}

func TestClient_IndexToStruct(t *testing.T) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	type Wheel struct {
		Color string
		Size  string
	}

	type Trailer struct {
		Color string
		Size  string
	}

	type Car struct {
		Name    string
		Wheel   *Wheel
		Trailer *Trailer
	}

	type User struct {
		Name  string
		Age   int
		Email string
		Car   *Car
	}

	user := User{
		Name:  "Max",
		Age:   25,
		Email: "max@mail.ru",
		Car: &Car{
			Name: "MyCar",
			Wheel: &Wheel{
				Color: "Black",
				Size:  "Big",
			},
			Trailer: &Trailer{
				Color: "Red",
				Size:  "Big",
			},
		},
	}

	_, err = db.StructToIndex(context.TODO(), "User1", user)
	if err != nil {
		t.Fatal(err)
	}

	us := &User{}
	err = db.IndexToStruct(context.TODO(), "User1", us)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp(user, *us) {
		t.Fatal("Not equal")
	}

}

func cmp[V any](a, b V) bool {
	aVal, bVal := reflect.ValueOf(a), reflect.ValueOf(b)
	return cmpReflect(aVal, bVal)
}

func cmpReflect(a, b reflect.Value) bool {
	for i := 0; i < a.NumField(); i++ {
		fa, fb := a.Field(i), b.Field(i)
		if fa.Type().Kind() == reflect.Ptr {
			if !cmpReflect(fa.Elem(), fb.Elem()) {
				return false
			}
		} else if fa.Interface() != fb.Interface() {
			fmt.Println(a.Field(i).Interface(), b.Field(i).Interface())
			return false
		}
	}
	return true
}

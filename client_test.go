package itisadb_test

import (
	"bufio"
	"context"
	"encoding/json"
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

var _ctx = context.TODO()

// TestSetGetOne to run this test, itisadb must be run on :8888.
func TestSetGetOne(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetOne(ctx, "qwe", "111", false)
	if err != nil {
		t.Fatal(err)
	}

	get, err := db.Get(ctx, "qwe")
	if err != nil {
		t.Fatal(err)
	}

	if get != "111" {
		t.Fatal("Wrong value")
	}
}

// TestSetToGetFrom to run this test, itisadb must be run on :8888.
func TestSetToGetFrom(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetTo(ctx, "fff", "qqq", 1, false)
	if err != nil {
		t.Fatal(err)
	}

	get, err := db.GetFrom(ctx, "fff", 1)
	if err != nil {
		t.Fatal(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetToAllGet to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetToAllGet(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetToAll(ctx, "all_key", "qqq", false)
	if err != nil {
		t.Fatal(err)
	}

	get, err := db.Get(ctx, "all_key")
	if err != nil {
		t.Fatal(err)
	}

	if get != "qqq" {
		t.Fatal("Wrong value")
	}
}

// TestSetManyGetMany to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetManyGetMany(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
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
		t.Fatal(err)
	}

	get, err := db.Get(ctx, "m2")
	if err != nil {
		t.Fatal(err)
	}

	if get != "k2" {
		t.Fatal("Wrong value")
	}

	k := []string{"m1", "m2", "m3", "m4", "m5"}
	res, err := db.GetMany(ctx, k)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, m) {
		t.Fatal("Wrong value")
	}
}

// TestSetManyOptsGetManyOpts to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetManyOptsGetManyOpts(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
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
		t.Fatal(err)
	}

	get, err := db.Get(ctx, "mo2")
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, me) {
		t.Fatal("Wrong value")
	}
}

func TestDelete(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	num := rand.Int31()
	n := fmt.Sprint(num)
	err = db.SetOne(ctx, "key_for_delete"+n, "value_for_delete", false)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Del(ctx, "key_for_delete"+n)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Get(ctx, "key_for_delete"+n)
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Key should be deleted, but %v", err)
	}
}

func TestDeleteObject(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteObject" + n

	indx, err := db.Object(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	err = indx.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if err != nil {
		t.Fatal(err)
	}

	err = indx.DeleteObject(ctx)
	if err != nil {
		t.Fatal(err)
	}

	indx, err = db.Object(ctx, name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Object should be deleted")
	}

	// TEST DELETE INNER OBJECT

	name = "inner_object"
	inner, err := indx.Object(ctx, name)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	err = inner.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	err = inner.DeleteObject(ctx)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	inner, err = indx.Object(ctx, name)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Inner Object should be deleted")
	}
}

func TestDeleteAttr(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteAttr" + n

	indx, err := db.Object(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	err = indx.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if err != nil {
		t.Fatal(err)
	}

	err = indx.DeleteAttr(ctx, "key_for_delete")
	if err != nil {
		t.Fatal(err)
	}

	indx, err = db.Object(ctx, name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = indx.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Key should be deleted")
	}

	// TEST DELETE ATTR INNER OBJECT

	name = "inner_object"
	inner, err := indx.Object(ctx, name)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	err = inner.Set(ctx, "key_for_delete", "value_for_delete", false)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	err = inner.DeleteAttr(ctx, "key_for_delete")
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	inner, err = indx.Object(ctx, name)
	if err != nil {
		log.Fatalf("Inner object %v: %v", name, err)
	}

	_, err = inner.Get(ctx, "key_for_delete")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Inner Object key should be deleted")
	}

	// TEST DELETE ATTR (INNER OBJECT) KEY DOES NOT EXIST

	err = inner.DeleteAttr(ctx, "key_for_delete_does_not_exist")
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Inner Object key shouldn't be deleted: %v", err)
	}

}

func TestAttachObject(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()

	originalObject := "TestAttachObject"
	indx, err := db.Object(ctx, originalObject)
	if err != nil {
		t.Fatal(err)
	}

	attachedObject := "TestAttachObject2"
	inner, err := db.Object(ctx, attachedObject)
	if err != nil {
		t.Fatal(err)
	}

	if err := indx.DeleteObject(ctx); err != nil {
		t.Fatalf("Delete object %v: %v", originalObject, err)
	}

	if err := inner.DeleteObject(ctx); err != nil {
		t.Fatalf("Delete object %v: %v", originalObject, err)
	}

	originalObject = "TestAttachObject"
	indx, err = db.Object(ctx, originalObject)
	if err != nil {
		t.Fatal(err)
	}

	attachedObject = "TestAttachObject2"
	inner, err = db.Object(ctx, attachedObject)
	if err != nil {
		t.Fatal(err)
	}

	err = inner.Set(ctx, "key_for_attach", "value_for_attach", false)
	if err != nil {
		log.Fatalf("set Inner object %v: %v", attachedObject, err)
	}

	err = indx.Attach(ctx, inner.GetName())
	if err != nil {
		t.Fatal(err)
	}

	innerCopy, err := indx.Object(ctx, attachedObject)
	if err != nil {
		log.Fatalf("main switch object %v: %v", attachedObject, err)
	}

	err = innerCopy.Set(ctx, "key_for_attach3", "value_for_attach3", false)
	if err != nil {
		log.Fatalf("set Inner object %v: %v", attachedObject, err)
	}

	err = inner.Set(ctx, "key_for_attach4", "value_for_attach4", false)
	if err != nil {
		log.Fatalf("set Inner object %v: %v", attachedObject, err)
	}

	originalAttached, err := inner.JSON(ctx)
	if err != nil {
		log.Fatalf("get Inner object %v: %v", attachedObject, err)
	}

	copiedAttached, err := innerCopy.JSON(ctx)
	if err != nil {
		log.Fatalf("get Inner object %v: %v", attachedObject, err)
	}

	if originalAttached != copiedAttached {
		t.Fatalf("Inner object not equal original object:  %v != %v", originalAttached, copiedAttached)
	}
}

func TestSetGetOneFromObject(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	object, err := db.Object(context.TODO(), "User")
	if err != nil {
		t.Fatal(err)
	}

	err = object.Set(context.TODO(), "Name", "Max", false)
	if err != nil {
		t.Fatal(err)
	}

	value, err := object.Get(context.TODO(), "Name")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Max" {
		t.Fatalf("Wrong value [%s] wanted [%s]\n", value, "Max")
	}

	/// CAR

	car, err := db.Object(context.TODO(), "Car")
	if err != nil {
		t.Fatal(err)
	}

	err = car.Set(context.TODO(), "Name", "MyCar", false)
	if err != nil {
		t.Fatal(err)
	}

	value, err = car.Get(context.TODO(), "Name")
	if err != nil {
		t.Fatal(err)
	}

	if value != "MyCar" {
		t.Fatal("Wrong value")
	}

	/// WHEEL

	wheel, err := car.Object(context.TODO(), "Wheel")
	if err != nil {
		t.Fatal(err)
	}

	err = wheel.Set(context.TODO(), "Color", "Black", false)
	if err != nil {
		t.Fatal(err)
	}

	value, err = wheel.Get(context.TODO(), "Color")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	/// TRAILER

	trailer, err := car.Object(context.TODO(), "Trailer")
	if err != nil {
		t.Fatal(err)
	}

	err = trailer.Set(context.TODO(), "Color", "Red", false)
	if err != nil {
		t.Fatal(err)
	}

	/// TEST WHEEL & TRAILER AREAS ARE STILL WORKING

	value, err = wheel.Get(context.TODO(), "Color")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	value, err = trailer.Get(context.TODO(), "Color")
	if err != nil {
		t.Fatal(err)
	}

	if value != "Red" {
		t.Fatal("Wrong value")
	}
}

// TestIsObject
func TestIsObject(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("%d", time.Now().Unix())
	object, err := db.Object(context.TODO(), name)
	if err != nil {
		log.Fatalf("Object error: %s\n", err)
	}

	time.Sleep(time.Second)

	ctx := context.TODO()
	if ok, err := db.IsObject(ctx, object.GetName()); err != nil || !ok {
		t.Fatal("Not an object, but should be")
	}

	name = fmt.Sprintf("%d", time.Now().Unix())
	newObject, err := object.Object(ctx, object.GetName())
	if err != nil {
		log.Fatalf("Object error: %s\n", err)
	}

	time.Sleep(time.Second)

	if ok, err := db.IsObject(ctx, newObject.GetName()); err != nil || !ok {
		t.Fatal("Not an object, but should be")
	}

	if ok, _ := db.IsObject(ctx, fmt.Sprintf("%d", time.Now().Unix())); ok {
		t.Fatal("Object, but should not be")
	}
}

func TestSize(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	object, err := db.Object(context.TODO(), fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		t.Fatal(err)
	}

	size, err := object.Size(context.TODO())
	if err != nil {
		t.Fatalf("Error %v\n", err)
	} else if size != 0 {
		t.Fatalf("Wrong size %d\n", size)
	}

	for i := 0; i < 100; i++ {
		k := fmt.Sprint(i)
		v := fmt.Sprint(i)
		if err = object.Set(context.TODO(), k, v, false); err != nil {
			t.Fatalf("Set error %v\n", err)
		}
	}

	if size, err := object.Size(context.TODO()); err != nil || size != 100 {
		t.Fatalf("Wrong size %d\n", size)
	}
}

// TestGetObject
func TestGetObject(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	name := fmt.Sprintf("%d", time.Now().Unix())
	object, err := db.Object(context.TODO(), name)
	if err != nil {
		log.Fatalf("Object error: %s\n", err)
	}

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	for k, v := range data {
		if err = object.Set(context.TODO(), k, v, false); err != nil {
			t.Fatalf("Set error %v\n", err)
		}
	}

	ctx := context.TODO()
	m, err := object.JSON(ctx)
	if err != nil {
		t.Fatalf("GetObject error %v\n", err)
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal error %v\n", err)
	}

	if m != string(dataJSON) {
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

//func TestClient_StructToObject(t *testing.T) {
//	db, err := itisadb.New(_ctx,":8888")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	type Wheel struct {
//		Color string
//		Size  string
//	}
//
//	type Trailer struct {
//		Color string
//		Size  string
//	}
//
//	type Car struct {
//		Name    string
//		Wheel   *Wheel
//		Trailer *Trailer
//	}
//
//	type IQ struct {
//		Count int
//	}
//
//	var t1 = "qwe"
//	var t2 = &t1
//	var t3 = &t2
//
//	type User struct {
//		Name  string
//		Age   int
//		Email string
//		Car   *Car
//		IQ    IQ
//		T     ***string
//	}
//
//	user := User{
//		Name:  "Max",
//		Age:   25,
//		Email: "max@mail.ru",
//		Car: &Car{
//			Name: "MyCar",
//			Wheel: &Wheel{
//				Color: "Black",
//				Size:  "Big",
//			},
//			Trailer: &Trailer{
//				Color: "Red",
//				Size:  "Big",
//			},
//		},
//		IQ: IQ{
//			Count: 1,
//		},
//		T: &t3,
//	}
//
//	object, err := db.StructToObject(context.TODO(), fmt.Sprintf("User%d", 1), user)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	mu := map[string]string{
//		"Name":  "Max",
//		"Car":   "\n\tName: MyCar\n\tWheel: \n\t\tSize: Big\n\t\tColor: Black\n\tTrailer: \n\t\t\tColor: Red\n\t\t\tSize: Big",
//		"Age":   "25",
//		"Email": "max@mail.ru",
//		"IQ":    "\n\t\tCount: 1",
//		"T":     "qwe",
//	}
//
//	mi, err := object.JSON(context.TODO())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	splittedMi := strings.Split(strings.Replace(strings.Replace(mi["Car"], "\t", "", -1), "\n", "", -1), "")
//	splittedMu := strings.Split(strings.Replace(strings.Replace(mu["Car"], "\t", "", -1), "\n", "", -1), "")
//
//	if !IsTheSameArray(splittedMi, splittedMu) {
//		t.Fatalf("want\n%v,\ngot \n%v", splittedMu, splittedMi)
//	}
//
//	delete(mu, "Car")
//	delete(mi, "Car")
//
//	if !reflect.DeepEqual(mi, mu) {
//		t.Fatalf("want\n%v,\ngot \n%v", mu, mi)
//	}
//}

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

func TestClient_ObjectToStruct(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		t.Fatal(err)
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
		Wheel   Wheel
		Trailer *Trailer
	}

	type User struct {
		Name  string
		Name2 *string
		Age   int
		Age2  int32
		Age3  int8
		Email string
		Car   *Car
	}

	name := "max2"
	user := User{
		Name:  "Max",
		Name2: &name,
		Age:   25,
		Age2:  25,
		Age3:  25,
		Email: "max@mail.ru",
		Car: &Car{
			Name: "MyCar",
			Wheel: Wheel{
				Color: "Black",
				Size:  "Big",
			},
			Trailer: &Trailer{
				Color: "Red",
				Size:  "Big",
			},
		},
	}

	_, err = db.StructToObject(context.TODO(), "User33", user)
	if err != nil {
		t.Fatal(err)
	}

	us := &User{}
	err = db.ObjectToStruct(context.TODO(), "User33", us)
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
			if fa.Type().Elem().Kind() != reflect.Struct {
				fa, fb = fa.Elem(), fb.Elem()
				if fa.Interface() != fb.Interface() {
					fmt.Printf("[%v] [%v]\n", fa.Interface(), fb.Interface())
					return false
				}
			} else if !cmpReflect(fa.Elem(), fb.Elem()) {
				return false
			}
		} else if fa.Interface() != fb.Interface() {
			fmt.Println(a.Field(i).Interface(), b.Field(i).Interface())
			return false
		}
	}
	return true
}

func TestClient_GetCmp(t *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	err = db.SetOne(ctx, "qwe", "123", false)
	if err != nil {
		t.Errorf("SetTo() error = %v, wantErr no", err)
		return
	}

	got, err := itisadb.GetCmp[string](ctx, db, "qwe")
	if err != nil {
		t.Errorf("GetCmp() error = %v, wantErr no", err)
		return
	}

	if got != "123" {
		t.Errorf("got != want\n%v!=%v", got, "123")
		return
	}

	iint, err := itisadb.GetCmp[int](ctx, db, "qwe")
	if err != nil {
		t.Errorf("GetCmp() error = %v, wantErr no", err)
		return
	}

	if iint != 123 {
		t.Errorf("got != want\n%v!=%v", iint, 123)
		return
	}
}

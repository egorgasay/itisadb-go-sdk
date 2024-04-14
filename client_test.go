package itisadb_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/egorgasay/itisadb-go-sdk"
)

var _ctx = context.TODO()

// TestSetGetOne to run this test, itisadb must be run on :8888.
func TestSetGetOne(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	r := db.SetOne(ctx, "qwe", "111")
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	r = db.SetOne(ctx, "qwe", "111")
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	r = db.SetOne(ctx, "qwe", "111")
	if r.Error() != nil {
		t.Fatal(r.Error())
	}
}

// TestSetToAllGet to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetToAllGet(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()
	ctx := context.TODO()
	db.SetToAll(ctx, "all_key", "qqq").Unwrap()
	get := db.GetOne(ctx, "all_key").Unwrap()

	want := "qqq"
	if get.Value != want {
		t.Fatalf("Wrong value [%s] wanted [%s]\n", get.Value, want)
	}
}

// TestSetManyGetMany to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetManyGetMany(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	m := map[string]string{
		"m1": "k1",
		"m2": "k2",
		"m3": "k3",
		"m4": "k4",
		"m5": "k5",
	}

	ctx := context.TODO()
	db.SetMany(ctx, m).Unwrap()

	get := db.GetOne(ctx, "m2").Unwrap()

	if get.Value != "k2" {
		t.Fatal("Wrong value")
	}

	k := []string{"m1", "m2", "m3", "m4", "m5"}
	res := db.GetMany(ctx, k).Unwrap()

	if !reflect.DeepEqual(res, m) {
		t.Fatal("Wrong value")
	}
}

// TestSetManyOptsGetManyOpts to run this test, itisadb must be run on :8888.
// TODO: Add edge cases.
func TestSetManyOptsGetManyOpts(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	me := map[string]string{
		"mo1": "k1",
		"mo2": "k2",
		"mo3": "k3",
		"mo4": "k4",
		"mo5": "k5",
	}

	m := map[string]itisadb.ValueSpec{
		"mo1": {Value: "k1", Options: itisadb.SetOptions{Server: 1}},
		"mo2": {Value: "k2", Options: itisadb.SetOptions{Server: 1}},
		"mo3": {Value: "k3", Options: itisadb.SetOptions{Server: -1}},
		"mo4": {Value: "k4", Options: itisadb.SetOptions{Server: -2}},
		"mo5": {Value: "k5", Options: itisadb.SetOptions{Server: -3}},
	}

	ctx := context.TODO()
	db.SetManyOpts(ctx, m).Unwrap()

	get := db.GetOne(ctx, "mo2").Unwrap()

	if get.Value != "k2" {
		t.Fatal("Wrong value")
	}

	k := []itisadb.KeySpec{
		{Key: "mo1", Options: itisadb.GetOptions{Server: 1}},
		{Key: "mo2", Options: itisadb.GetOptions{Server: 1}},
		{Key: "mo3", Options: itisadb.GetOptions{Server: -1}},
		{Key: "mo4", Options: itisadb.GetOptions{Server: 0}},
		{Key: "mo5", Options: itisadb.GetOptions{Server: 0}},
	}

	res := db.GetManyOpts(ctx, k).Unwrap()

	if !reflect.DeepEqual(res, me) {
		t.Fatal("Wrong value")
	}
}

func TestDelete(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()
	ctx := context.TODO()

	num := rand.Int31()
	n := fmt.Sprint(num)
	err := db.SetOne(ctx, "key_for_delete"+n, "value_for_delete").Error()
	if err != nil {
		t.Fatal(err)
	}

	err = db.DelOne(ctx, "key_for_delete"+n).Error()
	if err != nil {
		t.Fatal(err)
	}

	err = db.GetOne(ctx, "key_for_delete"+n).Error()
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Key should be deleted, but %v", err)
	}
}

func TestDeleteObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteObject" + n

	indx := db.Object(name)

	_ = indx.Set(ctx, "key_for_delete", "value_for_delete").Unwrap()

	_ = indx.Get(ctx, "key_for_delete").Unwrap()

	indx.DeleteObject(ctx).Unwrap()

	indx = db.Object(name)

	err := indx.Get(ctx, "key_for_delete").Error()
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Object should be deleted")
	}

	// TEST DELETE INNER OBJECT

	name = "inner_object"
	inner := indx.Object(name)

	name = "more_inner_object"
	moreInner := inner.Object(name)

	moreInner.DeleteObject(ctx).Unwrap()

	err = inner.Set(ctx, "key_for_delete", "value_for_delete").Error()
	if err != nil {
		t.Fatalf("Inner object %v: %v", name, err)
	}

	inner.Get(ctx, "key_for_delete").Unwrap()

	inner.DeleteObject(ctx).Unwrap()

	inner = indx.Object(name)

	err = inner.Get(ctx, "key_for_delete").Error()
}

func TestDeleteAttr(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()
	num := rand.Int31()
	n := fmt.Sprint(num)
	name := "TestDeleteAttr" + n

	obj := db.Object(name)

	obj.Set(ctx, "key_for_delete", "value_for_delete").Unwrap()

	obj.Get(ctx, "key_for_delete").Unwrap()

	obj.DeleteKey(ctx, "key_for_delete").Unwrap()

	obj = db.Object(name)

	err := obj.Get(ctx, "key_for_delete").Error()
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatal("Key should be deleted")
	}

	// TEST DELETE ATTR INNER OBJECT

	name = "inner_object"
	inner := obj.Object(name)

	inner.Set(ctx, "key_for_delete", "value_for_delete").Unwrap()

	inner.Get(ctx, "key_for_delete").Unwrap()

	inner.DeleteKey(ctx, "key_for_delete").Unwrap()

	inner = obj.Object(name)

	err = inner.Get(ctx, "key_for_delete").Error()

	// TEST DELETE ATTR (INNER OBJECT) KEY DOES NOT EXIST

	err = inner.DeleteKey(ctx, "key_for_delete_does_not_exist").Error()
	if !errors.Is(err, itisadb.ErrNotFound) {
		t.Fatalf("Inner Object key shouldn't be deleted: %v", err)
	}

}

func TestAttachObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	originalObject := "TestAttachObject"
	indx := db.Object(originalObject)
	indx.Create(ctx).Unwrap()

	attachedObject := "TestAttachObject2"
	inner := db.Object(attachedObject)
	inner.Create(ctx).Unwrap()

	if err := indx.DeleteObject(ctx).Error(); err != nil {
		t.Fatalf("Delete object %v: %v", originalObject, err)
	}

	if err := inner.DeleteObject(ctx).Error(); err != nil {
		t.Fatalf("Delete object %v: %v", originalObject, err)
	}

	originalObject = "TestAttachObject"
	indx = db.Object(originalObject)
	indx.Create(ctx).Unwrap()

	attachedObject = "TestAttachObject2"
	inner = db.Object(attachedObject)
	inner.Create(ctx).Unwrap()

	inner.Set(ctx, "key_for_attach", "value_for_attach").Unwrap()

	indx.Attach(ctx, inner.Name()).Unwrap()

	innerCopy := indx.Object(attachedObject)
	innerCopy.Create(ctx).Unwrap()

	innerCopy.Set(ctx, "key_for_attach3", "value_for_attach3").Unwrap()
	_ = inner.Set(ctx, "key_for_attach4", "value_for_attach4").Unwrap()

	originalAttached := inner.JSON(ctx).Unwrap()

	copiedAttached := innerCopy.JSON(ctx).Unwrap()

	if !cmpJSON(originalAttached, copiedAttached) {
		t.Fatalf("Inner object not equal original object:  %v != %v", originalAttached, copiedAttached)
	}
}

func TestSetGetOneFromObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	object := db.Object("User")
	object.Create(context.TODO()).Unwrap()

	object.Set(context.TODO(), "Name", "Max").Unwrap()

	value := object.Get(context.TODO(), "Name").Unwrap()

	if value != "Max" {
		t.Fatalf("Wrong value [%s] wanted [%s]\n", value, "Max")
	}

	/// CAR

	car := db.Object("Car")
	car.Create(context.TODO()).Unwrap()

	car.Set(context.TODO(), "Name", "MyCar").Unwrap()

	value = car.Get(context.TODO(), "Name").Unwrap()

	if value != "MyCar" {
		t.Fatal("Wrong value")
	}

	/// WHEEL

	wheel := car.Object("Wheel").Create(context.TODO()).Unwrap()

	_ = wheel.Set(context.TODO(), "Color", "Black").Unwrap()

	value = wheel.Get(context.TODO(), "Color").Unwrap()

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	/// TRAILER

	trailer := car.Object("Trailer").Create(context.TODO()).Unwrap()

	trailer.Set(context.TODO(), "Color", "Red").Unwrap()

	/// TEST WHEEL & TRAILER AREAS ARE STILL WORKING

	value = wheel.Get(context.TODO(), "Color").Unwrap()

	if value != "Black" {
		t.Fatal("Wrong value")
	}

	value = trailer.Get(context.TODO(), "Color").Unwrap()

	if value != "Red" {
		t.Fatal("Wrong value")
	}
}

// TestIsObject
func TestIsObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	name := fmt.Sprintf("%d", time.Now().Unix())
	object := db.Object(name).Create(context.TODO()).Unwrap()

	time.Sleep(time.Second)

	ctx := context.TODO()
	if ok := object.Is(ctx).Unwrap(); !ok {
		t.Fatal("Not an object, but should be")
	}

	name = fmt.Sprintf("%d", time.Now().Unix())
	newObject := object.Object(object.Name()).Create(context.TODO()).Unwrap()

	time.Sleep(time.Second)

	if ok := newObject.Is(ctx).Unwrap(); !ok {
		t.Fatal("Not an object, but should be")
	}

	if ok := object.Object(fmt.Sprintf("%d", time.Now().Unix())).Is(ctx).Unwrap(); ok {
		t.Fatal("Object, but should not be")
	}
}

func TestSize(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	object := db.Object(fmt.Sprintf("%d", time.Now().Unix())).Create(context.TODO()).Unwrap()

	size := object.Size(context.TODO()).Unwrap()
	if size != 0 {
		t.Fatalf("Wrong size %d\n", size)
	}

	for i := 0; i < 100; i++ {
		k := fmt.Sprint(i)
		v := fmt.Sprint(i)
		object.Set(context.TODO(), k, v).Unwrap()
	}

	if size := object.Size(context.TODO()).Unwrap(); size != 100 {
		t.Fatalf("Wrong size %d\n", size)
	}
}

// TestGetObject
func TestGetObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	time.Sleep(time.Second)

	name := fmt.Sprintf("%d", time.Now().Unix())
	object := db.Object(name).Create(context.TODO()).Unwrap()

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	for k, v := range data {
		object.Set(context.TODO(), k, v).Unwrap()
	}

	ctx := context.TODO()
	m := object.JSON(ctx).Unwrap()

	var mMap map[string]any
	err := json.Unmarshal([]byte(m), &mMap)
	if err != nil {
		t.Fatalf("Marshal error %v\n", err)
	}

	want := `{
	"name": "{{ name }}",
	"values": [
		{
			"name": "key3",
			"value": "value3"
		},
		{
			"name": "key2",
			"value": "value2"
		},
		{
			"name": "key4",
			"value": "value4"
		},
		{
			"name": "key5",
			"value": "value5"
		},
		{
			"name": "key1",
			"value": "value1"
		}
	]
}`

	want = strings.Replace(want, "{{ name }}", name, -1)

	if !cmpJSON(want, m) {
		t.Fatalf("Want %s, got %s\n", want, m)
	}
}

//func TestDistinct(t *testing.T) {
//	f:= os.Open("/tmp/log14/transactionLogger")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var keys = make(map[string]struct{}, 16000)
//
//	scanner := bufio.NewScanner(f)
//	for scanner.Scan() {
//		action := scanner.Text()
//		decode:= strutil.Base64Decode([]byte(action))
//		if err != nil {
//			return
//		}
//
//		split := strings.Split(string(decode), " ")
//		key := split[1]
//		keys[key] = struct{}{}
//	}
//
//	t.Log(len(keys))
//}

func TestClient_StructToObject(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

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

	type IQ struct {
		Count int
	}

	var t1 = "qwe"
	var t2 = &t1
	var t3 = &t2

	type User struct {
		Name  string
		Age   int
		Email string
		Car   *Car
		IQ    IQ
		T     ***string
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
		IQ: IQ{
			Count: 1,
		},
		T: &t3,
	}

	object, err := db.StructToObject(context.TODO(), fmt.Sprintf("User%d", 1), user)
	if err != nil {
		t.Fatal(err)
	}

	mi := object.JSON(context.TODO()).Unwrap()

	want := `{
	"name": "User1",
	"values": [
		{
			"name": "Name",
			"value": "Max"
		},
		{
			"name": "IQ",
			"values": [
				{
					"name": "Count",
					"value": "1"
				}
			]
		},
		{
			"name": "T",
			"value": "qwe"
		},
		{
			"name": "Age",
			"value": "25"
		},
		{
			"name": "Email",
			"value": "max@mail.ru"
		},
		{
			"name": "Car",
			"values": [
				{
					"name": "Wheel",
					"values": [
						{
							"name": "Color",
							"value": "Black"
						},
						{
							"name": "Size",
							"value": "Big"
						}
					]
				},
				{
					"name": "Trailer",
					"values": [
						{
							"name": "Color",
							"value": "Red"
						},
						{
							"name": "Size",
							"value": "Big"
						}
					]
				},
				{
					"name": "Name",
					"value": "MyCar"
				}
			]
		}
	]
}`

	if !cmpJSON(mi, want) {
		t.Fatalf("Want %s, got %s\n", want, mi)
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

func TestClient_ObjectToStruct(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

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

	_, err := db.StructToObject(context.TODO(), "User33", user)
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

//func TestClient_GetCmp(t *testing.T) {
//	db:= itisadb.New(_ctx, ":8888")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	ctx := context.Background()
//
//	err = db.SetOne(ctx, "qwe", "123").Error()
//	if err != nil {
//		t.Errorf("SetTo() error = %v, wantErr no", err)
//		return
//	}
//
//	got:= itisadb.GetCmp[string](ctx, db, "qwe")
//	if err != nil {
//		t.Errorf("GetCmp() error = %v, wantErr no", err)
//		return
//	}
//
//	if got != "123" {
//		t.Errorf("got != want\n%v!=%v", got, "123")
//		return
//	}
//
//	iint:= itisadb.GetCmp[int](ctx, db, "qwe")
//	if err != nil {
//		t.Errorf("GetCmp() error = %v, wantErr no", err)
//		return
//	}
//
//	if iint != 123 {
//		t.Errorf("got != want\n%v!=%v", iint, 123)
//		return
//	}
//}

func cmpJSON(want, got string) bool {
	var m1 = make(map[rune]int)
	var m2 = make(map[rune]int)

	for _, v := range want {
		if v == ' ' || v == '\n' || v == '\t' {
			continue
		}
		m1[v]++
	}
	for _, v := range got {
		if v == ' ' || v == '\n' || v == '\t' {
			continue
		}
		m2[v]++
	}

	return reflect.DeepEqual(m1, m2)
}

func TestClient_DeleteUser(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	err := db.NewUser(ctx, "max", "123").Error()
	if err != nil {
		t.Fatalf("NewUser: %v", err)
		return
	}

	res := db.DeleteUser(ctx, "max")
	if res.Unwrap() != true {
		t.Fatalf("user should exist")
		return
	}

	res = db.DeleteUser(ctx, "max2")
	if res.Unwrap() != false {
		t.Fatalf("user should not exist")
	}
}

func TestClient_NewUser(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	err := db.DeleteUser(ctx, "max").Error()
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
		return
	}

	err = db.NewUser(ctx, "max", "123").Error()
	if err != nil {
		t.Fatalf("NewUser: %v", err)
		return
	}

	conf := itisadb.Config{
		Credentials: itisadb.Credentials{
			Login: "max", Password: "123",
		},
	}

	db = itisadb.New(_ctx, ":8888", conf).Unwrap()
	if err != nil {
		t.Fatalf("New with config: %v", err)
		return
	}

	res := db.NewUser(ctx, "max2", "123", itisadb.NewUserOptions{
		Level: itisadb.DefaultLevel,
	})

	if res.Error() != nil {
		t.Fatalf("NewUser with options: %v", res.Error())
		return
	}

	//if res.Unwrap() != true {
	//	t.Fatalf("user should exist")
	//	return
	//}
}

func TestClient_ChangePassword(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	err := db.NewUser(ctx, "max", "123").Error()
	if err != nil {
		t.Fatalf("NewUser: %v", err)
		return
	}

	err = db.ChangePassword(ctx, "max", "1234").Error()
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
		return
	}

	db = itisadb.New(_ctx, ":8888", itisadb.Config{
		Credentials: itisadb.Credentials{
			Login: "max", Password: "1234",
		},
	}).Unwrap()

	err = itisadb.New(_ctx, ":8888", itisadb.Config{
		Credentials: itisadb.Credentials{
			Login: "max", Password: "63231fwe23e1e3",
		},
	}).Error()
	if err == nil {
		t.Fatal("no error with wrong password")
		return
	}
}

func TestClient_ChangeLevel(t *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	ctx := context.TODO()

	err := db.NewUser(ctx, "max", "123", itisadb.NewUserOptions{
		Level: itisadb.DefaultLevel,
	}).Error()
	if err != nil {
		t.Fatalf("NewUser: %v", err)
		return
	}

	// Create "max" with default level
	_ = db.Object("max").Create(ctx, itisadb.ObjectOptions{
		Level: itisadb.DefaultLevel,
	}).Unwrap()

	// Change "max" to restricted level
	err = db.ChangeLevel(ctx, "max", itisadb.RestrictedLevel).Error()
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
		return
	}

	_ = db.Object("max").Create(ctx, itisadb.ObjectOptions{
		Level: itisadb.DefaultLevel,
	}).Unwrap()
}

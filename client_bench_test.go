package itisadb_test

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"sync"
	"testing"
	"time"
)

// Test to define rps for SetOne.
func TestSetOneRPS(b *testing.T) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	const maxRPS = 100_000
	const gnum = maxRPS * 10

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		ints[i] = fmt.Sprint(i)
	}

	b.Log("Hops:", gnum/maxRPS)

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
				err := db.SetOne(ctx, ints[i], "qdw").Error()
				if err != nil {
					b.Log(err)
				}
				wgSent.Done()
			}(i)
		}
		wg.Wait()
		start := time.Now()
		wgSent.Wait()
		total += time.Since(start)
		b.Log(time.Since(start))
	}
	b.Logf("Total time: %v RPS: %v Ahead:%v", total, maxRPS, (gnum/maxRPS)*time.Second-total)
}

// Test to define rps for Get.
func BenchmarkGetOneRPS(b *testing.B) {
	db := itisadb.New(_ctx, ":8888").Unwrap()
	const maxRPS = 71_100
	const gnum = maxRPS * 10

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		ints[i] = fmt.Sprint(i)
	}

	b.Log("Hops:", gnum/maxRPS)

	var total time.Duration
	for tt := gnum / maxRPS; tt > 0; tt-- {
		var wg sync.WaitGroup
		wg.Add(maxRPS)

		var wgSent sync.WaitGroup
		wgSent.Add(maxRPS)
		ctx := context.TODO()
		for i := 0; i < maxRPS; i++ {
			go func(i int) {
				wg.Wait()
				err := db.GetOne(ctx, ints[i]).Error()
				if err != nil {
					b.Log(err)
				}
				wgSent.Done()
			}(i)
			wg.Done()
		}
		wg.Wait()
		start := time.Now()
		wgSent.Wait()
		total += time.Since(start)
		b.Log(time.Since(start))
	}
	b.Logf("Total time: %v RPS: %v OK:%v", total, maxRPS, total/time.Duration(maxRPS))
}

// Test to define rps for Get.
func BenchmarkSetToObjectRPS(b *testing.B) {
	db := itisadb.New(_ctx, ":800").Unwrap()

	ctx := context.Background()
	ind := db.Object(ctx, "User1").Unwrap()

	start := time.Now()
	mail := ind.Get(ctx, "Email").Unwrap()
	b.Log(time.Since(start))

	if mail != "max@mail.ru" {
		b.Error("Wrong email")
	}
}

// Test to define rps for Get.
func BenchmarkGetFromObjectRPS2(b *testing.B) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	const gnum = 1500000
	const maxRPS = 80_000

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	ctx := context.TODO()
	res := db.Object(ctx, "User")
	if res.IsErr() {
		b.Log(res.Error(), "[User]")
		return
	}
	ind := res.Unwrap()

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		x := "User" + fmt.Sprint(i)
		err := ind.Set(ctx, x, "xx").Error()
		if err != nil {
			b.Log(err)
			return
		}

		ints[i] = x
	}

	b.Log("Hops:", gnum/maxRPS)

	var total time.Duration
	for tt := gnum / maxRPS; tt > 0; tt-- {
		var wg sync.WaitGroup
		wg.Add(maxRPS)

		var wgSent sync.WaitGroup
		wgSent.Add(maxRPS)
		for i := 0; i < maxRPS; i++ {
			wg.Done()
			go func(i int) {
				wg.Wait()
				err := ind.Get(ctx, ints[i]).Error()
				if err != nil {
					b.Log(err)
					return
				}
				wgSent.Done()
			}(i)

		}
		wg.Wait()
		start := time.Now()
		wgSent.Wait()
		total += time.Since(start)
		b.Log(time.Since(start))
	}
	b.Log(total)
}

func BenchmarkSetToObjectRPS2(b *testing.B) {
	db := itisadb.New(_ctx, ":8888").Unwrap()

	const gnum = 1500000
	const maxRPS = 80_000
	const value = "xx"

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	ctx := context.TODO()
	ind := db.Object(ctx, "User").Unwrap()

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		x := "User" + fmt.Sprint(i)
		ints[i] = x
	}

	b.Log("Hops:", gnum/maxRPS)

	var total time.Duration
	for tt := gnum / maxRPS; tt > 0; tt-- {
		var wg sync.WaitGroup
		wg.Add(maxRPS)

		var wgSent sync.WaitGroup
		wgSent.Add(maxRPS)
		for i := 0; i < maxRPS; i++ {
			wg.Done()
			go func(i int) {
				wg.Wait()
				_ = ind.Set(ctx, ints[i], value).Unwrap()
				wgSent.Done()
			}(i)

		}
		wg.Wait()
		start := time.Now()
		wgSent.Wait()
		total += time.Since(start)
		b.Log(time.Since(start))
	}
	b.Log(total)
}

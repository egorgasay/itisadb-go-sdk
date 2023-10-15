package itisadb_test

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
	"sync"
	"testing"
	"time"
)

// Test to define rps for SetOne.
func TestSetOneRPS(b *testing.T) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		log.Fatalln(err)
	}
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
				err := db.SetOne(ctx, ints[i], "qdw").Err()
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
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		log.Fatalln(err)
	}
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
				_, err := db.GetOne(ctx, ints[i]).ValueAndErr()
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
	db, err := itisadb.New(_ctx, ":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	ind, err := db.Object(ctx, "User1").ValueAndErr()
	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()
	mail, err := ind.Get(ctx, "Email").ValueAndErr()
	b.Log(time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}

	if mail != "max@mail.ru" {
		b.Error("Wrong email")
	}
}

// Test to define rps for Get.
func BenchmarkGetFromObjectRPS2(b *testing.B) {
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		b.Fatal(err)
	}
	const gnum = 1500000
	const maxRPS = 80_000

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	ctx := context.TODO()
	ind, err := db.Object(ctx, "User").ValueAndErr()
	if err != nil {
		b.Log(err, "[User]")
		return
	}

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		x := "User" + fmt.Sprint(i)
		err := ind.Set(ctx, x, "xx").Err()
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
				_, err = ind.Get(ctx, ints[i]).ValueAndErr()
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
	db, err := itisadb.New(_ctx, ":8888")
	if err != nil {
		b.Fatal(err)
	}
	const gnum = 1500000
	const maxRPS = 80_000
	const value = "xx"

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	ctx := context.TODO()
	ind, err := db.Object(ctx, "User").ValueAndErr()
	if err != nil {
		b.Log(err, "[User]")
		return
	}

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
				_, err = ind.Set(ctx, ints[i], value).ValueAndErr()
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

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
func BenchmarkSetOneRPS(b *testing.B) {
	db, err := itisadb.New(":8888")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 10000

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
				err := db.SetOne(ctx, ints[i], "qdw", false)
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
	b.Log(total)
}

// Test to define rps for Get.
func BenchmarkGetOneRPS(b *testing.B) {
	db, err := itisadb.New(":8888")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 10000

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
				_, err := db.Get(ctx, ints[i])
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
	b.Log(total)
}

// Test to define rps for Get.
func BenchmarkGetFromDiskObjectRPS(b *testing.B) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	ind, err := db.Object(ctx, "User1")
	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()
	mail, err := ind.Get(ctx, "Email")
	b.Log(time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}

	if mail != "max@mail.ru" {
		b.Error("Wrong email")
	}
}

// Test to define rps for Get.
func BenchmarkGetFromDiskObjectRPS2(b *testing.B) {
	db, err := itisadb.New(":8888")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 10000

	b.Log("Total actions:", gnum)
	b.Log("RPS:", maxRPS)

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		ints[i] = "User" + fmt.Sprint(i)
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
				ind, err := db.Object(ctx, ints[i])
				if err != nil {
					b.Log(err, "["+ints[i]+"]")
					return
				}

				_, err = ind.Get(ctx, "Email")
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

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
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 40000

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
				db.SetOne(ctx, ints[i], "qdw", false)
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

// Test to define rps for GetOne.
func BenchmarkGetOneRPS(b *testing.B) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 20000

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
		b.Log(time.Since(start))
	}
	b.Log(total)
}

// Test to define rps for GetOne.
func BenchmarkGetFromDiskIndexRPS(b *testing.B) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	ind, err := db.Index(ctx, "User1")
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

// Test to define rps for GetOne.
func BenchmarkGetFromDiskIndexRPS2(b *testing.B) {
	db, err := itisadb.New(":800")
	if err != nil {
		log.Fatalln(err)
	}
	const gnum = 1500000
	const maxRPS = 10000

	log.Println("Total actions:", gnum)
	log.Println("RPS:", maxRPS)

	var ints = make([]string, maxRPS)
	for i := 0; i < maxRPS; i++ {
		ints[i] = "User" + fmt.Sprint(i)
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
				ind, err := db.Index(ctx, ints[i])
				if err != nil {
					log.Println(err, "["+ints[i]+"]")
				}

				_, err = ind.Get(ctx, "Email")
				if err != nil {
					log.Println(err)
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

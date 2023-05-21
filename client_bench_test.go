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

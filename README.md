# `itisadb-go-sdk` - itisadb Driver and Toolkit [![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/mod)](https://pkg.go.dev/golang.org/x/mod)

### [itisadb](https://github.com/egorgasay/itisadb) is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database.

# [Documentation](https://pkg.go.dev/github.com/egorgasay/itisadb-go-sdk)


# Installation
```bash
go get github.com/egorgasay/itisadb-go-sdk
```

# Quick start
```go
package main

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk"
	"log"
)

// main to run this test, itisadb must be run on :8888.
func main() {
	ctx := context.TODO()

	db := itisadb.New(ctx, ":8888").Unwrap()

	db.SetOne(ctx, "qwe", "111").Unwrap()

	if x := db.GetOne(ctx, "qwe").Unwrap().Value; x != "111" {
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", x)
	}
}
```

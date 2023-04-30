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

// main to run this test, itisadb must be run on :800.
func main() {
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
		log.Fatal("Wrong value")
	} else {
		log.Println("Value:", get)
	}
}
```

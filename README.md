# `grpcis-go-sdk` - gRPCis Driver and Toolkit [![PkgGoDev](https://pkg.go.dev/badge/github.com/egorgasay/grpcis-go-sdk](https://pkg.go.dev/github.com/redis/egorgasay/grpcis-go-sdk)

### [gRPCis](https://github.com/egorgasay/gRPCis) is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database.

# [Documentation](https://pkg.go.dev/github.com/egorgasay/grpcis-go-sdk)


# Installation  
```bash
go get https://github.com/egorgasay/grpcis-go-sdk
```

# Quick start  
```go
package main

import (
	"context"
	"github.com/egorgasay/grpcis-go-sdk"
	"log"
)

// main to run this test, gRPCis must be run on :800.
func main() {
	db, err := grpcis.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = db.SetOne(ctx, "qwe", "111")
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

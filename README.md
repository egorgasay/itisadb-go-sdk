# gRPCis client for Go [![PkgGoDev](https://pkg.go.dev/badge/github.com/redis/go-redis/v9)](https://pkg.go.dev/github.com/redis/go-redis/v9?tab=doc)

### [gRPCis](https://github.com/egorgasay/grpc-storage) is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database.

# [Documentation](https://pkg.go.dev/github.com/egorgasay/gRPCis-client)


# Installation  
```bash
go get https://github.com/egorgasay/gRPCis-client
```

# Quick start  
```go
package main

import (
	"context"
	grpcisclient "github.com/egorgasay/gRPCis-client"
	"log"
)

// main to run this test, gRPCis must be run on :800.
func main() {
	grpcis, err := grpcisclient.New(":800")
	if err != nil {
		return
	}

	ctx := context.TODO()
	err = grpcis.SetOne(ctx, "qwe", "111")
	if err != nil {
		log.Fatalln(err)
	}

	get, err := grpcis.GetOne(ctx, "qwe")
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

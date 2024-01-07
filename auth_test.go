package itisadb

import (
	"context"
	"reflect"
	"testing"

	api "github.com/egorgasay/itisadb-shared-proto/go"
	"google.golang.org/grpc/metadata"
)

func Test_withAuth(t *testing.T) {
	authMetadata = metadata.New(map[string]string{token: "XXX"})

	ctx := metadata.NewOutgoingContext(context.TODO(), authMetadata)

	got := withAuth(ctx)

	t.Log(got)
	if !reflect.DeepEqual(ctx, got) {
		t.Errorf("withAuth() = %v, want %v", got, ctx)
	}
}

func TestAuth(t *testing.T) {
	ctx := context.Background()
	cl := New(ctx, ":8888").Unwrap()

	// method that requires auth
	resp, err := cl.cl.Servers(withAuth(ctx), &api.ServersRequest{})
	if err != nil {
		t.Fatalf("Servers: %v", err)
	}

	t.Log("Auth OK")
	t.Log(resp.ServersInfo)
}

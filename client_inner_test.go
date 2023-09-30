package itisadb

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk/api"
	"testing"
)

func TestAuth(t *testing.T) {
	ctx := context.Background()
	cl, err := New(ctx, ":8888")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// method that requires auth
	resp, err := cl.cl.Servers(withAuth(ctx), &api.ServersRequest{})
	if err != nil {
		t.Fatalf("Servers: %v", err)
	}

	t.Log("Auth OK")
	t.Log(resp.ServersInfo)
}

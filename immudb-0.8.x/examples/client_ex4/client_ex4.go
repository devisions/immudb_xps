package main

import (
	"context"
	"fmt"
	"log"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
)

func main() {

	opts := immuclient.DefaultOptions()
	client, err := immuclient.NewImmuClient(opts)
	if err != nil {
		log.Fatalf(">>> Connection failed: %s\n", err)
	}
	ctx := context.Background()

	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`Immudb.1`))
	if err != nil {
		log.Fatalf(">>> Login failed: %s\n", err.Error())
	}
	defer func() { _ = client.Logout(ctx) }()

	fmt.Printf(">>> Login response:\n\t token=%s\n\t warning=%s\n", lr.Token, lr.Warning)

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	udr, err := client.UseDatabase(ctx, &schema.Database{Databasename: "lcsnapcyvvxl"})
	if err != nil {
		log.Fatal(">>> Failed to use the database. Reason:", err)
	}
	md.Set("authorization", udr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	idx := uint64(13)
	si, _ := client.ByIndex(ctx, idx)
	// Note here that the `StructuredValue`'s Value contains
	// a timestamp (.Timestamp) and the provided value (.Payload).
	log.Printf(">>> ByIndex(%d) => %+v\n", idx, si)

	root, _ := client.CurrentRoot(ctx)
	log.Printf(">>> CurrentRoot => %+v\n", root)

	fmt.Println()

}

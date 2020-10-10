package main

import (
	"context"
	"fmt"
	"log"

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

	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatalf(">>> Login failed: %s\n", err.Error())
	}
	fmt.Printf(">>> Login response:\n\t token=%s\n\t warning=%s\n", lr.Token, lr.Warning)

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	dlr, _ := client.DatabaseList(ctx)
	fmt.Printf(">>> List of databases:")
	for _, db := range dlr.Databases {
		fmt.Printf(" %s", db.Databasename)
	}
	fmt.Println()

	_ = client.Logout(ctx)

}

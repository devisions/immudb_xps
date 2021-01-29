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

	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatalf(">>> Login failed: %s\n", err.Error())
	}
	defer logout(ctx, client)

	fmt.Printf(">>> Login response:\n\t token=%s\n\t warning=%s\n", lr.Token, lr.Warning)

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	dbname := "mydb1"
	err = client.CreateDatabase(ctx, &schema.Database{Databasename: dbname})
	if err != nil {
		log.Fatalf(">>> Database creation failed: %s\n", err)
	} else {
		log.Printf(">>> Created new database named '%s'\n", dbname)
	}

}

func logout(ctx context.Context, client immuclient.ImmuClient) {
	err := client.Logout(ctx)
	if err != nil {
		log.Printf(">>> Logout failed: %s\n", err)
	}

}

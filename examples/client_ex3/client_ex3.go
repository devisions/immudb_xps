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
	defer func() { _ = client.Logout(ctx) }()

	fmt.Printf(">>> Login response:\n\t token=%s\n\t warning=%s\n", lr.Token, lr.Warning)

	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	idx1, err := client.Set(ctx, []byte("k1"), []byte("v1.0"))
	if err != nil {
		log.Println(">>> Failed to set k1,v1.0 Reason:", err)
	}
	log.Println("Set k1,v1.0 and got index", idx1.GetIndex())

	idx2, err := client.Set(ctx, []byte("k1"), []byte("v1.1"))
	if err != nil {
		log.Println(">>> Failed to set k1,v1.1 Reason:", err)
	}
	log.Println("Set k1,v1.1 and got index", idx2.GetIndex())

	si, _ := client.ByIndex(ctx, idx1.Index)
	log.Printf(">>> ByIndex(%d) => %+v\n", idx1.Index, si)

	root, _ := client.CurrentRoot(ctx)
	log.Printf(">>> CurrentRoot => %+v\n", root)

	fmt.Println()

}

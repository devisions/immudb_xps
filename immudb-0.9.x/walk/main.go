package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var cliArgs struct {
	host   string
	port   int
	user   string
	pass   string
	dbName string
}

func init() {
	flag.StringVar(&cliArgs.host, "host", "127.0.0.1", "hostname or IP adddress of immudb's listening endpoint")
	flag.IntVar(&cliArgs.port, "port", 3322, "port of immudb's listening endpoint")
	flag.StringVar(&cliArgs.user, "user", "immudb", "username to authenticate")
	flag.StringVar(&cliArgs.pass, "pass", "immudb", "password to authenticate")
	flag.StringVar(&cliArgs.dbName, "db", "defaultdb", "name of the database to use")
	flag.Parse()
}

func main() {

	client, err := client.NewImmuClient(
		client.DefaultOptions().
			WithAddress(cliArgs.host).
			WithPort(cliArgs.port),
	)
	if err != nil {
		log.Fatalf(">>> Failed to initialize the client. Reason: %v", err)
	}
	ctx := context.Background()

	lr, err := client.Login(ctx, []byte(cliArgs.user), []byte(cliArgs.pass))
	if err != nil {
		endNow(client, "Failed to login", err)
	}
	// Set up an authenticated context that will be used in all subsequent interactions.
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	udr, err := client.UseDatabase(ctx, &schema.Database{Databasename: cliArgs.dbName})
	if err != nil {
		endNow(client, "Failed to UseDatabase", err)
	}
	// Recollect the token that we get when using/switching the database.
	md = metadata.Pairs("authorization", udr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	txID := uint64(1)

	for {
		tx, err := client.TxByID(ctx, txID)
		if err != nil {
			e, ok := status.FromError(err)
			if ok {
				// This doesn't work since the error is "Reason: rpc error: code = Unknown desc = tx not found".
				// if e.Code() == codes.NotFound {
				// 	end(client)
				// }
				if strings.Contains(e.Message(), "tx not found") {
					end(client)
				}
			}
			endNow(client, fmt.Sprintf("Failed to TxByID with tx:%d", txID), err)
		}
		log.Printf("Tx ID: %d\n", txID)
		sc := *client.GetServiceClient()

		for _, txe := range tx.Entries {
			entry, err := sc.Get(ctx, &schema.KeyRequest{Key: txe.Key, AtTx: txID})
			if err != nil {
				log.Printf("Failed to get the value of key '%s'. Reason: %v\n", txe.Key, err)
				continue
			}
			log.Printf("TxEntry key: %s value: %s\n", txe.Key, entry.Value)
		}

		state, err := client.CurrentState(ctx)
		if err != nil {
			endNow(client, "Failed to CurrentState.", err)
		}
		if state.TxId == txID {
			end(client)
		}
		txID++
	}

}

func end(client client.ImmuClient) {
	if err := client.Disconnect(); err != nil {
		log.Printf(">>> Disconnect failed with error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func endNow(client client.ImmuClient, msg string, err error) {
	_ = client.Disconnect()
	errmsg := fmt.Sprintf(">>> %s Reason: %v", msg, err)
	log.Fatal(errmsg)
}

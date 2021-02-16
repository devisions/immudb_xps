package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	immustore "github.com/codenotary/immudb/embedded/store"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
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

	key := []byte("vsg")
	val := []byte("test-vsg")
	txmd, err := client.VerifiedSet(ctx, key, val)
	if err != nil {
		endNow(client, "Failed to Scan", err)
	}

	// Just print the details. Needed for a test between Go and Java ...

	fmt.Println("--{ schema.TxMetadata }--")
	fmt.Println("        ID:", txmd.Id)
	fmt.Println("   PrevAlh:", base64.StdEncoding.EncodeToString(txmd.PrevAlh[:]))
	fmt.Println("   PrevAlh:", txmd.PrevAlh)
	fmt.Println("        Ts:", txmd.Ts)
	fmt.Println("  NEntries:", txmd.Nentries)
	fmt.Println("        Eh:", base64.StdEncoding.EncodeToString(txmd.EH[:]))
	fmt.Println("        Eh:", txmd.EH)
	fmt.Println("    BlTxID:", txmd.BlTxId)
	fmt.Println("    BlRoot:", base64.StdEncoding.EncodeToString(txmd.BlRoot[:]))
	fmt.Println("    BlRoot:", txmd.BlRoot)

	itxmd := immustore.TxMetadata{
		ID:       txmd.Id,
		Ts:       txmd.Ts,
		NEntries: int(txmd.Nentries),
		BlTxID:   txmd.BlTxId,
	}
	copy(itxmd.PrevAlh[:], txmd.PrevAlh)
	copy(itxmd.Eh[:], txmd.EH)
	copy(itxmd.BlRoot[:], txmd.BlRoot)

	alh := itxmd.Alh()
	fmt.Println("--{ immustore.TxMetadata }--")

	fmt.Println("        ID:", itxmd.ID)
	fmt.Println("   PrevAlh:", base64.StdEncoding.EncodeToString(itxmd.PrevAlh[:]))
	fmt.Println("   PrevAlh:", itxmd.PrevAlh)
	fmt.Println("        Ts:", itxmd.Ts)
	fmt.Println("  NEntries:", itxmd.NEntries)
	fmt.Println("        Eh:", base64.StdEncoding.EncodeToString(itxmd.Eh[:]))
	fmt.Println("        Eh:", itxmd.Eh)
	fmt.Println("    BlTxID:", itxmd.BlTxID)
	fmt.Println("    BlRoot:", base64.StdEncoding.EncodeToString(itxmd.BlRoot[:]))
	fmt.Println("    BlRoot:", itxmd.BlRoot)
	fmt.Println("       Alh:", alh)
	fmt.Println("       Alh:", base64.StdEncoding.EncodeToString(alh[:]))

	end(client)

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

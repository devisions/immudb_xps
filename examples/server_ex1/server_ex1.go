package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
	immulogger "github.com/codenotary/immudb/pkg/logger"
	immuserver "github.com/codenotary/immudb/pkg/server"
	"google.golang.org/grpc/metadata"
)

func main() {

	log.Println(">>> Starting an embedded immudb server ...")

	const logfilename = "immuserver.log"
	flogger, file, err := immulogger.NewFileLogger("immuserver", logfilename)
	if err != nil {
		log.Fatalf(">>> Failed to init the logger: %s\n", err)
	} else {
		log.Printf(">>> Using '%s' file for logging.\n", logfilename)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalf(">>> Failed to close the log file: %s\n", err)
		}
	}()

	srvOpts := immuserver.DefaultOptions().WithLogfile(logfilename)
	srv := immuserver.DefaultServer()
	srv.WithOptions(srvOpts).WithLogger(flogger)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf(">>> Failed to start the server: %s\n", err)
		} else {
			log.Println(">>> Embedded immudb server started.")
		}
	}()
	defer func() {
		err := srv.Stop()
		// As this is for demo purposes, cleanup at the end.
		cleanup(srv.Options.Dir, srv.Options.Logfile)
		if err != nil {
			log.Fatalf(">>> Failed to stop the server: %s\n", err)
		}
	}()

	// wait for server to start
	time.Sleep(100 * time.Millisecond)

	log.Println(">>> Starting an immudb client ...")

	client, err := immuclient.NewImmuClient(immuclient.DefaultOptions())
	if err != nil {
		log.Fatalf(">>> Failed to start the immudb client: %s\n", err)
	}

	ctx, err := login(client)
	if err != nil {
		log.Fatalf(">>> Login failed: %s\n", err.Error())
	}
	log.Printf(">>> Client is connected? %t\n", client.IsConnected())

	root, err := client.CurrentRoot(ctx)
	if err != nil {
		log.Printf(">>> Failed to do CurrentRoot: %s\n", err)
	} else {
		log.Printf(">>> root: %s\n", root)
	}

}

func login(client immuclient.ImmuClient) (context.Context, error) {

	ctx := context.Background()
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		return ctx, err
	}
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return ctx, nil
}

func cleanup(dbDir string, logfile string) {

	fmt.Println(">>> Doing the cleanup now ...")
	os.RemoveAll(dbDir)
	os.Remove(logfile)
	// remove root
	files, err := filepath.Glob("./\\.root*")
	if err == nil {
		for _, f := range files {
			os.Remove(f)
		}
	}
}

package memdb_test

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
	"url-server/internal/service"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage/memdb"

	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	listenAddr = "127.0.0.1:8080"
)

func TestGetURL(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	defer finish()
	StartMyService(ctx, listenAddr)

	conn, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r1, _ := client.CreateShortURL(ctx, &pb.Request{Url: "http://welcome.com"})

	r, err := client.GetFullURL(ctx, &pb.Request{Url: r1.GetUrl()})
	if err != nil {
		t.Fatalf("could not get original url: %v", err)
	}

	result := &pb.Response{Url: "http://welcome.com"}

	if !reflect.DeepEqual(r.GetUrl(), result.GetUrl()) {
		t.Fatalf("logs2 dont match\nhave %+v\nwant %+v", r, result)
	}
}

func StartMyService(ctx context.Context, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := badger.DefaultOptions("")
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	defer db.Close()

	storage := memdb.NewStorage(db)
	service := service.NewURLServer(storage)

	server := grpc.NewServer()
	pb.RegisterURLServer(server, service)

	go func() {
		err = server.Serve(lis)
		if err != nil {
			fmt.Println("server Error")
		}

		<-ctx.Done()

		server.Stop()
		err = lis.Close()
		if err != nil {
			fmt.Println("close Error")
		}
	}()

	return nil
}

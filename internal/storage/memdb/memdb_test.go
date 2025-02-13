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

func TestURLFunctions(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	defer finish()
	err := StartMyService(ctx, listenAddr)
	if err != nil {
		t.Fatalf("fail to start server initial: %v", err)
	}

	conn, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	expected := []string{"http://welcome.com", "http://youtube.com/somevideo", "http://github.com/somerepo"}
	shortened := []string{}
	results := []string{}

	for i := range expected {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := client.CreateShortURL(ctx, &pb.Request{Url: expected[i]})
		if err != nil {
			t.Fatalf("could not make short url: %v", err)
		}
		shortened = append(shortened, r.GetUrl())
	}

	for i := range shortened {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := client.GetFullURL(ctx, &pb.Request{Url: shortened[i]})
		if err != nil {
			t.Fatalf("could not get original url: %v", err)
		}

		results = append(results, r.GetUrl())
	}

	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("original urls dont match\nhave %+v\nwant %+v", results, expected)
	}
}

func StartMyService(ctx context.Context, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := badger.DefaultOptions("/cmd")
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}

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
		db.Close()
		db.DropAll()
		err = lis.Close()
		if err != nil {
			fmt.Println("close Error")
		}
	}()

	return nil
}

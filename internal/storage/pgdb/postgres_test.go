package pgdb_test

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"url-server/internal/service"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage/pgdb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	listenAddr = "127.0.0.1:8080"
)

func TestCreateURL(t *testing.T) {
	conn, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//log.Fatalf("fail to dial: %v", err)
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.CreateShortURL(ctx, &pb.Request{Url: "http://welcome.com"})
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		t.Fatalf("could not make short url: %v", err)
	}
	//log.Printf("User: %s", r.GetName())
	fmt.Printf("Short URL: %s", r.GetUrl())

	result := &pb.Response{Url: "newshorturl"}

	if !reflect.DeepEqual(r.GetUrl(), result.GetUrl()) {
		t.Fatalf("logs2 dont match\nhave %+v\nwant %+v", r, result)
	}
}

func TestGetURL(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	defer finish()
	StartMyService(ctx, listenAddr)

	conn, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//log.Fatalf("fail to dial: %v", err)
		t.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r1, err := client.CreateShortURL(ctx, &pb.Request{Url: "http://welcome.com"})
	if err != nil {
		t.Fatalf("could not shorten url: %v", err)
	}

	//ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel1()

	r, err := client.GetFullURL(ctx, &pb.Request{Url: r1.GetUrl()})
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		t.Fatalf("could not get original url: %v", err)
	}
	//log.Printf("User: %s", r.GetName())
	fmt.Printf("Full URL: %s", r.GetUrl())

	result := &pb.Response{Url: "http://welcome.com"}

	if !reflect.DeepEqual(r.GetUrl(), result.GetUrl()) {
		t.Fatalf("logs2 dont match\nhave %+v\nwant %+v", r, result)
	}
}

func StartMyService(ctx context.Context, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
		//log.Fatalln("Cant listen port", err)
	}

	db, err := pgdb.NewPostgresDB(pgdb.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "Stacy",
	})
	if err != nil {
		//logrus.Fatalf("failed to initialize db: %s", err.Error())
		fmt.Printf("failed to initialize db: %s", err.Error())
	}
	storage := pgdb.NewStorage(db)
	service := service.NewURLServer(storage)

	server := grpc.NewServer()
	pb.RegisterURLServer(server, service)

	go func() {
		fmt.Println("Starting server at " + addr)
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

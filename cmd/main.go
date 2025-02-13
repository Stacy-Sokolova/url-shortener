package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	s "url-server/internal/app"
	"url-server/internal/service"
	"url-server/internal/storage"
	m "url-server/internal/storage/memdb"
	pg "url-server/internal/storage/pgdb"

	"github.com/dgraph-io/badger"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	//logrus

	//config
	if err := initConfig(); err != nil {
		//logrus.Fatalf("error initializing configs: %s", err.Error())
		fmt.Printf("error initializing configs: %s", err.Error())
	}

	dbType := dbParam()

	var storage storage.Storage
	var db *sqlx.DB
	var memdb *badger.DB

	if dbType == "postgres" {
		db, err := pg.NewPostgresDB(pg.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: viper.GetString("db.dbpassword"),
		})
		if err != nil {
			//logrus.Fatalf("failed to initialize db: %s", err.Error())
			fmt.Printf("failed to initialize db: %s", err.Error())
		}
		storage = pg.NewStorage(db)

	} else if dbType == "inmemory" {
		memdb, err := m.NewMemDB()
		if err != nil {
			fmt.Printf("failed to initialize db: %s", err.Error())
		}
		storage = m.NewStorage(memdb)
	} else {
		//fmt.Println("error in db param")
	}

	service := service.NewURLServer(storage)

	srv := new(s.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), service); err != nil {
			//logrus.Fatalf("error occured while running http server: %s", err.Error())
			fmt.Printf("error occured while running http server: %s", err.Error())
		}
	}()
	fmt.Println("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	srv.Shutdown()

	if dbType == "postgres" {
		if err := db.Close(); err != nil {
			//logrus.Errorf("error occured on db connection close: %s", err.Error())
			fmt.Printf("error occured on db connection close: %s", err.Error())
		}
	} else {
		if err := memdb.Close(); err != nil {
			//logrus.Errorf("error occured on db connection close: %s", err.Error())
			fmt.Printf("error occured on db connection close: %s", err.Error())
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func dbParam() string {
	var dbtype string
	if len(os.Args) > 1 {
		dbtype = os.Args[1]
	} else {
		fmt.Println("not correct param")
	}
	return dbtype
}

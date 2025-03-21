package main

import (
	"os"
	"os/signal"
	"syscall"
	s "url-server/internal/app"
	"url-server/internal/handler"
	"url-server/internal/service"
	"url-server/internal/storage"
	m "url-server/internal/storage/memdb"
	pg "url-server/internal/storage/pgdb"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	dbType := dbParam()
	var database any
	if dbType == "inmemory" {
		memdb, err := m.NewMemDB()
		if err != nil {
			logrus.Fatalf("failed to initialize in-memory db: %s", err.Error())
		}
		database = m.NewInmemStorage(memdb)
	} else {
		db, err := pg.NewPostgresDB(pg.Config{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Username: os.Getenv("POSTGRES_USER"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("DB_PASSWORD"),
		})
		if err != nil {
			logrus.Fatalf("failed to initialize db: %s", err.Error())
		}
		database = pg.NewSQLStorage(db)
	}

	storage := storage.NewStorage(database)
	service := service.NewService(storage)
	handler := handler.NewHandler(service)

	srv := new(s.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	srv.Shutdown()

	if err := storage.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func dbParam() string {
	//dbtype := ""
	/*if len(os.Args) > 1 {
		dbtype = os.Args[1]
	}*/
	return os.Getenv("STORAGE_TYPE")
}

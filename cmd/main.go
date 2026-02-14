package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DavelPurov777/todo-app-golang"
	"github.com/DavelPurov777/todo-app-golang/pkg/handler"
	"github.com/DavelPurov777/todo-app-golang/pkg/repository"
	"github.com/DavelPurov777/todo-app-golang/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running HTTP server %s", err.Error())
		}
	}()
	logrus.Print("Todo App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // строка для чтения из канала которая выполняет блокировку главной горутины main

	logrus.Print("Todo App shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occuring on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occuring on connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	// ReadInConfig считывает значения конфигов и записывает их во внутренний объект viper, а возвращает только ошибку
	return viper.ReadInConfig()
}

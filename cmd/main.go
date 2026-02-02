package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/DavelPurov777/todo-app-golang"
	"github.com/DavelPurov777/todo-app-golang/pkg/repository"
	"github.com/DavelPurov777/todo-app-golang/pkg/service"
	"github.com/DavelPurov777/todo-app-golang/pkg/handler"
	"os"
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
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize DB: %s", err.Error())
	}

	repos := repository.NewRepository(db);
	services := service.NewService(repos);
	handlers := handler.NewHandler(services);

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running HTTP server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	// ReadInConfig считывает значения конфигов и записывает их во внутренний объект viper, а возвращает только ошибку
	return viper.ReadInConfig()
}
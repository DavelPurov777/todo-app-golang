package main

import (
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/DavelPurov777/todo-app-golang"
	"github.com/DavelPurov777/todo-app-golang/pkg/repository"
	"github.com/DavelPurov777/todo-app-golang/pkg/service"
	"github.com/DavelPurov777/todo-app-golang/pkg/handler"
	"log"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables %s", err.Error())
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
		log.Fatalf("failed to initialize DB: %s", err.Error())
	}

	repos := repository.NewRepository(db);
	services := service.NewService(repos);
	handlers := handler.NewHandler(services);

	srv := new(todo.Server) // TODO: тут не понял что импортируется под todo
	if err := srv.Run(viper.GetString("8000"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running HTTP server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
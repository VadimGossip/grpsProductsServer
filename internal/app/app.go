package app

import (
	"context"
	"github.com/VadimGossip/grpsProductsServer/internal/config"
	"github.com/VadimGossip/grpsProductsServer/internal/repository"
	"github.com/VadimGossip/grpsProductsServer/internal/server"
	"github.com/VadimGossip/grpsProductsServer/internal/service"
	"github.com/VadimGossip/grpsProductsServer/pkg/database"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := database.NewMongoConnection(ctx, cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.URI)
	if err != nil {
		logrus.Fatalf("Mongo connection error %s", err)
	}
	db := dbClient.Database(cfg.Mongo.Database)

	productsRepo := repository.NewProducts(db)
	productsService := service.NewProductService(productsRepo)
	productsServer := server.NewProductServer(productsService)
	srv := server.NewServer(productsServer)

	logrus.Info("Products grps server started")
	if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
		logrus.Fatalf("error occured while running products grps server: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	srv.Stop()
	logrus.Info("Products grps server stopped")

}

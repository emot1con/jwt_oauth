package main

import (
	"auth/internal/delivery/http_delivery"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	handler := httpdelivery.NewHandler()

	srv := http.Server{
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: handler,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logrus.Infof("starting server at port %s", os.Getenv("SERVER_PORT"))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("error serving server: %v", err)
		}
	}()
	<-sig

	logrus.Info("server closed succesfully")
}

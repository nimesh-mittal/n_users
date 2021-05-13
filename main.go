package main

import (
	"n_users/handler"
	"n_users/server"

	"go.uber.org/zap"
)

func initLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()

	s := server.New()

	hh := handler.NewHealthHandler()
	s.Mount("/", hh.NewHealthRouter())

	s.StartServer(":8085")
}

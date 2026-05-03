package main

import (
	"TicTacToe/di"
	"TicTacToe/storage/config"
	"context"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/fx"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		panic(err.Error())
	}

	app := fx.New(
		fx.Provide(func() *config.DataBaseConfig { return cfg }),
		di.Module,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}

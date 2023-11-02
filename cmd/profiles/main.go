package main

import (
	"context"
	"fmt"

	"github.com/mnepesov/profiles/pkg/hasher"
	"github.com/mnepesov/profiles/service"
	"github.com/mnepesov/profiles/service/config"
	"github.com/mnepesov/profiles/service/db"
	"github.com/mnepesov/profiles/service/interceptor"
	"github.com/mnepesov/profiles/service/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig("./etc/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("read config: %v", err))
	}

	repo := db.NewInMemory()
	dataHasher := hasher.NewHasher()

	srv := server.New(cfg.Server)

	profileService := service.New(cfg, repo, dataHasher)
	srv.Register(profileService)
	srv.Serve(ctx, interceptor.AuthInterceptor(repo, dataHasher))
}

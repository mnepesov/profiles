package service

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/pkg/hasher"
	"github.com/mnepesov/profiles/service/config"
	"github.com/mnepesov/profiles/service/db"
	"github.com/mnepesov/profiles/service/domain"
	"github.com/mnepesov/profiles/service/handlers"
)

type ProfileService struct {
	handlers.Handler
}

func New(cfg *config.Root, repo *db.InMemory, dataHasher *hasher.Hasher) *ProfileService {
	hashedPassword, err := dataHasher.Get([]byte(cfg.Admin.Password))
	if err != nil {
		panic(fmt.Sprintf("can't get hash from password: %e", err))
	}

	err = repo.Add(domain.Profile{
		Id:       uuid.New().String(),
		Username: cfg.Admin.Username,
		Email:    cfg.Admin.Email,
		Password: hashedPassword,
		IsAdmin:  true,
	})
	if err != nil {
		panic(fmt.Sprintf("add admin: %v", err))
	}

	grpcHandler := handlers.Handler{
		Repo:           repo,
		Hasher:         dataHasher,
		DefaultProfile: cfg.Admin,
	}

	return &ProfileService{
		Handler: grpcHandler,
	}
}

func (p *ProfileService) GRPCRegisterer() func(s *grpc.Server) {
	return func(s *grpc.Server) {
		profiles.RegisterProfilesServer(s, p)
	}
}

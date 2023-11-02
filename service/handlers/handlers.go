package handlers

import (
	"github.com/mnepesov/profiles/service/config"
	"github.com/mnepesov/profiles/service/domain"
)

type InMemoryRepo interface {
	Get(Id string) (domain.Profile, error)
	Add(profile domain.Profile) error
	List() []domain.Profile
	Delete(id string) error
	Update(id string, profile domain.Profile) error
}

type Hasher interface {
	Get(data []byte) (string, error)
	Check(hashedData, data []byte) (bool, error)
}

type Handler struct {
	Repo           InMemoryRepo
	Hasher         Hasher
	DefaultProfile config.Admin
}

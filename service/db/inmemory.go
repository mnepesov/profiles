package db

import (
	"sync"

	"github.com/mnepesov/profiles/service/domain"
)

type InMemory struct {
	data      map[string]domain.Profile
	usernames map[string]string
	mu        sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		data:      make(map[string]domain.Profile),
		usernames: make(map[string]string),
	}
}

func (i *InMemory) Get(id string) (domain.Profile, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	profile, ok := i.data[id]
	if !ok {
		return domain.Profile{}, domain.NotFoundError
	}

	return profile, nil
}

func (i *InMemory) Add(profile domain.Profile) error {
	if i.checkUniqueUsername(profile.Username) {
		return domain.AlreadyExistError
	}
	if i.checkUniqueId(profile.Id) {
		return domain.AlreadyExistError
	}

	i.mu.Lock()
	i.data[profile.Id] = profile
	i.usernames[profile.Username] = profile.Id
	i.mu.Unlock()

	return nil
}

func (i *InMemory) Update(id string, profile domain.Profile) error {
	if _, err := i.Get(id); err != nil {
		return err
	}

	idByUsername, ok := i.usernames[profile.Username]
	if ok && idByUsername != id {
		return domain.AlreadyExistError
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[id] = profile
	i.usernames[profile.Username] = id

	return nil
}

func (i *InMemory) Delete(id string) error {
	profile, err := i.Get(id)
	if err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.data, id)
	delete(i.usernames, profile.Username)

	return nil
}

func (i *InMemory) List() []domain.Profile {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var profiles []domain.Profile
	for _, profile := range i.data {
		profiles = append(profiles, profile)
	}

	return profiles
}

func (i *InMemory) GetByUsername(username string) (domain.Profile, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	id, ok := i.usernames[username]
	if !ok {
		return domain.Profile{}, domain.NotFoundError
	}
	profile, ok := i.data[id]
	if !ok {
		return domain.Profile{}, domain.NotFoundError
	}

	return profile, nil
}

func (i *InMemory) checkUniqueUsername(username string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if _, ok := i.usernames[username]; ok {
		return true
	}

	return false
}

func (i *InMemory) checkUniqueId(id string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if _, ok := i.data[id]; ok {
		return true
	}

	return false
}

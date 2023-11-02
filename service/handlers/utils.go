package handlers

import (
	"context"
	"fmt"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/service/domain"
)

func isAdminStatusFromCtx(ctx context.Context) (bool, error) {
	isAdmin, ok := ctx.Value(domain.IsAdminCtxKey).(bool)
	if !ok {
		return false, fmt.Errorf("error get from context")
	}

	return isAdmin, nil
}

func profilesFromDomain(pList []domain.Profile) []*profiles.Profile {
	var profilesList []*profiles.Profile
	for _, profile := range pList {
		profilesList = append(profilesList, profileFromDomain(profile))
	}
	return profilesList
}

func profileFromDomain(profile domain.Profile) *profiles.Profile {
	return &profiles.Profile{
		Id:       profile.Id,
		Username: profile.Username,
		Email:    profile.Email,
		IsAdmin:  profile.IsAdmin,
	}
}

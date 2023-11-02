package handlers

import (
	"context"

	"github.com/mnepesov/profiles/api/proto/profiles"
)

func (h *Handler) GetList(ctx context.Context, req *profiles.GetListRequest) (*profiles.GetListResponse, error) {
	profilesList := h.Repo.List()
	return &profiles.GetListResponse{
		Profiles: profilesFromDomain(profilesList),
	}, nil
}

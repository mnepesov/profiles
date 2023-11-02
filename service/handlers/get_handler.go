package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/service/domain"
)

func (h *Handler) Get(ctx context.Context, req *profiles.GetRequest) (*profiles.GetResponse, error) {
	profile, err := h.Repo.Get(req.Id)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &profiles.GetResponse{
		Profile: profileFromDomain(profile),
	}, nil
}

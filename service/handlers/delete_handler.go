package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/service/domain"
)

func (h *Handler) Delete(ctx context.Context, req *profiles.DeleteRequest) (*profiles.DeleteResponse, error) {
	isAdmin, err := isAdminStatusFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	profile, err := h.Repo.Get(req.Id)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	if h.DefaultProfile.Username == profile.Username {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	err = h.Repo.Delete(req.Id)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return &profiles.DeleteResponse{}, nil
}

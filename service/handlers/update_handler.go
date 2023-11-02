package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/service/domain"
)

func (h *Handler) Update(ctx context.Context, req *profiles.UpdateRequest) (*profiles.UpdateResponse, error) {
	isAdmin, err := isAdminStatusFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	hashedPwd, err := h.Hasher.Get([]byte(req.Data.Password))
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	newProfile := domain.Profile{
		Id:       req.Id,
		Username: req.Data.Username,
		Email:    req.Data.Email,
		Password: hashedPwd,
		IsAdmin:  req.Data.IsAdmin,
	}
	err = h.Repo.Update(req.Id, newProfile)
	if err != nil {
		if errors.Is(err, domain.NotFoundError) {
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		}
		if errors.Is(err, domain.AlreadyExistError) {
			return nil, status.Error(codes.InvalidArgument, "username already exist")
		}
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &profiles.UpdateResponse{
		Profile: profileFromDomain(newProfile),
	}, nil
}

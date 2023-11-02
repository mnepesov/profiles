package handlers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mnepesov/profiles/api/proto/profiles"
	"github.com/mnepesov/profiles/service/domain"
)

func (h *Handler) Create(ctx context.Context, req *profiles.CreateRequest) (*profiles.CreateResponse, error) {
	isAdmin, err := isAdminStatusFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	hashedPwd, err := h.Hasher.Get([]byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	id := uuid.New().String()

	err = h.Repo.Add(domain.Profile{
		Id:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPwd,
		IsAdmin:  req.IsAdmin,
	})
	if err != nil {
		if errors.Is(err, domain.AlreadyExistError) {
			return nil, status.Error(codes.InvalidArgument, "username or id already exist")
		}
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &profiles.CreateResponse{
		Profile: &profiles.Profile{
			Id:       id,
			Username: req.Username,
			Email:    req.Email,
			IsAdmin:  req.IsAdmin,
		},
	}, nil
}

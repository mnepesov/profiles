package interceptor

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/mnepesov/profiles/service/domain"
)

type InMemoryRepo interface {
	GetByUsername(username string) (domain.Profile, error)
}

type Hasher interface {
	Check(hashedData, data []byte) (bool, error)
}

func AuthInterceptor(repo InMemoryRepo, dataHasher Hasher) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		username, password, ok := parseBasicAuthFromMetadata(md)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "basic authentication failed")
		}

		profile, err := repo.GetByUsername(username)
		if err != nil {
			if errors.Is(err, domain.NotFoundError) {
				return nil, status.Errorf(codes.Unauthenticated, "basic authentication failed")
			}
		}

		isCorrect, err := dataHasher.Check([]byte(profile.Password), []byte(password))
		if err != nil {
			return nil, status.Errorf(codes.Internal, codes.Internal.String())
		}

		if !isCorrect {
			return nil, status.Errorf(codes.Unauthenticated, "basic authentication failed")
		}

		ctx = context.WithValue(ctx, domain.IsAdminCtxKey, profile.IsAdmin)
		return handler(ctx, req)
	}
}

func parseBasicAuthFromMetadata(md metadata.MD) (string, string, bool) {
	auth, ok := md["authorization"]
	if !ok {
		return "", "", false
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth[0], "Basic "))
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return "", "", false
	}

	credentials := strings.Split(string(decoded), ":")
	if len(credentials) != 2 {
		return "", "", false
	}

	username := credentials[0]
	password := credentials[1]

	return username, password, true
}

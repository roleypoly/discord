package auth

import (
	"os"

	"context"

	pbAuth "github.com/roleypoly/rpc/auth/backend"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthConnector struct {
	Client pbAuth.AuthBackendClient
}

type sharedSecretCredential struct{}

func (sharedSecretCredential) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Shared " + os.Getenv("SHARED_SECRET"),
	}, nil
}

func (sharedSecretCredential) RequireTransportSecurity() bool {
	return true
}

func NewAuthConnector() (*AuthConnector, error) {
	var cred sharedSecretCredential

	grpcConn, err := grpc.Dial(
		os.Getenv("AUTH_SVC_DIAL"),
		grpc.WithPerRPCCredentials(cred),
		grpc.WithTransportCredentials(credentials.NewTLS(nil)),
	)
	if err != nil {
		return nil, err
	}
	return &AuthConnector{
		Client: pbAuth.NewAuthBackendClient(grpcConn),
	}, nil
}

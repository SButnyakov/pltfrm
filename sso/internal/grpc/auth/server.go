package auth

import (
	"context"
	ssov1 "github.com/SButnyakov/pltfrm-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(in); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := validateRegister(in)
	if err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(in); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(in *ssov1.LoginRequest) error {
	if in.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "no email provided")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "no password provided")
	}

	if in.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "no app id provided")
	}

	return nil
}

func validateRegister(in *ssov1.RegisterRequest) error {
	if in.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "no email provided")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "no password provided")
	}

	return nil
}

func validateIsAdmin(in *ssov1.IsAdminRequest) error {
	if in.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "no user id provided")
	}

	return nil
}
